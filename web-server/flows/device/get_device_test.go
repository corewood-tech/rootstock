package device

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"rootstock/web-server/config"
	deviceops "rootstock/web-server/ops/device"
	devicerepo "rootstock/web-server/repo/device"
	sqlmigrate "rootstock/web-server/repo/sql/migrate"
)

func setupDeviceFlowTest(t *testing.T) (*deviceops.Ops, *pgxpool.Pool) {
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
	t.Cleanup(func() {
		dRepo.Shutdown()
		pool.Close()
	})
	return dOps, pool
}

func TestGetDevice(t *testing.T) {
	dOps, _ := setupDeviceFlowTest(t)
	ctx := context.Background()

	created, err := dOps.CreateDevice(ctx, deviceops.CreateDeviceInput{
		OwnerID: "user-1", Class: "sensor", FirmwareVersion: "1.0.0", Tier: 1, Sensors: []string{"temp"},
	})
	if err != nil {
		t.Fatalf("CreateDevice(): %v", err)
	}

	flow := NewGetDeviceFlow(dOps)
	got, err := flow.Run(ctx, GetDeviceInput{DeviceID: created.ID})
	if err != nil {
		t.Fatalf("Run(): %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("ID = %q, want %q", got.ID, created.ID)
	}
	if got.Class != "sensor" {
		t.Errorf("Class = %q, want %q", got.Class, "sensor")
	}
}

func TestRevokeDevice(t *testing.T) {
	dOps, _ := setupDeviceFlowTest(t)
	ctx := context.Background()

	created, _ := dOps.CreateDevice(ctx, deviceops.CreateDeviceInput{
		OwnerID: "user-1", Class: "sensor", FirmwareVersion: "1.0.0", Tier: 1, Sensors: []string{"temp"},
	})
	dOps.UpdateDeviceStatus(ctx, created.ID, "active")

	flow := NewRevokeDeviceFlow(dOps)
	if err := flow.Run(ctx, RevokeDeviceInput{DeviceID: created.ID}); err != nil {
		t.Fatalf("Run(): %v", err)
	}

	got, _ := dOps.GetDevice(ctx, created.ID)
	if got.Status != "revoked" {
		t.Errorf("status = %q, want %q", got.Status, "revoked")
	}
}

func TestReinstateDevice(t *testing.T) {
	dOps, _ := setupDeviceFlowTest(t)
	ctx := context.Background()

	created, _ := dOps.CreateDevice(ctx, deviceops.CreateDeviceInput{
		OwnerID: "user-1", Class: "sensor", FirmwareVersion: "1.0.0", Tier: 1, Sensors: []string{"temp"},
	})
	dOps.UpdateDeviceStatus(ctx, created.ID, "active")
	dOps.UpdateDeviceStatus(ctx, created.ID, "revoked")

	flow := NewReinstateDeviceFlow(dOps)
	if err := flow.Run(ctx, ReinstateDeviceInput{DeviceID: created.ID}); err != nil {
		t.Fatalf("Run(): %v", err)
	}

	got, _ := dOps.GetDevice(ctx, created.ID)
	if got.Status != "active" {
		t.Errorf("status = %q, want %q", got.Status, "active")
	}
}
