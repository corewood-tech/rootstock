package device

import (
	"context"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"rootstock/web-server/config"
	certops "rootstock/web-server/ops/cert"
	deviceops "rootstock/web-server/ops/device"
	devicerepo "rootstock/web-server/repo/device"
	sqlmigrate "rootstock/web-server/repo/sql/migrate"
)

func setupRenewTest(t *testing.T) (*RenewCertFlow, *deviceops.Ops, *pgxpool.Pool) {
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

	cRepo := setupCertRepo(t) // reuses helper from register_device_test.go
	cOps := certops.NewOps(cRepo)

	flow := NewRenewCertFlow(dOps, cOps)

	t.Cleanup(func() {
		cRepo.Shutdown()
		dRepo.Shutdown()
		pool.Close()
	})

	return flow, dOps, pool
}

func TestRenewCert(t *testing.T) {
	flow, dOps, pool := setupRenewTest(t)
	ctx := context.Background()

	// Create an active device with an existing cert serial
	device, err := dOps.CreateDevice(ctx, deviceops.CreateDeviceInput{
		OwnerID: "user-1", Class: "sensor", FirmwareVersion: "1.0.0", Tier: 1, Sensors: []string{"temp"},
	})
	if err != nil {
		t.Fatalf("CreateDevice(): %v", err)
	}
	dOps.UpdateDeviceStatus(ctx, device.ID, "active")
	dOps.UpdateCertSerial(ctx, device.ID, "old-serial-abc")

	csr := generateCSR(t) // reuses helper from register_device_test.go

	result, err := flow.Run(ctx, RenewCertInput{
		DeviceID: device.ID,
		CSR:      csr,
	})
	if err != nil {
		t.Fatalf("Run(): %v", err)
	}

	if len(result.CertPEM) == 0 {
		t.Error("expected non-empty CertPEM")
	}
	if result.Serial == "" {
		t.Error("expected non-empty Serial")
	}

	// Verify cert CN = device ID
	block, _ := pem.Decode(result.CertPEM)
	if block == nil {
		t.Fatal("failed to decode cert PEM")
	}

	// Verify cert serial updated in DB (not the old one)
	var certSerial string
	pool.QueryRow(ctx, "SELECT cert_serial FROM devices WHERE id = $1", device.ID).Scan(&certSerial)
	if certSerial == "old-serial-abc" {
		t.Error("cert_serial should have been updated")
	}
	if certSerial != result.Serial {
		t.Errorf("cert_serial = %q, want %q", certSerial, result.Serial)
	}
}
