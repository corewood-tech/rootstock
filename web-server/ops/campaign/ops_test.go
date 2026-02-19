package campaign

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"rootstock/web-server/config"
	campaignrepo "rootstock/web-server/repo/campaign"
	sqlmigrate "rootstock/web-server/repo/sql/migrate"
)

func setupTest(t *testing.T) *Ops {
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

	repo := campaignrepo.NewRepository(pool)
	ops := NewOps(repo)

	t.Cleanup(func() {
		repo.Shutdown()
		pool.Close()
	})

	return ops
}

func TestCreateAndPublish(t *testing.T) {
	ops := setupTest(t)
	ctx := context.Background()

	c, err := ops.CreateCampaign(ctx, CreateCampaignInput{
		OrgID:     "org-1",
		CreatedBy: "user-1",
		Parameters: []ParameterInput{
			{Name: "temperature", Unit: "celsius"},
		},
	})
	if err != nil {
		t.Fatalf("CreateCampaign(): %v", err)
	}
	if c.Status != "draft" {
		t.Errorf("status = %q, want draft", c.Status)
	}

	if err := ops.PublishCampaign(ctx, c.ID); err != nil {
		t.Fatalf("PublishCampaign(): %v", err)
	}

	campaigns, err := ops.ListCampaigns(ctx, ListCampaignsInput{Status: "published"})
	if err != nil {
		t.Fatalf("ListCampaigns(): %v", err)
	}
	if len(campaigns) != 1 {
		t.Errorf("published count = %d, want 1", len(campaigns))
	}
}

func TestGetRulesAndEligibility(t *testing.T) {
	ops := setupTest(t)
	ctx := context.Background()

	min := 0.0
	max := 100.0
	c, err := ops.CreateCampaign(ctx, CreateCampaignInput{
		OrgID:     "org-1",
		CreatedBy: "user-1",
		Parameters: []ParameterInput{
			{Name: "humidity", Unit: "percent", MinRange: &min, MaxRange: &max},
		},
		Eligibility: []EligibilityInput{
			{DeviceClass: "weather-station", Tier: 1, RequiredSensors: []string{"humidity"}, FirmwareMin: "1.0.0"},
		},
	})
	if err != nil {
		t.Fatalf("CreateCampaign(): %v", err)
	}

	rules, err := ops.GetCampaignRules(ctx, c.ID)
	if err != nil {
		t.Fatalf("GetCampaignRules(): %v", err)
	}
	if len(rules.Parameters) != 1 {
		t.Fatalf("parameters = %d, want 1", len(rules.Parameters))
	}

	elig, err := ops.GetCampaignEligibility(ctx, c.ID)
	if err != nil {
		t.Fatalf("GetCampaignEligibility(): %v", err)
	}
	if len(elig) != 1 {
		t.Fatalf("eligibility = %d, want 1", len(elig))
	}
	if elig[0].DeviceClass != "weather-station" {
		t.Errorf("device_class = %q, want %q", elig[0].DeviceClass, "weather-station")
	}
}
