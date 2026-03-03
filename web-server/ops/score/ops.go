package score

import (
	"context"

	scorerepo "rootstock/web-server/repo/score"
)

// Ops holds score/gamification operations. Each method is one op.
type Ops struct {
	repo scorerepo.Repository
}

// NewOps creates score ops backed by the given repository.
func NewOps(repo scorerepo.Repository) *Ops {
	return &Ops{repo: repo}
}

// UpdateScore recomputes a scitizen's contribution score.
// Op #26: FR-034
func (o *Ops) UpdateScore(ctx context.Context, input UpsertScoreInput) (*Score, error) {
	result, err := o.repo.UpsertScore(ctx, toRepoUpsertInput(input))
	if err != nil {
		return nil, err
	}
	return fromRepoScore(result), nil
}

// CheckMilestones evaluates badge criteria against the current score.
// Op #27: FR-035
func (o *Ops) CheckMilestones(ctx context.Context, scitizenID string) (*Score, error) {
	result, err := o.repo.GetScore(ctx, scitizenID)
	if err != nil {
		return nil, err
	}
	return fromRepoScore(result), nil
}

// AwardBadge grants a badge to a scitizen.
// Op #28: FR-035
func (o *Ops) AwardBadge(ctx context.Context, scitizenID string, badgeType string) error {
	return o.repo.AwardBadge(ctx, scitizenID, badgeType)
}

// GrantSweepstakes adds entries at score milestones.
// Op #29: FR-036
func (o *Ops) GrantSweepstakes(ctx context.Context, input GrantSweepstakesInput) error {
	return o.repo.GrantSweepstakes(ctx, toRepoGrantSweepstakesInput(input))
}

func toRepoUpsertInput(in UpsertScoreInput) scorerepo.UpsertScoreInput {
	return scorerepo.UpsertScoreInput{
		ScitizenID:  in.ScitizenID,
		Volume:      in.Volume,
		QualityRate: in.QualityRate,
		Consistency: in.Consistency,
		Diversity:   in.Diversity,
		Total:       in.Total,
	}
}

func toRepoGrantSweepstakesInput(in GrantSweepstakesInput) scorerepo.GrantSweepstakesInput {
	return scorerepo.GrantSweepstakesInput{
		ScitizenID:       in.ScitizenID,
		Entries:          in.Entries,
		MilestoneTrigger: in.MilestoneTrigger,
	}
}

// GetLeaderboard returns a ranked leaderboard.
func (o *Ops) GetLeaderboard(ctx context.Context, input GetLeaderboardInput) (*LeaderboardResult, error) {
	result, err := o.repo.GetLeaderboard(ctx, toRepoGetLeaderboardInput(input))
	if err != nil {
		return nil, err
	}
	return fromRepoLeaderboardResult(result), nil
}

// GetBadges returns all badges awarded to a scitizen.
func (o *Ops) GetBadges(ctx context.Context, scitizenID string) ([]Badge, error) {
	results, err := o.repo.GetBadges(ctx, scitizenID)
	if err != nil {
		return nil, err
	}
	out := make([]Badge, len(results))
	for i, b := range results {
		out[i] = Badge{
			ID:        b.ID,
			BadgeType: b.BadgeType,
			AwardedAt: b.AwardedAt,
		}
	}
	return out, nil
}

func fromRepoScore(r *scorerepo.Score) *Score {
	return &Score{
		ScitizenID:  r.ScitizenID,
		Volume:      r.Volume,
		QualityRate: r.QualityRate,
		Consistency: r.Consistency,
		Diversity:   r.Diversity,
		Total:       r.Total,
		UpdatedAt:   r.UpdatedAt,
	}
}

func toRepoGetLeaderboardInput(in GetLeaderboardInput) scorerepo.GetLeaderboardInput {
	return scorerepo.GetLeaderboardInput{
		CampaignID:  in.CampaignID,
		TimePeriod:  in.TimePeriod,
		Limit:       in.Limit,
		Offset:      in.Offset,
		RequesterID: in.RequesterID,
	}
}

func fromRepoLeaderboardResult(r *scorerepo.LeaderboardResult) *LeaderboardResult {
	entries := make([]LeaderboardEntry, len(r.Entries))
	for i, e := range r.Entries {
		entries[i] = LeaderboardEntry{
			Rank:          e.Rank,
			ScitizenID:    e.ScitizenID,
			Score:         e.Score,
			BadgeCount:    e.BadgeCount,
			CampaignCount: e.CampaignCount,
		}
	}
	result := &LeaderboardResult{
		Entries: entries,
		Total:   r.Total,
	}
	if r.Requester != nil {
		result.Requester = &LeaderboardEntry{
			Rank:          r.Requester.Rank,
			ScitizenID:    r.Requester.ScitizenID,
			Score:         r.Requester.Score,
			BadgeCount:    r.Requester.BadgeCount,
			CampaignCount: r.Requester.CampaignCount,
		}
	}
	return result
}
