package device

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"rootstock/web-server/config"
	devicerepo "rootstock/web-server/repo/device"
	sqlmigrate "rootstock/web-server/repo/sql/migrate"
)

func setupTest(t *testing.T) (*Ops, *pgxpool.Pool) {
	t.Helper()
	cfg := config.PostgresConfig{
		Host: "app-postgres", Port: 5432, User: "rootstock", Password: "rootstock", DBName: "rootstock", SSLMode: "disable",
	}
	if err := sqlmigrate.Run(cfg); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	ctx := context.Background()
	pool.Exec(ctx, "TRUNCATE devices CASCADE")

	repo := devicerepo.NewRepository(pool)
	ops := NewOps(repo)
	t.Cleanup(func() {
		repo.Shutdown()
		pool.Close()
	})
	return ops, pool
}

func TestCreateAndGetDevice(t *testing.T) {
	ops, _ := setupTest(t)
	ctx := context.Background()

	d, err := ops.CreateDevice(ctx, CreateDeviceInput{
		OwnerID: "user-1", Class: "weather-station", FirmwareVersion: "1.0.0", Tier: 1, Sensors: []string{"temp"},
	})
	if err != nil {
		t.Fatalf("CreateDevice(): %v", err)
	}

	got, err := ops.GetDevice(ctx, d.ID)
	if err != nil {
		t.Fatalf("GetDevice(): %v", err)
	}
	if got.Class != "weather-station" {
		t.Errorf("class = %q, want %q", got.Class, "weather-station")
	}
}

func TestUpdateDeviceStatus(t *testing.T) {
	ops, _ := setupTest(t)
	ctx := context.Background()

	d, _ := ops.CreateDevice(ctx, CreateDeviceInput{
		OwnerID: "user-1", Class: "sensor", FirmwareVersion: "1.0.0", Tier: 1, Sensors: []string{"temp"},
	})

	if err := ops.UpdateDeviceStatus(ctx, d.ID, "active"); err != nil {
		t.Fatalf("UpdateDeviceStatus(): %v", err)
	}

	got, _ := ops.GetDevice(ctx, d.ID)
	if got.Status != "active" {
		t.Errorf("status = %q, want %q", got.Status, "active")
	}
}

func TestEnrollmentCodeLifecycle(t *testing.T) {
	ops, _ := setupTest(t)
	ctx := context.Background()

	d, _ := ops.CreateDevice(ctx, CreateDeviceInput{
		OwnerID: "user-1", Class: "sensor", FirmwareVersion: "1.0.0", Tier: 1, Sensors: []string{"temp"},
	})

	code, err := ops.GenerateEnrollmentCode(ctx, GenerateCodeInput{
		DeviceID: d.ID, Code: "XYZ789", TTL: 900,
	})
	if err != nil {
		t.Fatalf("GenerateEnrollmentCode(): %v", err)
	}
	if code.Used {
		t.Error("new code should not be used")
	}

	redeemed, err := ops.RedeemEnrollmentCode(ctx, "XYZ789")
	if err != nil {
		t.Fatalf("RedeemEnrollmentCode(): %v", err)
	}
	if !redeemed.Used {
		t.Error("redeemed code should be used")
	}
}

func TestGetCapabilitiesAndQuery(t *testing.T) {
	ops, _ := setupTest(t)
	ctx := context.Background()

	d, _ := ops.CreateDevice(ctx, CreateDeviceInput{
		OwnerID: "user-1", Class: "air-quality", FirmwareVersion: "2.0.0", Tier: 2, Sensors: []string{"pm25", "pm10"},
	})

	caps, err := ops.GetDeviceCapabilities(ctx, d.ID)
	if err != nil {
		t.Fatalf("GetDeviceCapabilities(): %v", err)
	}
	if caps.Tier != 2 {
		t.Errorf("tier = %d, want 2", caps.Tier)
	}

	devices, err := ops.QueryDevicesByClass(ctx, QueryByClassInput{Class: "air-quality"})
	if err != nil {
		t.Fatalf("QueryDevicesByClass(): %v", err)
	}
	if len(devices) != 1 {
		t.Errorf("query result = %d, want 1", len(devices))
	}
}
