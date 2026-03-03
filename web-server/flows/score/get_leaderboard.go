package score

import (
	"context"

	scoreops "rootstock/web-server/ops/score"
)

// GetLeaderboardFlow returns a ranked leaderboard with the requester's position.
type GetLeaderboardFlow struct {
	scoreOps *scoreops.Ops
}

// NewGetLeaderboardFlow creates the flow with its required ops.
func NewGetLeaderboardFlow(scoreOps *scoreops.Ops) *GetLeaderboardFlow {
	return &GetLeaderboardFlow{scoreOps: scoreOps}
}

// Run fetches the leaderboard.
func (f *GetLeaderboardFlow) Run(ctx context.Context, input GetLeaderboardInput) (*LeaderboardResult, error) {
	result, err := f.scoreOps.GetLeaderboard(ctx, scoreops.GetLeaderboardInput{
		CampaignID:  input.CampaignID,
		TimePeriod:  input.TimePeriod,
		Limit:       input.Limit,
		Offset:      input.Offset,
		RequesterID: input.RequesterID,
	})
	if err != nil {
		return nil, err
	}

	out := &LeaderboardResult{
		Total: result.Total,
	}
	for _, e := range result.Entries {
		out.Entries = append(out.Entries, LeaderboardEntry{
			Rank:          e.Rank,
			ScitizenID:    e.ScitizenID,
			Score:         e.Score,
			BadgeCount:    e.BadgeCount,
			CampaignCount: e.CampaignCount,
		})
	}
	if result.Requester != nil {
		out.Requester = &LeaderboardEntry{
			Rank:          result.Requester.Rank,
			ScitizenID:    result.Requester.ScitizenID,
			Score:         result.Requester.Score,
			BadgeCount:    result.Requester.BadgeCount,
			CampaignCount: result.Requester.CampaignCount,
		}
	}
	return out, nil
}
