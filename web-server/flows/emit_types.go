package flows

import "time"

// Campaign is the campaign record returned by CreateCampaignFlow.
type Campaign struct {
	ID          string
	OrgID       string
	Status      string
	WindowStart *time.Time
	WindowEnd   *time.Time
	CreatedBy   string
	CreatedAt   time.Time
}

// Reading is the reading record returned by IngestReadingFlow.
type Reading struct {
	ID               string
	DeviceID         string
	CampaignID       string
	Value            float64
	Timestamp        time.Time
	Geolocation      *string
	FirmwareVersion  string
	CertSerial       string
	IngestedAt       time.Time
	Status           string
	QuarantineReason *string
}
