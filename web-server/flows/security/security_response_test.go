package security

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"

	"rootstock/web-server/config"
	deviceops "rootstock/web-server/ops/device"
	notificationops "rootstock/web-server/ops/notification"
	readingops "rootstock/web-server/ops/reading"
	devicerepo "rootstock/web-server/repo/device"
	notificationrepo "rootstock/web-server/repo/notification"
	readingrepo "rootstock/web-server/repo/reading"
	sqlmigrate "rootstock/web-server/repo/sql/migrate"
)

func setupSecurityResponseTest(t *testing.T) (*SecurityResponseFlow, *pgxpool.Pool) {
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
	pool.Exec(ctx, "TRUNCATE devices, readings, scores, badges, sweepstakes_entries, device_campaigns, campaigns CASCADE")

	dRepo := devicerepo.NewRepository(pool)
	rRepo := readingrepo.NewRepository(pool)
	nRepo := notificationrepo.NewRepository()

	dOps := deviceops.NewOps(dRepo)
	rOps := readingops.NewOps(rRepo)
	nOps := notificationops.NewOps(nRepo)

	flow := NewSecurityResponseFlow(dOps, rOps, nOps)

	t.Cleanup(func() {
		dRepo.Shutdown()
		rRepo.Shutdown()
		nRepo.Shutdown()
		pool.Close()
	})

	return flow, pool
}

func insertDevice(t *testing.T, pool *pgxpool.Pool, ownerID, class, firmware string) string {
	t.Helper()
	id := ulid.Make().String()
	_, err := pool.Exec(context.Background(),
		`INSERT INTO devices (id, owner_id, status, class, firmware_version, tier, sensors)
		 VALUES ($1, $2, 'active', $3, $4, 1, '{"temperature"}')`,
		id, ownerID, class, firmware,
	)
	if err != nil {
		t.Fatalf("insert device: %v", err)
	}
	return id
}

func insertCampaign(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()
	id := ulid.Make().String()
	_, err := pool.Exec(context.Background(),
		`INSERT INTO campaigns (id, org_id, created_by) VALUES ($1, 'org-1', 'user-1')`, id,
	)
	if err != nil {
		t.Fatalf("insert campaign: %v", err)
	}
	return id
}

func insertReading(t *testing.T, pool *pgxpool.Pool, deviceID, campaignID string, ts time.Time) string {
	t.Helper()
	id := ulid.Make().String()
	_, err := pool.Exec(context.Background(),
		`INSERT INTO readings (id, device_id, campaign_id, value, timestamp, firmware_version, cert_serial, status)
		 VALUES ($1, $2, $3, 22.5, $4, '1.0.0', 'serial-1', 'accepted')`,
		id, deviceID, campaignID, ts,
	)
	if err != nil {
		t.Fatalf("insert reading: %v", err)
	}
	return id
}

func TestSecurityResponse(t *testing.T) {
	flow, pool := setupSecurityResponseTest(t)
	ctx := context.Background()

	dev1 := insertDevice(t, pool, "sci-1", "tier1", "1.0.0")
	dev2 := insertDevice(t, pool, "sci-2", "tier1", "1.0.0")
	campID := insertCampaign(t, pool)

	windowStart := time.Now().Add(-24 * time.Hour)
	windowEnd := time.Now()

	r1 := insertReading(t, pool, dev1, campID, windowStart.Add(1*time.Hour))
	insertReading(t, pool, dev2, campID, windowStart.Add(2*time.Hour))

	result, err := flow.Run(ctx, SecurityResponseInput{
		Class:       "tier1",
		FirmwareMin: "1.0.0",
		FirmwareMax: "1.0.0",
		WindowStart: windowStart,
		WindowEnd:   windowEnd,
		Reason:      "CVE-2025-0001: sensor firmware vulnerability",
	})
	if err != nil {
		t.Fatalf("Run(): %v", err)
	}

	if result.SuspendedCount != 2 {
		t.Errorf("SuspendedCount = %d, want 2", result.SuspendedCount)
	}
	if result.QuarantinedReadings != 2 {
		t.Errorf("QuarantinedReadings = %d, want 2", result.QuarantinedReadings)
	}
	if result.NotifiedScitizens != 2 {
		t.Errorf("NotifiedScitizens = %d, want 2", result.NotifiedScitizens)
	}

	// Verify devices are suspended in DB.
	var status string
	pool.QueryRow(ctx, "SELECT status FROM devices WHERE id = $1", dev1).Scan(&status)
	if status != "suspended" {
		t.Errorf("dev1 status = %q, want suspended", status)
	}
	pool.QueryRow(ctx, "SELECT status FROM devices WHERE id = $1", dev2).Scan(&status)
	if status != "suspended" {
		t.Errorf("dev2 status = %q, want suspended", status)
	}

	// Verify readings are quarantined.
	var readingStatus string
	pool.QueryRow(ctx, "SELECT status FROM readings WHERE id = $1", r1).Scan(&readingStatus)
	if readingStatus != "quarantined" {
		t.Errorf("r1 status = %q, want quarantined", readingStatus)
	}
}

func TestSecurityResponseNoMatchingDevices(t *testing.T) {
	flow, _ := setupSecurityResponseTest(t)
	ctx := context.Background()

	result, err := flow.Run(ctx, SecurityResponseInput{
		Class:       "tier1",
		FirmwareMin: "1.0.0",
		FirmwareMax: "1.0.0",
		WindowStart: time.Now().Add(-24 * time.Hour),
		WindowEnd:   time.Now(),
		Reason:      "test vulnerability",
	})
	if err != nil {
		t.Fatalf("Run(): %v", err)
	}

	if result.SuspendedCount != 0 {
		t.Errorf("SuspendedCount = %d, want 0", result.SuspendedCount)
	}
	if result.QuarantinedReadings != 0 {
		t.Errorf("QuarantinedReadings = %d, want 0", result.QuarantinedReadings)
	}
	if result.NotifiedScitizens != 0 {
		t.Errorf("NotifiedScitizens = %d, want 0", result.NotifiedScitizens)
	}
}

func TestSecurityResponseDeduplicatesOwners(t *testing.T) {
	flow, pool := setupSecurityResponseTest(t)
	ctx := context.Background()

	// Two devices owned by the same scitizen
	insertDevice(t, pool, "sci-4", "tier1", "1.0.0")
	insertDevice(t, pool, "sci-4", "tier1", "1.0.0")

	result, err := flow.Run(ctx, SecurityResponseInput{
		Class:       "tier1",
		FirmwareMin: "1.0.0",
		FirmwareMax: "1.0.0",
		WindowStart: time.Now().Add(-24 * time.Hour),
		WindowEnd:   time.Now(),
		Reason:      "test vulnerability",
	})
	if err != nil {
		t.Fatalf("Run(): %v", err)
	}

	if result.SuspendedCount != 2 {
		t.Errorf("SuspendedCount = %d, want 2", result.SuspendedCount)
	}
	// Only 1 notification even though 2 devices — same owner
	if result.NotifiedScitizens != 1 {
		t.Errorf("NotifiedScitizens = %d, want 1", result.NotifiedScitizens)
	}
}
