package campaign

import (
	"context"
	"fmt"
	"testing"
	"time"

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
	pool.Exec(ctx, "TRUNCATE campaigns CASCADE")

	repo := NewRepository(pool)
	t.Cleanup(func() {
		repo.Shutdown()
		pool.Close()
	})

	return repo, pool
}

func TestCreateAndList(t *testing.T) {
	repo, _ := setupTest(t)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Microsecond)
	end := now.Add(30 * 24 * time.Hour)

	minRange := 0.0
	maxRange := 100.0
	precision := 2

	created, err := repo.Create(ctx, CreateCampaignInput{
		OrgID:       "org-1",
		CreatedBy:   "user-1",
		WindowStart: &now,
		WindowEnd:   &end,
		Parameters: []ParameterInput{
			{Name: "temperature", Unit: "celsius", MinRange: &minRange, MaxRange: &maxRange, Precision: &precision},
		},
		Regions: []RegionInput{
			{GeoJSON: `{"type":"Point","coordinates":[-73.9857,40.7484]}`},
		},
		Eligibility: []EligibilityInput{
			{DeviceClass: "weather-station", Tier: 1, RequiredSensors: []string{"temp"}, FirmwareMin: "1.0.0"},
		},
	})
	if err != nil {
		t.Fatalf("Create(): %v", err)
	}
	if created.ID == "" {
		t.Fatal("Create() returned empty ID")
	}
	if created.Status != "draft" {
		t.Errorf("Create() status = %q, want %q", created.Status, "draft")
	}
	if created.OrgID != "org-1" {
		t.Errorf("Create() org_id = %q, want %q", created.OrgID, "org-1")
	}

	campaigns, err := repo.List(ctx, ListCampaignsInput{OrgID: "org-1"})
	if err != nil {
		t.Fatalf("List(): %v", err)
	}
	if len(campaigns) != 1 {
		t.Fatalf("List() returned %d campaigns, want 1", len(campaigns))
	}
	if campaigns[0].ID != created.ID {
		t.Errorf("List() ID = %q, want %q", campaigns[0].ID, created.ID)
	}
}

func TestPublish(t *testing.T) {
	repo, _ := setupTest(t)
	ctx := context.Background()

	created, err := repo.Create(ctx, CreateCampaignInput{
		OrgID:     "org-1",
		CreatedBy: "user-1",
	})
	if err != nil {
		t.Fatalf("Create(): %v", err)
	}

	if err := repo.Publish(ctx, created.ID); err != nil {
		t.Fatalf("Publish(): %v", err)
	}

	campaigns, err := repo.List(ctx, ListCampaignsInput{Status: "published"})
	if err != nil {
		t.Fatalf("List(): %v", err)
	}
	if len(campaigns) != 1 {
		t.Fatalf("List(published) returned %d, want 1", len(campaigns))
	}

	// Publishing again should fail (already published, not draft)
	err = repo.Publish(ctx, created.ID)
	if err == nil {
		t.Error("Publish() second call should fail, got nil")
	}
}

func TestGetRules(t *testing.T) {
	repo, _ := setupTest(t)
	ctx := context.Background()

	minR := 10.0
	maxR := 50.0
	prec := 1

	created, err := repo.Create(ctx, CreateCampaignInput{
		OrgID:     "org-1",
		CreatedBy: "user-1",
		Parameters: []ParameterInput{
			{Name: "humidity", Unit: "percent", MinRange: &minR, MaxRange: &maxR, Precision: &prec},
		},
		Regions: []RegionInput{
			{GeoJSON: `{"type":"Polygon","coordinates":[[[-74,40],[-73,40],[-73,41],[-74,41],[-74,40]]]}`},
		},
	})
	if err != nil {
		t.Fatalf("Create(): %v", err)
	}

	rules, err := repo.GetRules(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetRules(): %v", err)
	}
	if len(rules.Parameters) != 1 {
		t.Fatalf("GetRules() parameters = %d, want 1", len(rules.Parameters))
	}
	if rules.Parameters[0].Name != "humidity" {
		t.Errorf("parameter name = %q, want %q", rules.Parameters[0].Name, "humidity")
	}
	if len(rules.Regions) != 1 {
		t.Fatalf("GetRules() regions = %d, want 1", len(rules.Regions))
	}
}

func TestGetEligibility(t *testing.T) {
	repo, _ := setupTest(t)
	ctx := context.Background()

	created, err := repo.Create(ctx, CreateCampaignInput{
		OrgID:     "org-1",
		CreatedBy: "user-1",
		Eligibility: []EligibilityInput{
			{DeviceClass: "air-quality", Tier: 2, RequiredSensors: []string{"pm25", "pm10"}, FirmwareMin: "2.0.0"},
		},
	})
	if err != nil {
		t.Fatalf("Create(): %v", err)
	}

	elig, err := repo.GetEligibility(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetEligibility(): %v", err)
	}
	if len(elig) != 1 {
		t.Fatalf("GetEligibility() = %d, want 1", len(elig))
	}
	if elig[0].DeviceClass != "air-quality" {
		t.Errorf("device_class = %q, want %q", elig[0].DeviceClass, "air-quality")
	}
	if len(elig[0].RequiredSensors) != 2 {
		t.Errorf("required_sensors = %d, want 2", len(elig[0].RequiredSensors))
	}
}

func TestListFilterByStatus(t *testing.T) {
	repo, _ := setupTest(t)
	ctx := context.Background()

	// Create two campaigns, publish one
	c1, _ := repo.Create(ctx, CreateCampaignInput{OrgID: "org-1", CreatedBy: "user-1"})
	repo.Create(ctx, CreateCampaignInput{OrgID: "org-1", CreatedBy: "user-1"})
	repo.Publish(ctx, c1.ID)

	drafts, err := repo.List(ctx, ListCampaignsInput{Status: "draft"})
	if err != nil {
		t.Fatalf("List(draft): %v", err)
	}
	if len(drafts) != 1 {
		t.Errorf("List(draft) = %d, want 1", len(drafts))
	}

	published, err := repo.List(ctx, ListCampaignsInput{Status: "published"})
	if err != nil {
		t.Fatalf("List(published): %v", err)
	}
	if len(published) != 1 {
		t.Errorf("List(published) = %d, want 1", len(published))
	}
}
