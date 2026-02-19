package campaign

import "time"

// CreateCampaignInput is what callers send to CreateCampaignFlow.
type CreateCampaignInput struct {
	OrgID       string
	CreatedBy   string
	WindowStart *time.Time
	WindowEnd   *time.Time
	Parameters  []ParameterInput
	Regions     []RegionInput
	Eligibility []EligibilityInput
}

type ParameterInput struct {
	Name      string
	Unit      string
	MinRange  *float64
	MaxRange  *float64
	Precision *int
}

type RegionInput struct {
	GeoJSON string
}

type EligibilityInput struct {
	DeviceClass     string
	Tier            int
	RequiredSensors []string
	FirmwareMin     string
}

// BrowseCampaignsInput is what callers send to BrowseCampaignsFlow.
type BrowseCampaignsInput struct {
	Status    string
	OrgID     string
	Longitude *float64
	Latitude  *float64
	RadiusKm  *float64
}
