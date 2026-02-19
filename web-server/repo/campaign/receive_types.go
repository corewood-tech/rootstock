package campaign

import "time"

// CreateCampaignInput is what the CreateCampaign op sends to the repository.
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
	GeoJSON string // GeoJSON geometry
}

type EligibilityInput struct {
	DeviceClass     string
	Tier            int
	RequiredSensors []string
	FirmwareMin     string
}

// ListCampaignsInput is what the ListCampaigns op sends to the repository.
type ListCampaignsInput struct {
	Status    string
	OrgID     string
	Longitude *float64
	Latitude  *float64
	RadiusKm  *float64
}
