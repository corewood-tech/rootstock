package campaign

import "time"

// Campaign is the campaign record returned by campaign ops.
type Campaign struct {
	ID          string
	OrgID       string
	Status      string
	WindowStart *time.Time
	WindowEnd   *time.Time
	CreatedBy   string
	CreatedAt   time.Time
}

// CampaignRules holds validation criteria for ingestion.
type CampaignRules struct {
	CampaignID  string
	Parameters  []Parameter
	Regions     []Region
	WindowStart *time.Time
	WindowEnd   *time.Time
}

type Parameter struct {
	Name      string
	Unit      string
	MinRange  *float64
	MaxRange  *float64
	Precision *int
}

type Region struct {
	GeoJSON string
}

// EligibilityCriteria holds what devices can participate.
type EligibilityCriteria struct {
	DeviceClass     string
	Tier            int
	RequiredSensors []string
	FirmwareMin     string
}
