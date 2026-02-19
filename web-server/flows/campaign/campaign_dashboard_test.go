package campaign

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	campaignops "rootstock/web-server/ops/campaign"
	readingops "rootstock/web-server/ops/reading"
	"rootstock/web-server/config"
	campaignrepo "rootstock/web-server/repo/campaign"
	readingrepo "rootstock/web-server/repo/reading"
	sqlmigrate "rootstock/web-server/repo/sql/migrate"
)

func setupCampaignDashboardTest(t *testing.T) (*CreateCampaignFlow, *DashboardFlow, *pgxpool.Pool) {
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
	pool.Exec(ctx, "TRUNCATE readings CASCADE")

	cRepo := campaignrepo.NewRepository(pool)
	rRepo := readingrepo.NewRepository(pool)
	cOps := campaignops.NewOps(cRepo)
	rOps := readingops.NewOps(rRepo)
	createFlow := NewCreateCampaignFlow(cOps)
	dashboardFlow := NewDashboardFlow(rOps)

	t.Cleanup(func() {
		cRepo.Shutdown()
		rRepo.Shutdown()
		pool.Close()
	})

	return createFlow, dashboardFlow, pool
}

func TestCampaignDashboardEmpty(t *testing.T) {
	createFlow, dashboardFlow, _ := setupCampaignDashboardTest(t)
	ctx := context.Background()

	c, err := createFlow.Run(ctx, CreateCampaignInput{
		OrgID:     "org-1",
		CreatedBy: "user-1",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	dashboard, err := dashboardFlow.Run(ctx, c.ID)
	if err != nil {
		t.Fatalf("dashboard: %v", err)
	}
	if dashboard.CampaignID != c.ID {
		t.Errorf("campaign_id = %q, want %q", dashboard.CampaignID, c.ID)
	}
	if dashboard.AcceptedCount != 0 {
		t.Errorf("accepted = %d, want 0", dashboard.AcceptedCount)
	}
	if dashboard.QuarantineCount != 0 {
		t.Errorf("quarantine = %d, want 0", dashboard.QuarantineCount)
	}
}
