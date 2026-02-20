package score

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"rootstock/web-server/config"
	deviceops "rootstock/web-server/ops/device"
	readingops "rootstock/web-server/ops/reading"
	scoreops "rootstock/web-server/ops/score"
	devicerepo "rootstock/web-server/repo/device"
	readingrepo "rootstock/web-server/repo/reading"
	scorerepo "rootstock/web-server/repo/score"
	sqlmigrate "rootstock/web-server/repo/sql/migrate"
)

func setupRefreshScoreTest(t *testing.T) (*RefreshScitizenScoreFlow, *pgxpool.Pool) {
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
	pool.Exec(ctx, "TRUNCATE devices, readings, scores, badges, sweepstakes_entries, device_campaigns CASCADE")

	dRepo := devicerepo.NewRepository(pool)
	rRepo := readingrepo.NewRepository(pool)
	sRepo := scorerepo.NewRepository(pool)

	dOps := deviceops.NewOps(dRepo)
	rOps := readingops.NewOps(rRepo)
	sOps := scoreops.NewOps(sRepo)

	flow := NewRefreshScitizenScoreFlow(dOps, rOps, sOps)

	t.Cleanup(func() {
		dRepo.Shutdown()
		rRepo.Shutdown()
		sRepo.Shutdown()
		pool.Close()
	})

	return flow, pool
}

func TestRefreshScitizenScore(t *testing.T) {
	flow, pool := setupRefreshScoreTest(t)
	ctx := context.Background()

	// Insert a device owned by scitizen sci-1
	pool.Exec(ctx, `INSERT INTO devices (id, owner_id, status, class, firmware_version, tier, sensors)
		VALUES ('dev-1', 'sci-1', 'active', 'tier1', '1.0.0', 1, '{"temperature"}')`)

	// Insert an accepted reading for the device
	pool.Exec(ctx, `INSERT INTO readings (id, device_id, campaign_id, value, timestamp, firmware_version, cert_serial, status)
		VALUES ('r-1', 'dev-1', 'camp-1', 22.5, $1, '1.0.0', 'serial-1', 'accepted')`,
		time.Now())

	result, err := flow.Run(ctx, RefreshScitizenScoreInput{DeviceID: "dev-1"})
	if err != nil {
		t.Fatalf("Run(): %v", err)
	}

	if result.Score.Volume != 1 {
		t.Errorf("volume = %d, want 1", result.Score.Volume)
	}
	if result.Score.ScitizenID != "sci-1" {
		t.Errorf("scitizen_id = %s, want sci-1", result.Score.ScitizenID)
	}

	// Should have awarded "first-contribution" badge
	found := false
	for _, b := range result.BadgesAwarded {
		if b == "first-contribution" {
			found = true
		}
	}
	if !found {
		t.Error("expected first-contribution badge to be awarded")
	}
}

func TestRefreshScitizenScoreNoReadings(t *testing.T) {
	flow, pool := setupRefreshScoreTest(t)
	ctx := context.Background()

	// Device exists but no readings
	pool.Exec(ctx, `INSERT INTO devices (id, owner_id, status, class, firmware_version, tier, sensors)
		VALUES ('dev-2', 'sci-2', 'active', 'tier1', '1.0.0', 1, '{"temperature"}')`)

	result, err := flow.Run(ctx, RefreshScitizenScoreInput{DeviceID: "dev-2"})
	if err != nil {
		t.Fatalf("Run(): %v", err)
	}

	if result.Score.Volume != 0 {
		t.Errorf("volume = %d, want 0", result.Score.Volume)
	}
}
