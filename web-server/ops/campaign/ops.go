package campaign

import (
	"context"

	campaignrepo "rootstock/web-server/repo/campaign"
)

// Ops holds campaign operations. Each method is one op.
type Ops struct {
	repo campaignrepo.Repository
}

// NewOps creates campaign ops backed by the given repository.
func NewOps(repo campaignrepo.Repository) *Ops {
	return &Ops{repo: repo}
}

// CreateCampaign creates a new campaign with parameters, regions, window, and eligibility.
// Op #6: FR-005–008
func (o *Ops) CreateCampaign(ctx context.Context, input CreateCampaignInput) (*Campaign, error) {
	result, err := o.repo.Create(ctx, toRepoCreateInput(input))
	if err != nil {
		return nil, err
	}
	return fromRepoCampaign(result), nil
}

// PublishCampaign transitions a campaign from draft to published.
// Op #7: FR-009
func (o *Ops) PublishCampaign(ctx context.Context, id string) error {
	return o.repo.Publish(ctx, id)
}

// ListCampaigns queries campaigns with filters.
// Op #8: FR-009, FR-012
func (o *Ops) ListCampaigns(ctx context.Context, input ListCampaignsInput) ([]Campaign, error) {
	results, err := o.repo.List(ctx, toRepoListInput(input))
	if err != nil {
		return nil, err
	}
	out := make([]Campaign, len(results))
	for i, r := range results {
		out[i] = *fromRepoCampaign(&r)
	}
	return out, nil
}

// GetCampaignRules returns validation criteria for ingestion.
// Op #9: FR-022
func (o *Ops) GetCampaignRules(ctx context.Context, campaignID string) (*CampaignRules, error) {
	result, err := o.repo.GetRules(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	return fromRepoRules(result), nil
}

// GetCampaignEligibility returns eligibility criteria for a campaign.
// Op #10: FR-019
func (o *Ops) GetCampaignEligibility(ctx context.Context, campaignID string) ([]EligibilityCriteria, error) {
	results, err := o.repo.GetEligibility(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	out := make([]EligibilityCriteria, len(results))
	for i, r := range results {
		out[i] = EligibilityCriteria{
			DeviceClass:     r.DeviceClass,
			Tier:            r.Tier,
			RequiredSensors: r.RequiredSensors,
			FirmwareMin:     r.FirmwareMin,
		}
	}
	return out, nil
}

func toRepoCreateInput(in CreateCampaignInput) campaignrepo.CreateCampaignInput {
	params := make([]campaignrepo.ParameterInput, len(in.Parameters))
	for i, p := range in.Parameters {
		params[i] = campaignrepo.ParameterInput{
			Name:      p.Name,
			Unit:      p.Unit,
			MinRange:  p.MinRange,
			MaxRange:  p.MaxRange,
			Precision: p.Precision,
		}
	}
	regions := make([]campaignrepo.RegionInput, len(in.Regions))
	for i, r := range in.Regions {
		regions[i] = campaignrepo.RegionInput{GeoJSON: r.GeoJSON}
	}
	elig := make([]campaignrepo.EligibilityInput, len(in.Eligibility))
	for i, e := range in.Eligibility {
		elig[i] = campaignrepo.EligibilityInput{
			DeviceClass:     e.DeviceClass,
			Tier:            e.Tier,
			RequiredSensors: e.RequiredSensors,
			FirmwareMin:     e.FirmwareMin,
		}
	}
	return campaignrepo.CreateCampaignInput{
		OrgID:       in.OrgID,
		CreatedBy:   in.CreatedBy,
		WindowStart: in.WindowStart,
		WindowEnd:   in.WindowEnd,
		Parameters:  params,
		Regions:     regions,
		Eligibility: elig,
	}
}

func toRepoListInput(in ListCampaignsInput) campaignrepo.ListCampaignsInput {
	return campaignrepo.ListCampaignsInput{
		Status:    in.Status,
		OrgID:     in.OrgID,
		Longitude: in.Longitude,
		Latitude:  in.Latitude,
		RadiusKm:  in.RadiusKm,
	}
}

func fromRepoCampaign(r *campaignrepo.Campaign) *Campaign {
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

func fromRepoRules(r *campaignrepo.CampaignRules) *CampaignRules {
	params := make([]Parameter, len(r.Parameters))
	for i, p := range r.Parameters {
		params[i] = Parameter{
			Name:      p.Name,
			Unit:      p.Unit,
			MinRange:  p.MinRange,
			MaxRange:  p.MaxRange,
			Precision: p.Precision,
		}
	}
	regions := make([]Region, len(r.Regions))
	for i, rg := range r.Regions {
		regions[i] = Region{GeoJSON: rg.GeoJSON}
	}
	return &CampaignRules{
		CampaignID:  r.CampaignID,
		Parameters:  params,
		Regions:     regions,
		WindowStart: r.WindowStart,
		WindowEnd:   r.WindowEnd,
	}
}
