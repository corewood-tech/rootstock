package score

import (
	"context"

	deviceops "rootstock/web-server/ops/device"
	readingops "rootstock/web-server/ops/reading"
	scoreops "rootstock/web-server/ops/score"
)

// RefreshScitizenScoreFlow resolves device→scitizen, computes stats from DB,
// updates the score, checks milestones, awards badges, and grants sweepstakes.
type RefreshScitizenScoreFlow struct {
	deviceOps  *deviceops.Ops
	readingOps *readingops.Ops
	scoreOps   *scoreops.Ops
	milestones []Milestone
}

// NewRefreshScitizenScoreFlow creates the flow with its required ops.
func NewRefreshScitizenScoreFlow(deviceOps *deviceops.Ops, readingOps *readingops.Ops, scoreOps *scoreops.Ops) *RefreshScitizenScoreFlow {
	return &RefreshScitizenScoreFlow{
		deviceOps:  deviceOps,
		readingOps: readingOps,
		scoreOps:   scoreOps,
		milestones: DefaultMilestones,
	}
}

// Run refreshes a scitizen's score based on the device that just submitted a reading.
func (f *RefreshScitizenScoreFlow) Run(ctx context.Context, input RefreshScitizenScoreInput) (*UpdateContributionScoreResult, error) {
	// 1. Resolve device → scitizen
	device, err := f.deviceOps.GetDevice(ctx, input.DeviceID)
	if err != nil {
		return nil, err
	}
	scitizenID := device.OwnerID

	// 2. Compute stats from DB
	stats, err := f.readingOps.GetScitizenReadingStats(ctx, scitizenID)
	if err != nil {
		return nil, err
	}

	// 3. Compute total score
	total := float64(stats.Volume)*0.4 +
		stats.QualityRate*100*0.3 +
		stats.Consistency*100*0.2 +
		float64(stats.Diversity)*10*0.1

	// 4. Update score
	updated, err := f.scoreOps.UpdateScore(ctx, scoreops.UpsertScoreInput{
		ScitizenID:  scitizenID,
		Volume:      stats.Volume,
		QualityRate: stats.QualityRate,
		Consistency: stats.Consistency,
		Diversity:   stats.Diversity,
		Total:       total,
	})
	if err != nil {
		return nil, err
	}

	// 5. Check milestones
	current, err := f.scoreOps.CheckMilestones(ctx, scitizenID)
	if err != nil {
		return nil, err
	}

	// 6. Award badges and grant sweepstakes for crossed milestones
	var badgesAwarded []string
	var totalSweepEntries int

	for _, m := range f.milestones {
		if current.Volume >= m.Volume {
			if err := f.scoreOps.AwardBadge(ctx, scitizenID, m.BadgeType); err != nil {
				return nil, err
			}
			badgesAwarded = append(badgesAwarded, m.BadgeType)

			if m.SweepEntries > 0 {
				if err := f.scoreOps.GrantSweepstakes(ctx, scoreops.GrantSweepstakesInput{
					ScitizenID:       scitizenID,
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
