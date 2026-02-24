package scitizen

import (
	"context"

	scitizenops "rootstock/web-server/ops/scitizen"
)

// CampaignDetailFlow returns full campaign detail for enrollment decision.
// Graph node: 0x25 — implements FR-082 (0x6)
type CampaignDetailFlow struct {
	scitizenOps *scitizenops.Ops
}

// NewCampaignDetailFlow creates the flow with its required ops.
func NewCampaignDetailFlow(scitizenOps *scitizenops.Ops) *CampaignDetailFlow {
	return &CampaignDetailFlow{scitizenOps: scitizenOps}
}

// Run returns full campaign detail.
func (f *CampaignDetailFlow) Run(ctx context.Context, campaignID string) (*CampaignDetail, error) {
	result, err := f.scitizenOps.GetCampaignDetail(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	params := make([]Parameter, len(result.Parameters))
	for i, p := range result.Parameters {
		params[i] = Parameter{Name: p.Name, Unit: p.Unit, MinRange: p.MinRange, MaxRange: p.MaxRange, Precision: p.Precision}
	}
	regions := make([]Region, len(result.Regions))
	for i, r := range result.Regions {
		regions[i] = Region{GeoJSON: r.GeoJSON}
	}
	elig := make([]EligibilityCriteria, len(result.Eligibility))
	for i, e := range result.Eligibility {
		elig[i] = EligibilityCriteria{
			DeviceClass: e.DeviceClass, Tier: e.Tier,
			RequiredSensors: e.RequiredSensors, FirmwareMin: e.FirmwareMin,
		}
	}

	return &CampaignDetail{
		CampaignID:      result.CampaignID,
		Status:          result.Status,
		WindowStart:     result.WindowStart,
		WindowEnd:       result.WindowEnd,
		Parameters:      params,
		Regions:         regions,
		Eligibility:     elig,
		EnrollmentCount: result.EnrollmentCount,
		ProgressPercent: result.ProgressPercent,
	}, nil
}
