package reading

import "time"

// Reading is the core reading record returned by the repository.
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

// QualityMetrics holds aggregated campaign quality data.
type QualityMetrics struct {
	CampaignID      string
	AcceptedCount   int
	QuarantineCount int
}
