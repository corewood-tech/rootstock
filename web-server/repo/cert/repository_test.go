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
)

func setupTest(t *testing.T) Repository {
	t.Helper()

	// Generate a test CA in a temp dir
	dir := t.TempDir()

	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate ca key: %v", err)
	}

	caTemplate := &x509.Certificate{
		SerialNumber:          mustSerial(t),
		Subject:               pkix.Name{CommonName: "Test CA"},
		NotBefore:             mustNow(),
		NotAfter:              mustNow().AddDate(1, 0, 0),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create ca cert: %v", err)
	}

	// Write CA cert PEM
	certPath := filepath.Join(dir, "ca.crt")
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER})
	if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
		t.Fatalf("write ca cert: %v", err)
	}

	// Write CA key PEM
	keyPath := filepath.Join(dir, "ca.key")
	keyDER, err := x509.MarshalECPrivateKey(caKey)
	if err != nil {
		t.Fatalf("marshal ca key: %v", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		t.Fatalf("write ca key: %v", err)
	}

	cfg := config.CertConfig{
		CACertPath:       certPath,
		CAKeyPath:        keyPath,
		CertLifetimeDays: 90,
	}

	repo, err := NewRepository(cfg)
	if err != nil {
		t.Fatalf("NewRepository(): %v", err)
	}
	t.Cleanup(func() {
		repo.Shutdown()
	})

	return repo
}

func TestIssueCert(t *testing.T) {
	repo := setupTest(t)
	ctx := context.Background()

	csr := generateTestCSR(t)

	issued, err := repo.IssueCert(ctx, IssueCertInput{
		DeviceID: "device-abc-123",
		CSR:      csr,
	})
	if err != nil {
		t.Fatalf("IssueCert(): %v", err)
	}

	if len(issued.CertPEM) == 0 {
		t.Fatal("IssueCert() returned empty CertPEM")
	}
	if issued.Serial == "" {
		t.Fatal("IssueCert() returned empty Serial")
	}

	// Parse issued cert and verify properties
	block, _ := pem.Decode(issued.CertPEM)
	if block == nil {
		t.Fatal("failed to decode issued cert PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("parse issued cert: %v", err)
	}

	if cert.Subject.CommonName != "device-abc-123" {
		t.Errorf("CN = %q, want %q", cert.Subject.CommonName, "device-abc-123")
	}
	if cert.Issuer.CommonName != "Test CA" {
		t.Errorf("Issuer CN = %q, want %q", cert.Issuer.CommonName, "Test CA")
	}
	if cert.KeyUsage&x509.KeyUsageDigitalSignature == 0 {
		t.Error("expected DigitalSignature key usage")
	}
	if len(cert.ExtKeyUsage) != 1 || cert.ExtKeyUsage[0] != x509.ExtKeyUsageClientAuth {
		t.Error("expected ClientAuth extended key usage")
	}
}

func TestIssueCertSetsDeviceIDNotCSRSubject(t *testing.T) {
	repo := setupTest(t)
	ctx := context.Background()

	// CSR has CN=some-csr-cn but we pass DeviceID=real-device-id
	csr := generateTestCSR(t)

	issued, err := repo.IssueCert(ctx, IssueCertInput{
		DeviceID: "real-device-id",
		CSR:      csr,
	})
	if err != nil {
		t.Fatalf("IssueCert(): %v", err)
	}

	block, _ := pem.Decode(issued.CertPEM)
	cert, _ := x509.ParseCertificate(block.Bytes)

	if cert.Subject.CommonName != "real-device-id" {
		t.Errorf("CN = %q, want %q (should use DeviceID, not CSR subject)", cert.Subject.CommonName, "real-device-id")
	}
}

func TestIssueCertRejectsWeakKey(t *testing.T) {
	repo := setupTest(t)
	ctx := context.Background()

	// Generate a P-224 key (too small, < 256 bits)
	key, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	template := &x509.CertificateRequest{Subject: pkix.Name{CommonName: "weak"}}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, template, key)
	if err != nil {
		t.Fatalf("create csr: %v", err)
	}

	_, err = repo.IssueCert(ctx, IssueCertInput{DeviceID: "dev-1", CSR: csrDER})
	if err == nil {
		t.Fatal("IssueCert() should reject weak key")
	}
}

func TestGetCACert(t *testing.T) {
	repo := setupTest(t)
	ctx := context.Background()

	ca, err := repo.GetCACert(ctx)
	if err != nil {
		t.Fatalf("GetCACert(): %v", err)
	}

	block, _ := pem.Decode(ca.CertPEM)
	if block == nil {
		t.Fatal("GetCACert() returned invalid PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("parse ca cert: %v", err)
	}
	if cert.Subject.CommonName != "Test CA" {
		t.Errorf("CA CN = %q, want %q", cert.Subject.CommonName, "Test CA")
	}
}

// --- helpers ---

func generateTestCSR(t *testing.T) []byte {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	template := &x509.CertificateRequest{
		Subject: pkix.Name{CommonName: "some-csr-cn"},
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, template, key)
	if err != nil {
		t.Fatalf("create csr: %v", err)
	}
	return csrDER
}

func mustSerial(t *testing.T) *big.Int {
	t.Helper()
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	s, err := rand.Int(rand.Reader, max)
	if err != nil {
		t.Fatalf("generate serial: %v", err)
	}
	return s
}

func mustNow() time.Time {
	return time.Now().UTC()
}
