package campaign

import (
	"context"

	campaignops "rootstock/web-server/ops/campaign"
)

// PublishCampaignFlow transitions a campaign from draft to published.
type PublishCampaignFlow struct {
	campaignOps *campaignops.Ops
}

// NewPublishCampaignFlow creates the flow with its required ops.
func NewPublishCampaignFlow(campaignOps *campaignops.Ops) *PublishCampaignFlow {
	return &PublishCampaignFlow{campaignOps: campaignOps}
}

// Run publishes a campaign by ID.
func (f *PublishCampaignFlow) Run(ctx context.Context, campaignID string) error {
	return f.campaignOps.PublishCampaign(ctx, campaignID)
}
