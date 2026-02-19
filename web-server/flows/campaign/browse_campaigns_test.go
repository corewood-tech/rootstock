package campaign

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	campaignops "rootstock/web-server/ops/campaign"
	"rootstock/web-server/config"
	campaignrepo "rootstock/web-server/repo/campaign"
	sqlmigrate "rootstock/web-server/repo/sql/migrate"
)

func setupBrowseCampaignsTest(t *testing.T) (*CreateCampaignFlow, *BrowseCampaignsFlow, *pgxpool.Pool) {
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
	createFlow := NewCreateCampaignFlow(cOps)
	browseFlow := NewBrowseCampaignsFlow(cOps)

	t.Cleanup(func() {
		cRepo.Shutdown()
		pool.Close()
	})

	return createFlow, browseFlow, pool
}

func TestBrowseCampaignsEmpty(t *testing.T) {
	_, browseFlow, _ := setupBrowseCampaignsTest(t)
	ctx := context.Background()

	campaigns, err := browseFlow.Run(ctx, BrowseCampaignsInput{})
	if err != nil {
		t.Fatalf("browse: %v", err)
	}
	if len(campaigns) != 0 {
		t.Errorf("got %d campaigns, want 0", len(campaigns))
	}
}

func TestBrowseCampaignsReturnsCreated(t *testing.T) {
	createFlow, browseFlow, _ := setupBrowseCampaignsTest(t)
	ctx := context.Background()

	_, err := createFlow.Run(ctx, CreateCampaignInput{
		OrgID:     "org-1",
		CreatedBy: "user-1",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	campaigns, err := browseFlow.Run(ctx, BrowseCampaignsInput{OrgID: "org-1"})
	if err != nil {
		t.Fatalf("browse: %v", err)
	}
	if len(campaigns) != 1 {
		t.Errorf("got %d campaigns, want 1", len(campaigns))
	}
}
