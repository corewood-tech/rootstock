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
