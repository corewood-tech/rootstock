package reading

import "time"

// Reading is the reading record returned by reading ops.
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

// ScitizenReadingStats holds aggregated reading statistics for a scitizen.
type ScitizenReadingStats struct {
	Volume      int
	QualityRate float64
	Consistency float64
	Diversity   int
}
