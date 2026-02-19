package reading

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"rootstock/web-server/config"
	readingrepo "rootstock/web-server/repo/reading"
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
	pool.Exec(ctx, "TRUNCATE readings, devices, campaigns CASCADE")

	repo := readingrepo.NewRepository(pool)
	ops := NewOps(repo)
	t.Cleanup(func() {
		repo.Shutdown()
		pool.Close()
	})
	return ops, pool
}

func createFixtures(t *testing.T, pool *pgxpool.Pool) (deviceID, campaignID string) {
	t.Helper()
	ctx := context.Background()
	pool.QueryRow(ctx,
		`INSERT INTO devices (owner_id, class, firmware_version, tier, sensors, status)
		 VALUES ('user-1', 'sensor', '1.0.0', 1, '{temp}', 'active') RETURNING id`,
	).Scan(&deviceID)
	pool.QueryRow(ctx,
		`INSERT INTO campaigns (org_id, created_by) VALUES ('org-1', 'user-1') RETURNING id`,
	).Scan(&campaignID)
	return
}

func TestPersistAndQuery(t *testing.T) {
	ops, pool := setupTest(t)
	ctx := context.Background()
	deviceID, campaignID := createFixtures(t, pool)

	rd, err := ops.PersistReading(ctx, PersistReadingInput{
		DeviceID: deviceID, CampaignID: campaignID, Value: 23.5,
		Timestamp: time.Now().UTC(), FirmwareVersion: "1.0.0", CertSerial: "s1",
	})
	if err != nil {
		t.Fatalf("PersistReading(): %v", err)
	}
	if rd.Status != "accepted" {
		t.Errorf("status = %q, want accepted", rd.Status)
	}

	readings, err := ops.QueryReadings(ctx, QueryReadingsInput{CampaignID: campaignID})
	if err != nil {
		t.Fatalf("QueryReadings(): %v", err)
	}
	if len(readings) != 1 {
		t.Errorf("count = %d, want 1", len(readings))
	}
}

func TestQuarantineReading(t *testing.T) {
	ops, pool := setupTest(t)
	ctx := context.Background()
	deviceID, campaignID := createFixtures(t, pool)

	rd, _ := ops.PersistReading(ctx, PersistReadingInput{
		DeviceID: deviceID, CampaignID: campaignID, Value: 999.0,
		Timestamp: time.Now().UTC(), FirmwareVersion: "1.0.0", CertSerial: "s1",
	})

	if err := ops.QuarantineReading(ctx, rd.ID, "outlier"); err != nil {
		t.Fatalf("QuarantineReading(): %v", err)
	}

	quality, err := ops.GetCampaignQuality(ctx, campaignID)
	if err != nil {
		t.Fatalf("GetCampaignQuality(): %v", err)
	}
	if quality.QuarantineCount != 1 {
		t.Errorf("quarantined = %d, want 1", quality.QuarantineCount)
	}
}
