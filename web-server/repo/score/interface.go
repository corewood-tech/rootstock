package score

import "context"

// Repository defines the interface for score/gamification operations.
type Repository interface {
	UpsertScore(ctx context.Context, input UpsertScoreInput) (*Score, error)
	GetScore(ctx context.Context, scitizenID string) (*Score, error)
	AwardBadge(ctx context.Context, scitizenID string, badgeType string) error
	GetBadges(ctx context.Context, scitizenID string) ([]Badge, error)
	GrantSweepstakes(ctx context.Context, input GrantSweepstakesInput) error
	GetSweepstakesEntries(ctx context.Context, scitizenID string) ([]SweepstakesEntry, error)
	Shutdown()
}
