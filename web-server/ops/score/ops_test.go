package score

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"rootstock/web-server/config"
	scorerepo "rootstock/web-server/repo/score"
	sqlmigrate "rootstock/web-server/repo/sql/migrate"
)

func setupTest(t *testing.T) *Ops {
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
	pool.Exec(ctx, "TRUNCATE scores, badges, sweepstakes_entries CASCADE")

	repo := scorerepo.NewRepository(pool)
	ops := NewOps(repo)
	t.Cleanup(func() {
		repo.Shutdown()
		pool.Close()
	})
	return ops
}

func TestUpdateAndCheckScore(t *testing.T) {
	ops := setupTest(t)
	ctx := context.Background()

	s, err := ops.UpdateScore(ctx, UpsertScoreInput{
		ScitizenID: "sci-1", Volume: 50, QualityRate: 0.95, Consistency: 0.8, Diversity: 3, Total: 100.0,
	})
	if err != nil {
		t.Fatalf("UpdateScore(): %v", err)
	}
	if s.Volume != 50 {
		t.Errorf("volume = %d, want 50", s.Volume)
	}

	checked, err := ops.CheckMilestones(ctx, "sci-1")
	if err != nil {
		t.Fatalf("CheckMilestones(): %v", err)
	}
	if checked.Total != 100.0 {
		t.Errorf("total = %f, want 100.0", checked.Total)
	}
}

func TestAwardBadgeAndGrantSweepstakes(t *testing.T) {
	ops := setupTest(t)
	ctx := context.Background()

	if err := ops.AwardBadge(ctx, "sci-1", "first-contribution"); err != nil {
		t.Fatalf("AwardBadge(): %v", err)
	}

	if err := ops.GrantSweepstakes(ctx, GrantSweepstakesInput{
		ScitizenID: "sci-1", Entries: 5, MilestoneTrigger: "100-readings",
	}); err != nil {
		t.Fatalf("GrantSweepstakes(): %v", err)
	}
}
