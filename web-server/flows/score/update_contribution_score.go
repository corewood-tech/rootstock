package score

import (
	"context"

	scoreops "rootstock/web-server/ops/score"
)

// Milestone defines a badge and sweepstakes grant at a volume threshold.
type Milestone struct {
	Volume        int
	BadgeType     string
	SweepEntries  int
}

// DefaultMilestones are the milestones defined by FR-035/FR-036.
var DefaultMilestones = []Milestone{
	{Volume: 1, BadgeType: "first-contribution", SweepEntries: 1},
	{Volume: 100, BadgeType: "100-readings", SweepEntries: 5},
	{Volume: 1000, BadgeType: "1000-readings", SweepEntries: 25},
}

// UpdateContributionScoreFlow orchestrates score update, milestone checks,
// badge awards, and sweepstakes grants.
type UpdateContributionScoreFlow struct {
	scoreOps   *scoreops.Ops
	milestones []Milestone
}

// NewUpdateContributionScoreFlow creates the flow with its required ops.
func NewUpdateContributionScoreFlow(scoreOps *scoreops.Ops) *UpdateContributionScoreFlow {
	return &UpdateContributionScoreFlow{
		scoreOps:   scoreOps,
		milestones: DefaultMilestones,
	}
}

// Run updates the score, checks milestones, awards badges, and grants sweepstakes.
func (f *UpdateContributionScoreFlow) Run(ctx context.Context, input UpdateContributionScoreInput) (*UpdateContributionScoreResult, error) {
	// 1. Update the score
	updated, err := f.scoreOps.UpdateScore(ctx, scoreops.UpsertScoreInput{
		ScitizenID:  input.ScitizenID,
		Volume:      input.Volume,
		QualityRate: input.QualityRate,
		Consistency: input.Consistency,
		Diversity:   input.Diversity,
		Total:       input.Total,
	})
	if err != nil {
		return nil, err
	}

	// 2. Check milestones — get current score to evaluate thresholds
	current, err := f.scoreOps.CheckMilestones(ctx, input.ScitizenID)
	if err != nil {
		return nil, err
	}

	// 3. Determine which milestones are newly crossed
	var badgesAwarded []string
	var totalSweepEntries int

	for _, m := range f.milestones {
		if current.Volume >= m.Volume {
			// Try to award — AwardBadge is idempotent (INSERT, duplicate key is a no-op concern at repo level)
			if err := f.scoreOps.AwardBadge(ctx, input.ScitizenID, m.BadgeType); err != nil {
				return nil, err
			}
			badgesAwarded = append(badgesAwarded, m.BadgeType)

			if m.SweepEntries > 0 {
				if err := f.scoreOps.GrantSweepstakes(ctx, scoreops.GrantSweepstakesInput{
					ScitizenID:       input.ScitizenID,
					Entries:          m.SweepEntries,
					MilestoneTrigger: m.BadgeType,
				}); err != nil {
					return nil, err
				}
				totalSweepEntries += m.SweepEntries
			}
		}
	}

	return &UpdateContributionScoreResult{
		Score: Score{
			ScitizenID:  updated.ScitizenID,
			Volume:      updated.Volume,
			QualityRate: updated.QualityRate,
			Consistency: updated.Consistency,
			Diversity:   updated.Diversity,
			Total:       updated.Total,
			UpdatedAt:   updated.UpdatedAt,
		},
		BadgesAwarded: badgesAwarded,
		SweepEntries:  totalSweepEntries,
	}, nil
}
