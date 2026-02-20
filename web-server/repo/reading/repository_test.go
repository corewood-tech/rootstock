package reading

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"

	"rootstock/web-server/config"
	sqlmigrate "rootstock/web-server/repo/sql/migrate"
)

func setupTest(t *testing.T) (Repository, *pgxpool.Pool) {
	t.Helper()
	cfg := config.PostgresConfig{
		Host:     "app-postgres",
		Port:     5432,
		User:     "rootstock",
		Password: "rootstock",
		DBName:   "rootstock",
		SSLMode:  "disable",
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

	repo := NewRepository(pool)
	t.Cleanup(func() {
		repo.Shutdown()
		pool.Close()
	})

	return repo, pool
}

// createFixtures inserts a device and campaign for FK constraints.
func createFixtures(t *testing.T, pool *pgxpool.Pool) (deviceID, campaignID string) {
	t.Helper()
	ctx := context.Background()

	deviceID = ulid.Make().String()
	_, err := pool.Exec(ctx,
		`INSERT INTO devices (id, owner_id, class, firmware_version, tier, sensors, status)
		 VALUES ($1, 'user-1', 'sensor', '1.0.0', 1, '{temp}', 'active')`, deviceID)
	if err != nil {
		t.Fatalf("insert device: %v", err)
	}

	campaignID = ulid.Make().String()
	_, err = pool.Exec(ctx,
		`INSERT INTO campaigns (id, org_id, created_by) VALUES ($1, 'org-1', 'user-1')`, campaignID)
	if err != nil {
		t.Fatalf("insert campaign: %v", err)
	}
	return
}

func TestPersistAndQuery(t *testing.T) {
	repo, pool := setupTest(t)
	ctx := context.Background()
	deviceID, campaignID := createFixtures(t, pool)

	ts := time.Now().UTC().Truncate(time.Microsecond)
	created, err := repo.Persist(ctx, PersistReadingInput{
		DeviceID:        deviceID,
		CampaignID:      campaignID,
		Value:           23.5,
		Timestamp:       ts,
		Geolocation:     `{"type":"Point","coordinates":[-73.98,40.74]}`,
		FirmwareVersion: "1.0.0",
		CertSerial:      "serial-001",
	})
	if err != nil {
		t.Fatalf("Persist(): %v", err)
	}
	if created.ID == "" {
		t.Fatal("Persist() returned empty ID")
	}
	if created.Status != "accepted" {
		t.Errorf("status = %q, want %q", created.Status, "accepted")
	}

	readings, err := repo.Query(ctx, QueryReadingsInput{CampaignID: campaignID})
	if err != nil {
		t.Fatalf("Query(): %v", err)
	}
	if len(readings) != 1 {
		t.Fatalf("Query() = %d, want 1", len(readings))
	}
	if readings[0].Value != 23.5 {
		t.Errorf("value = %f, want 23.5", readings[0].Value)
	}
}

func TestQuarantine(t *testing.T) {
	repo, pool := setupTest(t)
	ctx := context.Background()
	deviceID, campaignID := createFixtures(t, pool)

	rd, _ := repo.Persist(ctx, PersistReadingInput{
		DeviceID: deviceID, CampaignID: campaignID, Value: 999.0,
		Timestamp: time.Now().UTC(), FirmwareVersion: "1.0.0", CertSerial: "s1",
	})

	if err := repo.Quarantine(ctx, rd.ID, "outlier"); err != nil {
		t.Fatalf("Quarantine(): %v", err)
	}

	readings, _ := repo.Query(ctx, QueryReadingsInput{CampaignID: campaignID, Status: "quarantined"})
	if len(readings) != 1 {
		t.Fatalf("quarantined count = %d, want 1", len(readings))
	}
	if readings[0].QuarantineReason == nil || *readings[0].QuarantineReason != "outlier" {
		t.Errorf("quarantine_reason = %v, want 'outlier'", readings[0].QuarantineReason)
	}
}

func TestQuarantineByWindow(t *testing.T) {
	repo, pool := setupTest(t)
	ctx := context.Background()
	deviceID, campaignID := createFixtures(t, pool)

	base := time.Now().UTC().Add(-1 * time.Hour)
	// Insert 3 readings
	for i := 0; i < 3; i++ {
		repo.Persist(ctx, PersistReadingInput{
			DeviceID: deviceID, CampaignID: campaignID, Value: float64(i),
			Timestamp: base.Add(time.Duration(i*10) * time.Minute), FirmwareVersion: "1.0.0", CertSerial: "s1",
		})
	}

	// Quarantine readings in the first 15 minutes
	affected, err := repo.QuarantineByWindow(ctx, QuarantineByWindowInput{
		DeviceIDs: []string{deviceID},
		Since:     base,
		Until:     base.Add(15 * time.Minute),
		Reason:    "vulnerability-CVE-2025-0001",
	})
	if err != nil {
		t.Fatalf("QuarantineByWindow(): %v", err)
	}
	if affected != 2 {
		t.Errorf("affected = %d, want 2", affected)
	}
}

func TestGetCampaignQualityEmpty(t *testing.T) {
	repo, pool := setupTest(t)
	ctx := context.Background()
	_, campaignID := createFixtures(t, pool)

	quality, err := repo.GetCampaignQuality(ctx, campaignID)
	if err != nil {
		t.Fatalf("GetCampaignQuality(): %v", err)
	}
	if quality.AcceptedCount != 0 {
		t.Errorf("accepted = %d, want 0", quality.AcceptedCount)
	}
	if quality.QuarantineCount != 0 {
		t.Errorf("quarantined = %d, want 0", quality.QuarantineCount)
	}
}

func TestGetCampaignQuality(t *testing.T) {
	repo, pool := setupTest(t)
	ctx := context.Background()
	deviceID, campaignID := createFixtures(t, pool)

	ts := time.Now().UTC()
	repo.Persist(ctx, PersistReadingInput{
		DeviceID: deviceID, CampaignID: campaignID, Value: 1.0,
		Timestamp: ts, FirmwareVersion: "1.0.0", CertSerial: "s1",
	})
	rd, _ := repo.Persist(ctx, PersistReadingInput{
		DeviceID: deviceID, CampaignID: campaignID, Value: 2.0,
		Timestamp: ts, FirmwareVersion: "1.0.0", CertSerial: "s1",
	})
	repo.Quarantine(ctx, rd.ID, "outlier")

	quality, err := repo.GetCampaignQuality(ctx, campaignID)
	if err != nil {
		t.Fatalf("GetCampaignQuality(): %v", err)
	}
	if quality.AcceptedCount != 1 {
		t.Errorf("accepted = %d, want 1", quality.AcceptedCount)
	}
	if quality.QuarantineCount != 1 {
		t.Errorf("quarantined = %d, want 1", quality.QuarantineCount)
	}
}
