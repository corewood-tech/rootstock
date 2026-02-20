package server

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	mochi "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"

	"rootstock/web-server/config"
)

// NewMQTTServer creates an embedded Mochi MQTT broker with:
//   - InlineClient enabled (for server-side publish/subscribe)
//   - mTLS auth hook (device identity from cert CN, topic ACL)
//   - TLS listener on cfg.MQTT.Port with RequireAndVerifyClientCert
//
// The broker generates an ephemeral server certificate from the CA at startup.
// Devices trust the CA, so they trust the server cert. No persistent server cert needed.
//
// Call Serve() on the returned server to start accepting connections.
// The cleanup function closes the broker gracefully.
func NewMQTTServer(cfg *config.Config) (*mochi.Server, func(), error) {
	// Create broker with inline client
	server := mochi.New(&mochi.Options{
		InlineClient: true,
	})

	// Load CA cert + key
	caCert, caSigner, caCertPEM, err := loadCA(cfg.Cert.CACertPath, cfg.Cert.CAKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("load ca for mqtt: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCertPEM) {
		return nil, nil, fmt.Errorf("failed to add ca cert to mqtt pool")
	}

	// Add mTLS auth hook
	if err := server.AddHook(&MQTTAuthHook{}, &MQTTAuthHookConfig{
		CACertPool:      caCertPool,
		GracePeriodDays: cfg.MQTT.GracePeriodDays,
	}); err != nil {
		return nil, nil, fmt.Errorf("add mqtt auth hook: %w", err)
	}

	// Generate ephemeral server cert signed by the CA
	serverCert, err := generateServerCert(caCert, caSigner, cfg.MQTT.ServerSANs)
	if err != nil {
		return nil, nil, fmt.Errorf("generate mqtt server cert: %w", err)
	}

	// TLS listener with mTLS + grace period for expired certs.
	// RequireAnyClientCert ensures a cert is presented but skips Go's built-in
	// expiry check. VerifyPeerCertificate does full chain validation with a
	// relaxed expiry window so devices with recently-expired certs can still
	// connect to renew.
	gracePeriodDays := cfg.MQTT.GracePeriodDays
	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAnyClientCert,
		VerifyPeerCertificate: func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
			if len(rawCerts) == 0 {
				return fmt.Errorf("no client certificate")
			}
			cert, err := x509.ParseCertificate(rawCerts[0])
			if err != nil {
				return fmt.Errorf("parse client cert: %w", err)
			}

			// Beyond grace period → reject
			now := time.Now()
			graceCutoff := cert.NotAfter.Add(time.Duration(gracePeriodDays) * 24 * time.Hour)
			if now.After(graceCutoff) {
				return fmt.Errorf("certificate expired beyond grace period")
			}

			// Verify chain with time clamped to just before expiry
			// (bypasses Go's expiry check while still validating chain + signature)
			verifyTime := now
			if now.After(cert.NotAfter) {
				verifyTime = cert.NotAfter.Add(-time.Second)
			}

			intermediates := x509.NewCertPool()
			for _, raw := range rawCerts[1:] {
				if ic, err := x509.ParseCertificate(raw); err == nil {
					intermediates.AddCert(ic)
				}
			}

			_, err = cert.Verify(x509.VerifyOptions{
				Roots:         caCertPool,
				Intermediates: intermediates,
				CurrentTime:   verifyTime,
				KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			})
			return err
		},
	}

	tcp := listeners.NewTCP(listeners.Config{
		ID:        "mqtt-tls",
		Address:   fmt.Sprintf(":%d", cfg.MQTT.Port),
		TLSConfig: tlsCfg,
	})
	if err := server.AddListener(tcp); err != nil {
		return nil, nil, fmt.Errorf("add mqtt listener: %w", err)
	}

	cleanup := func() {
		server.Close()
	}

	return server, cleanup, nil
}

// loadCA reads and parses the CA cert and key files.
func loadCA(certPath, keyPath string) (*x509.Certificate, crypto.Signer, []byte, error) {
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read ca cert: %w", err)
	}
	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return nil, nil, nil, fmt.Errorf("decode ca cert: no PEM block")
	}
	caCert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parse ca cert: %w", err)
	}

	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read ca key: %w", err)
	}
	var keyBlock *pem.Block
	rest := keyPEM
	for {
		keyBlock, rest = pem.Decode(rest)
		if keyBlock == nil {
			return nil, nil, nil, fmt.Errorf("decode ca key: no private key block found")
		}
		if keyBlock.Type == "EC PRIVATE KEY" {
			break
		}
		if len(rest) == 0 {
			return nil, nil, nil, fmt.Errorf("decode ca key: no EC PRIVATE KEY block found")
		}
	}
	caKey, err := x509.ParseECPrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parse ca key: %w", err)
	}

	return caCert, caKey, certPEM, nil
}

// generateServerCert creates an ephemeral ECDSA server certificate signed by
// the CA for the MQTT TLS listener. Valid for 24 hours — the broker restarts
// will generate a new one.
func generateServerCert(caCert *x509.Certificate, caSigner crypto.Signer, sans []string) (tls.Certificate, error) {
	if len(sans) == 0 {
		return tls.Certificate{}, fmt.Errorf("server_sans must not be empty")
	}

	serverKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("generate server key: %w", err)
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("generate serial: %w", err)
	}

	now := time.Now().UTC()
	template := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: sans[0]},
		DNSNames:     sans,
		NotBefore:    now,
		NotAfter:     now.Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, caCert, &serverKey.PublicKey, caSigner)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("sign server cert: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyDER, err := x509.MarshalECPrivateKey(serverKey)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("marshal server key: %w", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})

	return tls.X509KeyPair(certPEM, keyPEM)
}
