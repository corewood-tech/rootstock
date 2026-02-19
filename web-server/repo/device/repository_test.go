package device

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

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
	pool.Exec(ctx, "TRUNCATE devices CASCADE")

	repo := NewRepository(pool)
	t.Cleanup(func() {
		repo.Shutdown()
		pool.Close()
	})

	return repo, pool
}

func TestCreateAndGet(t *testing.T) {
	repo, _ := setupTest(t)
	ctx := context.Background()

	created, err := repo.Create(ctx, CreateDeviceInput{
		OwnerID:         "user-1",
		Class:           "weather-station",
		FirmwareVersion: "1.2.0",
		Tier:            1,
		Sensors:         []string{"temp", "humidity"},
	})
	if err != nil {
		t.Fatalf("Create(): %v", err)
	}
	if created.ID == "" {
		t.Fatal("Create() returned empty ID")
	}
	if created.Status != "pending" {
		t.Errorf("status = %q, want %q", created.Status, "pending")
	}

	got, err := repo.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get(): %v", err)
	}
	if got.Class != "weather-station" {
		t.Errorf("class = %q, want %q", got.Class, "weather-station")
	}
	if len(got.Sensors) != 2 {
		t.Errorf("sensors = %d, want 2", len(got.Sensors))
	}
}

func TestUpdateStatus(t *testing.T) {
	repo, _ := setupTest(t)
	ctx := context.Background()

	d, _ := repo.Create(ctx, CreateDeviceInput{
		OwnerID: "user-1", Class: "sensor", FirmwareVersion: "1.0.0", Tier: 1, Sensors: []string{"temp"},
	})

	if err := repo.UpdateStatus(ctx, d.ID, "active"); err != nil {
		t.Fatalf("UpdateStatus(): %v", err)
	}

	got, _ := repo.Get(ctx, d.ID)
	if got.Status != "active" {
		t.Errorf("status = %q, want %q", got.Status, "active")
	}
}

func TestGetCapabilities(t *testing.T) {
	repo, _ := setupTest(t)
	ctx := context.Background()

	d, _ := repo.Create(ctx, CreateDeviceInput{
		OwnerID: "user-1", Class: "air-quality", FirmwareVersion: "2.1.0", Tier: 2, Sensors: []string{"pm25", "pm10"},
	})

	caps, err := repo.GetCapabilities(ctx, d.ID)
	if err != nil {
		t.Fatalf("GetCapabilities(): %v", err)
	}
	if caps.Class != "air-quality" {
		t.Errorf("class = %q, want %q", caps.Class, "air-quality")
	}
	if caps.Tier != 2 {
		t.Errorf("tier = %d, want 2", caps.Tier)
	}
}

func TestEnrollmentCodeLifecycle(t *testing.T) {
	repo, _ := setupTest(t)
	ctx := context.Background()

	d, _ := repo.Create(ctx, CreateDeviceInput{
		OwnerID: "user-1", Class: "sensor", FirmwareVersion: "1.0.0", Tier: 1, Sensors: []string{"temp"},
	})

	code, err := repo.GenerateEnrollmentCode(ctx, GenerateCodeInput{
		DeviceID: d.ID,
		Code:     "ABC123",
		TTL:      900,
	})
	if err != nil {
		t.Fatalf("GenerateEnrollmentCode(): %v", err)
	}
	if code.Code != "ABC123" {
		t.Errorf("code = %q, want %q", code.Code, "ABC123")
	}
	if code.Used {
		t.Error("new code should not be used")
	}

	redeemed, err := repo.RedeemEnrollmentCode(ctx, "ABC123")
	if err != nil {
		t.Fatalf("RedeemEnrollmentCode(): %v", err)
	}
	if !redeemed.Used {
		t.Error("redeemed code should be used")
	}

	// Redeeming again should fail
	_, err = repo.RedeemEnrollmentCode(ctx, "ABC123")
	if err == nil {
		t.Error("second redeem should fail")
	}
}

func TestQueryByClass(t *testing.T) {
	repo, _ := setupTest(t)
	ctx := context.Background()

	repo.Create(ctx, CreateDeviceInput{
		OwnerID: "user-1", Class: "weather-station", FirmwareVersion: "1.0.0", Tier: 1, Sensors: []string{"temp"},
	})
	repo.Create(ctx, CreateDeviceInput{
		OwnerID: "user-1", Class: "weather-station", FirmwareVersion: "2.0.0", Tier: 1, Sensors: []string{"temp"},
	})
	repo.Create(ctx, CreateDeviceInput{
		OwnerID: "user-1", Class: "air-quality", FirmwareVersion: "1.0.0", Tier: 2, Sensors: []string{"pm25"},
	})

	devices, err := repo.QueryByClass(ctx, QueryByClassInput{Class: "weather-station"})
	if err != nil {
		t.Fatalf("QueryByClass(): %v", err)
	}
	if len(devices) != 2 {
		t.Errorf("QueryByClass(weather-station) = %d, want 2", len(devices))
	}

	devices, err = repo.QueryByClass(ctx, QueryByClassInput{
		Class:          "weather-station",
		FirmwareMinGte: "1.5.0",
	})
	if err != nil {
		t.Fatalf("QueryByClass(gte): %v", err)
	}
	if len(devices) != 1 {
		t.Errorf("QueryByClass(gte 1.5.0) = %d, want 1", len(devices))
	}
}

func TestEnrollInCampaign(t *testing.T) {
	repo, pool := setupTest(t)
	ctx := context.Background()

	d, _ := repo.Create(ctx, CreateDeviceInput{
		OwnerID: "user-1", Class: "sensor", FirmwareVersion: "1.0.0", Tier: 1, Sensors: []string{"temp"},
	})

	// Create a campaign directly in DB for the FK
	var campaignID string
	err := pool.QueryRow(ctx,
		`INSERT INTO campaigns (org_id, created_by) VALUES ('org-1', 'user-1') RETURNING id`,
	).Scan(&campaignID)
	if err != nil {
		t.Fatalf("insert campaign: %v", err)
	}

	if err := repo.EnrollInCampaign(ctx, d.ID, campaignID); err != nil {
		t.Fatalf("EnrollInCampaign(): %v", err)
	}

	// Enrolling again should fail (duplicate PK)
	err = repo.EnrollInCampaign(ctx, d.ID, campaignID)
	if err == nil {
		t.Error("duplicate enrollment should fail")
	}
}
