package reading

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	campaignops "rootstock/web-server/ops/campaign"
	readingops "rootstock/web-server/ops/reading"
	"rootstock/web-server/config"
	campaignrepo "rootstock/web-server/repo/campaign"
	readingrepo "rootstock/web-server/repo/reading"
	sqlmigrate "rootstock/web-server/repo/sql/migrate"
)

func setupIngestTest(t *testing.T) (*IngestReadingFlow, *pgxpool.Pool) {
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

	cRepo := campaignrepo.NewRepository(pool)
	rRepo := readingrepo.NewRepository(pool)

	cOps := campaignops.NewOps(cRepo)
	rOps := readingops.NewOps(rRepo)

	flow := NewIngestReadingFlow(cOps, rOps)

	t.Cleanup(func() {
		cRepo.Shutdown()
		rRepo.Shutdown()
		pool.Close()
	})

	return flow, pool
}

func TestIngestValidReading(t *testing.T) {
	flow, pool := setupIngestTest(t)
	ctx := context.Background()

	now := time.Now().UTC()
	start := now.Add(-1 * time.Hour)
	end := now.Add(1 * time.Hour)
	min := 0.0
	max := 100.0

	cRepo := campaignrepo.NewRepository(pool)
	defer cRepo.Shutdown()
	campaign, err := cRepo.Create(ctx, campaignrepo.CreateCampaignInput{
		OrgID:       "org-1",
		CreatedBy:   "user-1",
		WindowStart: &start,
		WindowEnd:   &end,
		Parameters:  []campaignrepo.ParameterInput{{Name: "temp", Unit: "celsius", MinRange: &min, MaxRange: &max}},
	})
	if err != nil {
		t.Fatalf("create campaign: %v", err)
	}

	var deviceID string
	pool.QueryRow(ctx,
		`INSERT INTO devices (owner_id, class, firmware_version, tier, sensors, status)
		 VALUES ('user-1', 'sensor', '1.0.0', 1, '{temp}', 'active') RETURNING id`,
	).Scan(&deviceID)

	rd, err := flow.Run(ctx, IngestReadingInput{
		DeviceID:        deviceID,
		CampaignID:      campaign.ID,
		Value:           23.5,
		Timestamp:       now,
		FirmwareVersion: "1.0.0",
		CertSerial:      "serial-1",
	})
	if err != nil {
		t.Fatalf("Run(): %v", err)
	}
	if rd.Status != "accepted" {
		t.Errorf("status = %q, want accepted", rd.Status)
	}
}

func TestIngestOutOfRangeReading(t *testing.T) {
	flow, pool := setupIngestTest(t)
	ctx := context.Background()

	now := time.Now().UTC()
	start := now.Add(-1 * time.Hour)
	end := now.Add(1 * time.Hour)
	min := 0.0
	max := 50.0

	cRepo := campaignrepo.NewRepository(pool)
	defer cRepo.Shutdown()
	campaign, _ := cRepo.Create(ctx, campaignrepo.CreateCampaignInput{
		OrgID:       "org-1",
		CreatedBy:   "user-1",
		WindowStart: &start,
		WindowEnd:   &end,
		Parameters:  []campaignrepo.ParameterInput{{Name: "temp", Unit: "celsius", MinRange: &min, MaxRange: &max}},
	})

	var deviceID string
	pool.QueryRow(ctx,
		`INSERT INTO devices (owner_id, class, firmware_version, tier, sensors, status)
		 VALUES ('user-1', 'sensor', '1.0.0', 1, '{temp}', 'active') RETURNING id`,
	).Scan(&deviceID)

	rd, err := flow.Run(ctx, IngestReadingInput{
		DeviceID:        deviceID,
		CampaignID:      campaign.ID,
		Value:           999.0,
		Timestamp:       now,
		FirmwareVersion: "1.0.0",
		CertSerial:      "serial-1",
	})
	if err != nil {
		t.Fatalf("Run(): %v", err)
	}
	if rd.Status != "quarantined" {
		t.Errorf("status = %q, want quarantined", rd.Status)
	}
	if rd.QuarantineReason == nil {
		t.Error("quarantine_reason should not be nil")
	}
}

func TestIngestOutsideWindowReading(t *testing.T) {
	flow, pool := setupIngestTest(t)
	ctx := context.Background()

	now := time.Now().UTC()
	start := now.Add(1 * time.Hour)
	end := now.Add(2 * time.Hour)

	cRepo := campaignrepo.NewRepository(pool)
	defer cRepo.Shutdown()
	campaign, _ := cRepo.Create(ctx, campaignrepo.CreateCampaignInput{
		OrgID:       "org-1",
		CreatedBy:   "user-1",
		WindowStart: &start,
		WindowEnd:   &end,
	})

	var deviceID string
	pool.QueryRow(ctx,
		`INSERT INTO devices (owner_id, class, firmware_version, tier, sensors, status)
		 VALUES ('user-1', 'sensor', '1.0.0', 1, '{temp}', 'active') RETURNING id`,
	).Scan(&deviceID)

	rd, err := flow.Run(ctx, IngestReadingInput{
		DeviceID:        deviceID,
		CampaignID:      campaign.ID,
		Value:           23.5,
		Timestamp:       now,
		FirmwareVersion: "1.0.0",
		CertSerial:      "serial-1",
	})
	if err != nil {
		t.Fatalf("Run(): %v", err)
	}
	if rd.Status != "quarantined" {
		t.Errorf("status = %q, want quarantined", rd.Status)
	}
}
