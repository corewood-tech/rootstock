package campaign

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	campaignops "rootstock/web-server/ops/campaign"
	graphops "rootstock/web-server/ops/graph"
	"rootstock/web-server/config"
	campaignrepo "rootstock/web-server/repo/campaign"
	graphrepo "rootstock/web-server/repo/graph"
	sqlmigrate "rootstock/web-server/repo/sql/migrate"
)

func setupPublishCampaignTest(t *testing.T) (*CreateCampaignFlow, *PublishCampaignFlow, *pgxpool.Pool) {
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

	gRepo, err := graphrepo.NewDgraphRepository("dgraph-alpha:9080")
	if err != nil {
		t.Fatalf("create graph repo: %v", err)
	}
	gOps := graphops.NewOps(gRepo)

	createFlow := NewCreateCampaignFlow(cOps, gOps)
	publishFlow := NewPublishCampaignFlow(cOps, gOps)

	t.Cleanup(func() {
		cRepo.Shutdown()
		gRepo.Shutdown()
		pool.Close()
	})

	return createFlow, publishFlow, pool
}

func TestPublishCampaign(t *testing.T) {
	createFlow, publishFlow, _ := setupPublishCampaignTest(t)
	ctx := context.Background()

	campaign, err := createFlow.Run(ctx, CreateCampaignInput{
		OrgID:     "org-1",
		CreatedBy: "user-1",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if campaign.Status != "draft" {
		t.Fatalf("status = %q, want draft", campaign.Status)
	}

	if err := publishFlow.Run(ctx, campaign.ID); err != nil {
		t.Fatalf("publish: %v", err)
	}
}
