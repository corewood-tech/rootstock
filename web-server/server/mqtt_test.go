package server

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"rootstock/web-server/config"
)

func writeTestCA(t *testing.T, dir string) (*x509.Certificate, *ecdsa.PrivateKey) {
	t.Helper()

	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate ca key: %v", err)
	}

	serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	caTemplate := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: "Test MQTT CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create ca cert: %v", err)
	}
	caCert, _ := x509.ParseCertificate(caDER)

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER})
	os.WriteFile(filepath.Join(dir, "ca.crt"), certPEM, 0644)

	keyDER, _ := x509.MarshalECPrivateKey(caKey)
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	os.WriteFile(filepath.Join(dir, "ca.key"), keyPEM, 0600)

	return caCert, caKey
}

func makeDeviceTLSCert(t *testing.T, caCert *x509.Certificate, caKey *ecdsa.PrivateKey, deviceID string) tls.Certificate {
	t.Helper()

	devKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))

	template := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: deviceID},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().AddDate(0, 3, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	devDER, err := x509.CreateCertificate(rand.Reader, template, caCert, &devKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create device cert: %v", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: devDER})
	keyDER, _ := x509.MarshalECPrivateKey(devKey)
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})

	pair, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		t.Fatalf("x509 key pair: %v", err)
	}
	return pair
}

func TestNewMQTTServer_StartsAndAcceptsMTLS(t *testing.T) {
	dir := t.TempDir()
	caCert, caKey := writeTestCA(t, dir)

	// Use a random available port
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("find free port: %v", err)
	}
	port := lis.Addr().(*net.TCPAddr).Port
	lis.Close()

	cfg := &config.Config{
		Cert: config.CertConfig{
			CACertPath: filepath.Join(dir, "ca.crt"),
			CAKeyPath:  filepath.Join(dir, "ca.key"),
		},
		MQTT: config.MQTTConfig{
			Port:       port,
			ServerSANs: []string{"localhost"},
		},
	}

	mqttServer, cleanup, err := NewMQTTServer(cfg)
	if err != nil {
		t.Fatalf("NewMQTTServer(): %v", err)
	}
	defer cleanup()

	go mqttServer.Serve()

	// Give broker a moment to start listening
	time.Sleep(200 * time.Millisecond)

	// Build TLS config for device with valid cert
	caPool := x509.NewCertPool()
	caPool.AddCert(caCert)
	deviceCert := makeDeviceTLSCert(t, caCert, caKey, "test-device")

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{deviceCert},
		RootCAs:      caPool,
		ServerName:   "localhost",
	}

	// Connect with valid mTLS cert — should succeed
	conn, err := tls.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port), tlsCfg)
	if err != nil {
		t.Fatalf("TLS dial with valid cert failed: %v", err)
	}
	conn.Close()
}

func TestNewMQTTServer_RejectsNoCert(t *testing.T) {
	dir := t.TempDir()
	caCert, _ := writeTestCA(t, dir)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("find free port: %v", err)
	}
	port := lis.Addr().(*net.TCPAddr).Port
	lis.Close()

	cfg := &config.Config{
		Cert: config.CertConfig{
			CACertPath: filepath.Join(dir, "ca.crt"),
			CAKeyPath:  filepath.Join(dir, "ca.key"),
		},
		MQTT: config.MQTTConfig{
			Port:       port,
			ServerSANs: []string{"localhost"},
		},
	}

	mqttServer, cleanup, err := NewMQTTServer(cfg)
	if err != nil {
		t.Fatalf("NewMQTTServer(): %v", err)
	}
	defer cleanup()

	go mqttServer.Serve()
	time.Sleep(200 * time.Millisecond)

	// Connect without client cert — should be rejected by TLS handshake
	caPool := x509.NewCertPool()
	caPool.AddCert(caCert)

	tlsCfg := &tls.Config{
		RootCAs:    caPool,
		ServerName: "localhost",
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port), tlsCfg)
	if err != nil {
		// TLS handshake rejected — this is the expected outcome
		return
	}
	defer conn.Close()

	// Some TLS implementations complete the handshake but fail on first I/O.
	// Try to trigger the server-side rejection by doing a handshake explicitly.
	if err := conn.Handshake(); err != nil {
		return // rejected during handshake — expected
	}

	// If we got here, the connection was established without a client cert.
	// Read to trigger the server-side close.
	buf := make([]byte, 1)
	_, err = conn.Read(buf)
	if err == nil {
		t.Fatal("expected connection to be rejected without client cert")
	}
}
