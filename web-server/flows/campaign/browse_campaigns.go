package campaign

import (
	"context"

	campaignops "rootstock/web-server/ops/campaign"
)

// BrowseCampaignsFlow lists campaigns with optional filters.
type BrowseCampaignsFlow struct {
	campaignOps *campaignops.Ops
}

// NewBrowseCampaignsFlow creates the flow with its required ops.
func NewBrowseCampaignsFlow(campaignOps *campaignops.Ops) *BrowseCampaignsFlow {
	return &BrowseCampaignsFlow{campaignOps: campaignOps}
}

// Run queries campaigns matching the given filters.
func (f *BrowseCampaignsFlow) Run(ctx context.Context, input BrowseCampaignsInput) ([]Campaign, error) {
	results, err := f.campaignOps.ListCampaigns(ctx, toOpsListInput(input))
	if err != nil {
		return nil, err
	}
	out := make([]Campaign, len(results))
	for i, r := range results {
		out[i] = *fromOpsCampaign(&r)
	}
	return out, nil
}

func toOpsListInput(in BrowseCampaignsInput) campaignops.ListCampaignsInput {
	return campaignops.ListCampaignsInput{
		Status:    in.Status,
		OrgID:     in.OrgID,
		Longitude: in.Longitude,
		Latitude:  in.Latitude,
		RadiusKm:  in.RadiusKm,
	}
}
