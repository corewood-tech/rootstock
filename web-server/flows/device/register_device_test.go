package device

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"rootstock/web-server/config"
	certops "rootstock/web-server/ops/cert"
	deviceops "rootstock/web-server/ops/device"
	certrepo "rootstock/web-server/repo/cert"
	devicerepo "rootstock/web-server/repo/device"
	sqlmigrate "rootstock/web-server/repo/sql/migrate"
)

func setupRegisterTest(t *testing.T) (*RegisterDeviceFlow, *deviceops.Ops, *pgxpool.Pool) {
	t.Helper()
	pgCfg := config.PostgresConfig{
		Host: "app-postgres", Port: 5432, User: "rootstock", Password: "rootstock", DBName: "rootstock", SSLMode: "disable",
	}
	if err := sqlmigrate.Run(pgCfg); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		pgCfg.User, pgCfg.Password, pgCfg.Host, pgCfg.Port, pgCfg.DBName, pgCfg.SSLMode,
	)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	ctx := context.Background()
	pool.Exec(ctx, "TRUNCATE devices CASCADE")

	dRepo := devicerepo.NewRepository(pool)
	dOps := deviceops.NewOps(dRepo)

	cRepo := setupCertRepo(t)
	cOps := certops.NewOps(cRepo)

	flow := NewRegisterDeviceFlow(dOps, cOps)

	t.Cleanup(func() {
		cRepo.Shutdown()
		dRepo.Shutdown()
		pool.Close()
	})

	return flow, dOps, pool
}

func TestRegisterDevice(t *testing.T) {
	flow, dOps, pool := setupRegisterTest(t)
	ctx := context.Background()

	// Create a pending device and enrollment code
	device, err := dOps.CreateDevice(ctx, deviceops.CreateDeviceInput{
		OwnerID: "user-1", Class: "sensor", FirmwareVersion: "1.0.0", Tier: 1, Sensors: []string{"temp"},
	})
	if err != nil {
		t.Fatalf("CreateDevice(): %v", err)
	}

	code, err := dOps.GenerateEnrollmentCode(ctx, deviceops.GenerateCodeInput{
		DeviceID: device.ID, Code: "REG-TEST-001", TTL: 900,
	})
	if err != nil {
		t.Fatalf("GenerateEnrollmentCode(): %v", err)
	}

	csr := generateCSR(t)

	result, err := flow.Run(ctx, RegisterDeviceInput{
		EnrollmentCode: code.Code,
		CSR:            csr,
	})
	if err != nil {
		t.Fatalf("Run(): %v", err)
	}

	if result.DeviceID != device.ID {
		t.Errorf("DeviceID = %q, want %q", result.DeviceID, device.ID)
	}
	if len(result.CertPEM) == 0 {
		t.Error("expected non-empty CertPEM")
	}
	if result.Serial == "" {
		t.Error("expected non-empty Serial")
	}

	// Verify cert CN = device ID
	block, _ := pem.Decode(result.CertPEM)
	cert, _ := x509.ParseCertificate(block.Bytes)
	if cert.Subject.CommonName != device.ID {
		t.Errorf("cert CN = %q, want %q", cert.Subject.CommonName, device.ID)
	}

	// Verify device is now active with cert serial
	var status, certSerial string
	pool.QueryRow(ctx, "SELECT status, cert_serial FROM devices WHERE id = $1", device.ID).Scan(&status, &certSerial)
	if status != "active" {
		t.Errorf("device status = %q, want %q", status, "active")
	}
	if certSerial != result.Serial {
		t.Errorf("cert_serial = %q, want %q", certSerial, result.Serial)
	}

	// Verify enrollment code is used
	var used bool
	pool.QueryRow(ctx, "SELECT used FROM enrollment_codes WHERE code = $1", code.Code).Scan(&used)
	if !used {
		t.Error("enrollment code should be marked used")
	}
}

func TestRegisterDeviceInvalidCode(t *testing.T) {
	flow, _, _ := setupRegisterTest(t)
	ctx := context.Background()

	csr := generateCSR(t)

	_, err := flow.Run(ctx, RegisterDeviceInput{
		EnrollmentCode: "nonexistent-code",
		CSR:            csr,
	})
	if err == nil {
		t.Fatal("expected error for invalid enrollment code")
	}
}

// --- helpers ---

func setupCertRepo(t *testing.T) certrepo.Repository {
	t.Helper()
	dir := t.TempDir()

	caKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
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
	caDER, _ := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)

	certPath := filepath.Join(dir, "ca.crt")
	os.WriteFile(certPath, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER}), 0644)

	keyPath := filepath.Join(dir, "ca.key")
	keyDER, _ := x509.MarshalECPrivateKey(caKey)
	os.WriteFile(keyPath, pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}), 0600)

	repo, err := certrepo.NewRepository(config.CertConfig{
		CACertPath: certPath, CAKeyPath: keyPath, CertLifetimeDays: 90,
	})
	if err != nil {
		t.Fatalf("NewRepository(): %v", err)
	}
	return repo
}

func generateCSR(t *testing.T) []byte {
	t.Helper()
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		Subject: pkix.Name{CommonName: "test"},
	}, key)
	if err != nil {
		t.Fatalf("create csr: %v", err)
	}
	return csrDER
}
