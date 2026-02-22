package campaign

import (
	"context"
	"log/slog"

	campaignops "rootstock/web-server/ops/campaign"
	graphops "rootstock/web-server/ops/graph"
)

// CreateCampaignFlow orchestrates campaign creation.
type CreateCampaignFlow struct {
	campaignOps *campaignops.Ops
	graphOps    *graphops.Ops
}

// NewCreateCampaignFlow creates the flow with its required ops.
func NewCreateCampaignFlow(campaignOps *campaignops.Ops, graphOps *graphops.Ops) *CreateCampaignFlow {
	return &CreateCampaignFlow{campaignOps: campaignOps, graphOps: graphOps}
}

// Run creates a campaign with parameters, regions, window, and eligibility.
func (f *CreateCampaignFlow) Run(ctx context.Context, input CreateCampaignInput) (*Campaign, error) {
	result, err := f.campaignOps.CreateCampaign(ctx, toOpsCampaignInput(input))
	if err != nil {
		return nil, err
	}

	// Initialize campaign state machine in graph (best-effort)
	if f.graphOps != nil {
		if _, err := f.graphOps.InitCampaignState(ctx, result.ID); err != nil {
			slog.WarnContext(ctx, "failed to init campaign graph state", "campaign_id", result.ID, "error", err)
		}
	}

	return fromOpsCampaign(result), nil
}

func toOpsCampaignInput(in CreateCampaignInput) campaignops.CreateCampaignInput {
	params := make([]campaignops.ParameterInput, len(in.Parameters))
	for i, p := range in.Parameters {
		params[i] = campaignops.ParameterInput{
			Name:      p.Name,
			Unit:      p.Unit,
			MinRange:  p.MinRange,
			MaxRange:  p.MaxRange,
			Precision: p.Precision,
		}
	}
	regions := make([]campaignops.RegionInput, len(in.Regions))
	for i, r := range in.Regions {
		regions[i] = campaignops.RegionInput{GeoJSON: r.GeoJSON}
	}
	elig := make([]campaignops.EligibilityInput, len(in.Eligibility))
	for i, e := range in.Eligibility {
		elig[i] = campaignops.EligibilityInput{
			DeviceClass:     e.DeviceClass,
			Tier:            e.Tier,
			RequiredSensors: e.RequiredSensors,
			FirmwareMin:     e.FirmwareMin,
		}
	}
	return campaignops.CreateCampaignInput{
		OrgID:       in.OrgID,
		CreatedBy:   in.CreatedBy,
		WindowStart: in.WindowStart,
		WindowEnd:   in.WindowEnd,
		Parameters:  params,
		Regions:     regions,
		Eligibility: elig,
	}
}

func fromOpsCampaign(r *campaignops.Campaign) *Campaign {
	return &Campaign{
		ID:          r.ID,
		OrgID:       r.OrgID,
		Status:      r.Status,
		WindowStart: r.WindowStart,
		WindowEnd:   r.WindowEnd,
		CreatedBy:   r.CreatedBy,
		CreatedAt:   r.CreatedAt,
	}
}
