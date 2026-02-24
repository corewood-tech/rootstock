package scitizen

import (
	"context"

	scitizenops "rootstock/web-server/ops/scitizen"
)

// CampaignProgressFlow returns contributions and score for a scitizen.
// Graph node: 0x31 — implements FR-095 (0xf)
type CampaignProgressFlow struct {
	scitizenOps *scitizenops.Ops
}

// NewCampaignProgressFlow creates the flow with its required ops.
func NewCampaignProgressFlow(scitizenOps *scitizenops.Ops) *CampaignProgressFlow {
	return &CampaignProgressFlow{scitizenOps: scitizenOps}
}

// Run returns contributions and score.
func (f *CampaignProgressFlow) Run(ctx context.Context, userID string) (*ContributionsResult, error) {
	histories, err := f.scitizenOps.GetContributions(ctx, userID)
	if err != nil {
		return nil, err
	}

	dashboard, err := f.scitizenOps.GetDashboard(ctx, userID)
	if err != nil {
		return nil, err
	}

	out := make([]ReadingHistory, len(histories))
	for i, h := range histories {
		out[i] = ReadingHistory{
			DeviceID: h.DeviceID, CampaignID: h.CampaignID,
			Total: h.Total, Accepted: h.Accepted, Rejected: h.Rejected,
		}
	}

	badges := make([]Badge, len(dashboard.Badges))
	for i, b := range dashboard.Badges {
		badges[i] = Badge{ID: b.ID, BadgeType: b.BadgeType, AwardedAt: b.AwardedAt}
	}

	return &ContributionsResult{
		Histories:         out,
		ContributionScore: dashboard.ContributionScore,
		Badges:            badges,
	}, nil
}
