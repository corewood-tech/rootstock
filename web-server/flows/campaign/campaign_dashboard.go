package campaign

import (
	"context"

	readingops "rootstock/web-server/ops/reading"
)

// DashboardFlow returns quality metrics for a campaign.
type DashboardFlow struct {
	readingOps *readingops.Ops
}

// NewDashboardFlow creates the flow with its required ops.
func NewDashboardFlow(readingOps *readingops.Ops) *DashboardFlow {
	return &DashboardFlow{readingOps: readingOps}
}

// Run fetches campaign quality metrics.
func (f *DashboardFlow) Run(ctx context.Context, campaignID string) (*Dashboard, error) {
	metrics, err := f.readingOps.GetCampaignQuality(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	return &Dashboard{
		CampaignID:      metrics.CampaignID,
		AcceptedCount:   metrics.AcceptedCount,
		QuarantineCount: metrics.QuarantineCount,
	}, nil
}
