package campaign

import (
	"context"
	"log/slog"

	campaignops "rootstock/web-server/ops/campaign"
	graphops "rootstock/web-server/ops/graph"
)

// PublishCampaignFlow transitions a campaign from draft to published.
type PublishCampaignFlow struct {
	campaignOps *campaignops.Ops
	graphOps    *graphops.Ops
}

// NewPublishCampaignFlow creates the flow with its required ops.
func NewPublishCampaignFlow(campaignOps *campaignops.Ops, graphOps *graphops.Ops) *PublishCampaignFlow {
	return &PublishCampaignFlow{campaignOps: campaignOps, graphOps: graphOps}
}

// Run publishes a campaign by ID.
func (f *PublishCampaignFlow) Run(ctx context.Context, campaignID string) error {
	if err := f.campaignOps.PublishCampaign(ctx, campaignID); err != nil {
		return err
	}

	// Transition graph state machine (best-effort — SQL is source of truth)
	if _, err := f.graphOps.TransitionState(ctx, graphops.TransitionStateInput{
		CampaignRef: campaignID,
		EventName:   "publish",
	}); err != nil {
		slog.WarnContext(ctx, "failed to transition campaign graph state", "campaign_id", campaignID, "error", err)
	}

	return nil
}
