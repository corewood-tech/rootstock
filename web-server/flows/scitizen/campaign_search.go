package scitizen

import (
	"context"

	scitizenops "rootstock/web-server/ops/scitizen"
)

// CampaignSearchFlow performs full-text search across published campaigns.
// Graph node: 0x30 — implements FR-088 (0x15)
type CampaignSearchFlow struct {
	scitizenOps *scitizenops.Ops
}

// NewCampaignSearchFlow creates the flow with its required ops.
func NewCampaignSearchFlow(scitizenOps *scitizenops.Ops) *CampaignSearchFlow {
	return &CampaignSearchFlow{scitizenOps: scitizenOps}
}

// Run searches published campaigns by query string.
func (f *CampaignSearchFlow) Run(ctx context.Context, input SearchInput) ([]CampaignSummary, int, error) {
	results, total, err := f.scitizenOps.SearchCampaigns(ctx, scitizenops.SearchInput{
		Query:  input.Query,
		Limit:  input.Limit,
		Offset: input.Offset,
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
