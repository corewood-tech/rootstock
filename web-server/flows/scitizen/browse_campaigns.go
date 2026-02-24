package scitizen

import (
	"context"

	scitizenops "rootstock/web-server/ops/scitizen"
)

// BrowseCampaignsFlow lists published campaigns for scitizen browsing.
// Graph node: 0x32 — implements FR-009 (0x3), FR-012 (0xa)
type BrowseCampaignsFlow struct {
	scitizenOps *scitizenops.Ops
}

// NewBrowseCampaignsFlow creates the flow with its required ops.
func NewBrowseCampaignsFlow(scitizenOps *scitizenops.Ops) *BrowseCampaignsFlow {
	return &BrowseCampaignsFlow{scitizenOps: scitizenOps}
}

// Run returns published campaigns matching the given filters.
func (f *BrowseCampaignsFlow) Run(ctx context.Context, input BrowseInput) ([]CampaignSummary, int, error) {
	results, total, err := f.scitizenOps.BrowseCampaigns(ctx, scitizenops.BrowseInput{
		Longitude:  input.Longitude,
		Latitude:   input.Latitude,
		RadiusKm:   input.RadiusKm,
		SensorType: input.SensorType,
		Limit:      input.Limit,
		Offset:     input.Offset,
	})
	if err != nil {
		return nil, 0, err
	}

	out := make([]CampaignSummary, len(results))
	for i, r := range results {
		out[i] = CampaignSummary{
			ID: r.ID, Status: r.Status, WindowStart: r.WindowStart, WindowEnd: r.WindowEnd,
			EnrollmentCount: r.EnrollmentCount, RequiredSensors: r.RequiredSensors, CreatedAt: r.CreatedAt,
		}
	}
	return out, total, nil
}
