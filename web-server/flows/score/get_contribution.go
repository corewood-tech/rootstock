package score

import (
	"context"

	scoreops "rootstock/web-server/ops/score"
)

// GetContributionFlow returns a scitizen's contribution profile (score + badges).
type GetContributionFlow struct {
	scoreOps *scoreops.Ops
}

// NewGetContributionFlow creates the flow with its required ops.
func NewGetContributionFlow(scoreOps *scoreops.Ops) *GetContributionFlow {
	return &GetContributionFlow{scoreOps: scoreOps}
}

// Run fetches the score and badges for a scitizen.
func (f *GetContributionFlow) Run(ctx context.Context, scitizenID string) (*Contribution, error) {
	score, err := f.scoreOps.CheckMilestones(ctx, scitizenID)
	if err != nil {
		return nil, err
	}

	opsBadges, err := f.scoreOps.GetBadges(ctx, scitizenID)
	if err != nil {
		return nil, err
	}

	badges := make([]Badge, len(opsBadges))
	for i, b := range opsBadges {
		badges[i] = Badge{
			ID:        b.ID,
			BadgeType: b.BadgeType,
			AwardedAt: b.AwardedAt,
		}
	}

	return &Contribution{
		ScitizenID:  score.ScitizenID,
		Volume:      score.Volume,
		QualityRate: score.QualityRate,
		Consistency: score.Consistency,
		Diversity:   score.Diversity,
		Total:       score.Total,
		UpdatedAt:   score.UpdatedAt,
		Badges:      badges,
	}, nil
}
