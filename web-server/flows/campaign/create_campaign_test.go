package campaign

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	campaignops "rootstock/web-server/ops/campaign"
	"rootstock/web-server/config"
	campaignrepo "rootstock/web-server/repo/campaign"
	sqlmigrate "rootstock/web-server/repo/sql/migrate"
)

func setupCreateCampaignTest(t *testing.T) (*CreateCampaignFlow, *pgxpool.Pool) {
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
	pool.Exec(ctx, "TRUNCATE campaigns CASCADE")

	cRepo := campaignrepo.NewRepository(pool)
	cOps := campaignops.NewOps(cRepo)
	flow := NewCreateCampaignFlow(cOps)

	t.Cleanup(func() {
		cRepo.Shutdown()
		pool.Close()
	})

	return flow, pool
}

func TestCreateCampaignBasic(t *testing.T) {
	flow, _ := setupCreateCampaignTest(t)
	ctx := context.Background()

	campaign, err := flow.Run(ctx, CreateCampaignInput{
		OrgID:     "org-1",
		CreatedBy: "user-1",
	})
	if err != nil {
		t.Fatalf("Run(): %v", err)
	}
	if campaign.ID == "" {
		t.Error("expected non-empty campaign ID")
	}
	if campaign.Status != "draft" {
		t.Errorf("status = %q, want draft", campaign.Status)
	}
}

func TestCreateCampaignWithParameters(t *testing.T) {
	flow, pool := setupCreateCampaignTest(t)
	ctx := context.Background()

	min := 0.0
	max := 100.0
	start := time.Now().UTC().Add(-1 * time.Hour)
	end := time.Now().UTC().Add(1 * time.Hour)

	campaign, err := flow.Run(ctx, CreateCampaignInput{
		OrgID:       "org-1",
		CreatedBy:   "user-1",
		WindowStart: &start,
		WindowEnd:   &end,
		Parameters:  []ParameterInput{{Name: "temp", Unit: "celsius", MinRange: &min, MaxRange: &max}},
	})
	if err != nil {
		t.Fatalf("Run(): %v", err)
	}

	var count int
	pool.QueryRow(ctx, "SELECT count(*) FROM campaign_parameters WHERE campaign_id = $1", campaign.ID).Scan(&count)
	if count != 1 {
		t.Errorf("parameter count = %d, want 1", count)
	}
}
