package cert

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"rootstock/web-server/config"
	certrepo "rootstock/web-server/repo/cert"
)

func setupTest(t *testing.T) *Ops {
	t.Helper()

	dir := t.TempDir()

	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate ca key: %v", err)
	}

	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serial, _ := rand.Int(rand.Reader, max)

	caTemplate := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: "Test CA"},
		NotBefore:             time.Now().UTC(),
		NotAfter:              time.Now().UTC().AddDate(1, 0, 0),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create ca cert: %v", err)
	}

	certPath := filepath.Join(dir, "ca.crt")
	os.WriteFile(certPath, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER}), 0644)

	keyPath := filepath.Join(dir, "ca.key")
	keyDER, _ := x509.MarshalECPrivateKey(caKey)
	os.WriteFile(keyPath, pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}), 0600)

	repo, err := certrepo.NewRepository(config.CertConfig{
		CACertPath:       certPath,
		CAKeyPath:        keyPath,
		CertLifetimeDays: 30,
	})
	if err != nil {
		t.Fatalf("NewRepository(): %v", err)
	}

	ops := NewOps(repo)
	t.Cleanup(func() {
		repo.Shutdown()
	})
	return ops
}

func TestIssueCert(t *testing.T) {
	ops := setupTest(t)
	ctx := context.Background()

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	csrDER, _ := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		Subject: pkix.Name{CommonName: "ignored"},
	}, key)

	issued, err := ops.IssueCert(ctx, IssueCertInput{
		DeviceID: "dev-42",
		CSR:      csrDER,
	})
	if err != nil {
		t.Fatalf("IssueCert(): %v", err)
	}

	if issued.Serial == "" {
		t.Error("expected non-empty serial")
	}

	// Verify 30-day lifetime
	diff := issued.NotAfter.Sub(issued.NotBefore)
	if diff.Hours() < 29*24 || diff.Hours() > 31*24 {
		t.Errorf("lifetime = %v, want ~30 days", diff)
	}
}

func TestGetCACert(t *testing.T) {
	ops := setupTest(t)
	ctx := context.Background()

	ca, err := ops.GetCACert(ctx)
	if err != nil {
		t.Fatalf("GetCACert(): %v", err)
	}
	if len(ca.CertPEM) == 0 {
		t.Error("expected non-empty CertPEM")
	}
}
