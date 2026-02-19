package score

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"rootstock/web-server/config"
	sqlmigrate "rootstock/web-server/repo/sql/migrate"
)

func setupTest(t *testing.T) Repository {
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
	pool.Exec(ctx, "DELETE FROM sweepstakes_entries")
	pool.Exec(ctx, "DELETE FROM badges")
	pool.Exec(ctx, "DELETE FROM scores")

	repo := NewRepository(pool)
	t.Cleanup(func() {
		repo.Shutdown()
		pool.Close()
	})

	return repo
}

func TestUpsertAndGetScore(t *testing.T) {
	repo := setupTest(t)
	ctx := context.Background()

	s, err := repo.UpsertScore(ctx, UpsertScoreInput{
		ScitizenID:  "sci-1",
		Volume:      42,
		QualityRate: 0.95,
		Consistency: 0.8,
		Diversity:   3,
		Total:       120.5,
	})
	if err != nil {
		t.Fatalf("UpsertScore(): %v", err)
	}
	if s.Volume != 42 {
		t.Errorf("volume = %d, want 42", s.Volume)
	}

	got, err := repo.GetScore(ctx, "sci-1")
	if err != nil {
		t.Fatalf("GetScore(): %v", err)
	}
	if got.Total != 120.5 {
		t.Errorf("total = %f, want 120.5", got.Total)
	}

	// Upsert again — should update
	s2, err := repo.UpsertScore(ctx, UpsertScoreInput{
		ScitizenID:  "sci-1",
		Volume:      100,
		QualityRate: 0.98,
		Consistency: 0.9,
		Diversity:   5,
		Total:       250.0,
	})
	if err != nil {
		t.Fatalf("UpsertScore() update: %v", err)
	}
	if s2.Volume != 100 {
		t.Errorf("updated volume = %d, want 100", s2.Volume)
	}
}

func TestAwardBadge(t *testing.T) {
	repo := setupTest(t)
	ctx := context.Background()

	if err := repo.AwardBadge(ctx, "sci-1", "first-contribution"); err != nil {
		t.Fatalf("AwardBadge(): %v", err)
	}
	if err := repo.AwardBadge(ctx, "sci-1", "100-readings"); err != nil {
		t.Fatalf("AwardBadge() second: %v", err)
	}

	badges, err := repo.GetBadges(ctx, "sci-1")
	if err != nil {
		t.Fatalf("GetBadges(): %v", err)
	}
	if len(badges) != 2 {
		t.Fatalf("badges = %d, want 2", len(badges))
	}
	if badges[0].BadgeType != "first-contribution" {
		t.Errorf("badge[0] = %q, want %q", badges[0].BadgeType, "first-contribution")
	}
}

func TestGrantSweepstakes(t *testing.T) {
	repo := setupTest(t)
	ctx := context.Background()

	if err := repo.GrantSweepstakes(ctx, GrantSweepstakesInput{
		ScitizenID:       "sci-1",
		Entries:          5,
		MilestoneTrigger: "100-readings",
	}); err != nil {
		t.Fatalf("GrantSweepstakes(): %v", err)
	}

	entries, err := repo.GetSweepstakesEntries(ctx, "sci-1")
	if err != nil {
		t.Fatalf("GetSweepstakesEntries(): %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("entries = %d, want 1", len(entries))
	}
	if entries[0].Entries != 5 {
		t.Errorf("entries count = %d, want 5", entries[0].Entries)
	}
	if entries[0].MilestoneTrigger != "100-readings" {
		t.Errorf("trigger = %q, want %q", entries[0].MilestoneTrigger, "100-readings")
	}
}

func TestGetScoreNotFound(t *testing.T) {
	repo := setupTest(t)
	ctx := context.Background()

	_, err := repo.GetScore(ctx, "nonexistent")
	if err == nil {
		t.Error("GetScore(nonexistent) should fail")
	}
}
