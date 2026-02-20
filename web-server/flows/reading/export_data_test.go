package reading

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"rootstock/web-server/config"
	readingops "rootstock/web-server/ops/reading"
	readingrepo "rootstock/web-server/repo/reading"
	sqlmigrate "rootstock/web-server/repo/sql/migrate"
)

func setupExportTest(t *testing.T) (*ExportDataFlow, *pgxpool.Pool) {
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
	pool.Exec(ctx, "TRUNCATE readings CASCADE")

	rRepo := readingrepo.NewRepository(pool)
	rOps := readingops.NewOps(rRepo)
	flow := NewExportDataFlow(rOps)

	t.Cleanup(func() {
		rRepo.Shutdown()
		pool.Close()
	})

	return flow, pool
}

func TestExportDataPseudonymizes(t *testing.T) {
	flow, pool := setupExportTest(t)
	ctx := context.Background()

	// Ensure device exists for FK
	pool.Exec(ctx, `INSERT INTO devices (id, owner_id, status, class, firmware_version, tier, sensors)
		VALUES ('dev-exp-1', 'sci-1', 'active', 'tier1', '1.0.0', 1, '{"temperature"}')
		ON CONFLICT (id) DO NOTHING`)

	now := time.Now()
	pool.Exec(ctx, `INSERT INTO readings (id, device_id, campaign_id, value, timestamp, firmware_version, cert_serial, status)
		VALUES ('r-exp-1', 'dev-exp-1', 'camp-exp-1', 22.5, $1, '1.0.0', 'serial-1', 'accepted')`, now)
	pool.Exec(ctx, `INSERT INTO readings (id, device_id, campaign_id, value, timestamp, firmware_version, cert_serial, status)
		VALUES ('r-exp-2', 'dev-exp-1', 'camp-exp-1', 23.1, $1, '1.0.0', 'serial-1', 'quarantined')`, now)

	result, err := flow.Run(ctx, ExportDataInput{
		CampaignID: "camp-exp-1",
		Secret:     "test-secret",
		Limit:      100,
		Offset:     0,
	})
	if err != nil {
		t.Fatalf("Run(): %v", err)
	}

	// Should only have accepted readings
	if len(result.Readings) != 1 {
		t.Fatalf("expected 1 reading (accepted only), got %d", len(result.Readings))
	}

	// Device ID should be pseudonymized
	r := result.Readings[0]
	if r.PseudoDeviceID == "dev-exp-1" {
		t.Error("device ID should be pseudonymized, not the original")
	}
	if r.PseudoDeviceID == "" {
		t.Error("pseudo device ID should not be empty")
	}
	if r.Value != 22.5 {
		t.Errorf("value = %f, want 22.5", r.Value)
	}
}

func TestExportDataPagination(t *testing.T) {
	flow, pool := setupExportTest(t)
	ctx := context.Background()

	pool.Exec(ctx, `INSERT INTO devices (id, owner_id, status, class, firmware_version, tier, sensors)
		VALUES ('dev-exp-2', 'sci-1', 'active', 'tier1', '1.0.0', 1, '{"temperature"}')
		ON CONFLICT (id) DO NOTHING`)

	now := time.Now()
	for i := 0; i < 5; i++ {
		pool.Exec(ctx, `INSERT INTO readings (id, device_id, campaign_id, value, timestamp, firmware_version, cert_serial, status)
			VALUES ($1, 'dev-exp-2', 'camp-exp-2', $2, $3, '1.0.0', 'serial-1', 'accepted')`,
			fmt.Sprintf("r-page-%d", i), float64(20+i), now.Add(time.Duration(-i)*time.Minute))
	}

	// First page
	result1, err := flow.Run(ctx, ExportDataInput{
		CampaignID: "camp-exp-2",
		Secret:     "test-secret",
		Limit:      2,
		Offset:     0,
	})
	if err != nil {
		t.Fatalf("Run() page 1: %v", err)
	}
	if len(result1.Readings) != 2 {
		t.Errorf("page 1: expected 2 readings, got %d", len(result1.Readings))
	}

	// Second page
	result2, err := flow.Run(ctx, ExportDataInput{
		CampaignID: "camp-exp-2",
		Secret:     "test-secret",
		Limit:      2,
		Offset:     2,
	})
	if err != nil {
		t.Fatalf("Run() page 2: %v", err)
	}
	if len(result2.Readings) != 2 {
		t.Errorf("page 2: expected 2 readings, got %d", len(result2.Readings))
	}
}
