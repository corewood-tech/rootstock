package flows

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

// IngestReadingInput is what callers send to IngestReadingFlow.
type IngestReadingInput struct {
	DeviceID        string
	CampaignID      string
	Value           float64
	Timestamp       time.Time
	Geolocation     string
	FirmwareVersion string
	CertSerial      string
}
