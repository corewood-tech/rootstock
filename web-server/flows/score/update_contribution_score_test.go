package score

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"rootstock/web-server/config"
	scoreops "rootstock/web-server/ops/score"
	scorerepo "rootstock/web-server/repo/score"
	sqlmigrate "rootstock/web-server/repo/sql/migrate"
)

func setupUpdateScoreTest(t *testing.T) (*UpdateContributionScoreFlow, *pgxpool.Pool) {
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
	ops := scoreops.NewOps(repo)
	flow := NewUpdateContributionScoreFlow(ops)

	t.Cleanup(func() {
		repo.Shutdown()
		pool.Close()
	})

	return flow, pool
}

func TestUpdateScoreFirstContribution(t *testing.T) {
	flow, pool := setupUpdateScoreTest(t)
	ctx := context.Background()

	result, err := flow.Run(ctx, UpdateContributionScoreInput{
		ScitizenID:  "sci-1",
		Volume:      1,
		QualityRate: 1.0,
		Consistency: 1.0,
		Diversity:   1,
		Total:       10.0,
	})
	if err != nil {
		t.Fatalf("Run(): %v", err)
	}

	if result.Score.Volume != 1 {
		t.Errorf("volume = %d, want 1", result.Score.Volume)
	}

	// Should have awarded "first-contribution" badge
	found := false
	for _, b := range result.BadgesAwarded {
		if b == "first-contribution" {
			found = true
		}
	}
	if !found {
		t.Error("expected first-contribution badge to be awarded")
	}

	// Verify badge persisted
	var badgeCount int
	pool.QueryRow(ctx, "SELECT count(*) FROM badges WHERE scitizen_id = 'sci-1'").Scan(&badgeCount)
	if badgeCount != 1 {
		t.Errorf("badge count = %d, want 1", badgeCount)
	}

	// Verify sweepstakes entry persisted
	var sweepCount int
	pool.QueryRow(ctx, "SELECT count(*) FROM sweepstakes_entries WHERE scitizen_id = 'sci-1'").Scan(&sweepCount)
	if sweepCount != 1 {
		t.Errorf("sweepstakes count = %d, want 1", sweepCount)
	}
}

func TestUpdateScoreIdempotent(t *testing.T) {
	flow, pool := setupUpdateScoreTest(t)
	ctx := context.Background()

	input := UpdateContributionScoreInput{
		ScitizenID:  "sci-1",
		Volume:      1,
		QualityRate: 1.0,
		Consistency: 1.0,
		Diversity:   1,
		Total:       10.0,
	}

	// Run twice with same volume
	if _, err := flow.Run(ctx, input); err != nil {
		t.Fatalf("first Run(): %v", err)
	}
	if _, err := flow.Run(ctx, input); err != nil {
		t.Fatalf("second Run(): %v", err)
	}

	// Should still have exactly 1 badge, not 2
	var badgeCount int
	pool.QueryRow(ctx, "SELECT count(*) FROM badges WHERE scitizen_id = 'sci-1'").Scan(&badgeCount)
	if badgeCount != 1 {
		t.Errorf("badge count = %d, want 1 (idempotent)", badgeCount)
	}

	var sweepCount int
	pool.QueryRow(ctx, "SELECT count(*) FROM sweepstakes_entries WHERE scitizen_id = 'sci-1'").Scan(&sweepCount)
	if sweepCount != 1 {
		t.Errorf("sweepstakes count = %d, want 1 (idempotent)", sweepCount)
	}
}

func TestUpdateScoreMultipleMilestones(t *testing.T) {
	flow, pool := setupUpdateScoreTest(t)
	ctx := context.Background()

	// Jump straight to 100 readings — should trigger both first-contribution and 100-readings
	result, err := flow.Run(ctx, UpdateContributionScoreInput{
		ScitizenID:  "sci-1",
		Volume:      100,
		QualityRate: 0.95,
		Consistency: 0.9,
		Diversity:   3,
		Total:       200.0,
	})
	if err != nil {
		t.Fatalf("Run(): %v", err)
	}

	if len(result.BadgesAwarded) != 2 {
		t.Errorf("badges awarded = %d, want 2", len(result.BadgesAwarded))
	}

	// 1 entry for first-contribution + 5 entries for 100-readings = 6
	if result.SweepEntries != 6 {
		t.Errorf("sweep entries = %d, want 6", result.SweepEntries)
	}

	var badgeCount int
	pool.QueryRow(ctx, "SELECT count(*) FROM badges WHERE scitizen_id = 'sci-1'").Scan(&badgeCount)
	if badgeCount != 2 {
		t.Errorf("badge count = %d, want 2", badgeCount)
	}
}
