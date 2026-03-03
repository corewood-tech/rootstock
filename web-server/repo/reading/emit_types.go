package reading

import "time"

// ReadingValue is a single parameter measurement within a reading.
type ReadingValue struct {
	ID               string
	ReadingID        string
	ParameterName    string
	Value            float64
	Status           string
	QuarantineReason *string
}

// Reading is the core reading record returned by the repository.
type Reading struct {
	ID               string
	DeviceID         string
	CampaignID       string
	Value            *float64 // legacy single-value, nullable
	Values           []ReadingValue
	Timestamp        time.Time
	Geolocation      *string
	FirmwareVersion  string
	CertSerial       string
	IngestedAt       time.Time
	Status           string
	QuarantineReason *string
}

// ParameterQuality holds per-parameter quality metrics.
type ParameterQuality struct {
	ParameterName   string
	AcceptedCount   int
	QuarantinedCount int
}

// QualityMetrics holds aggregated campaign quality data.
type QualityMetrics struct {
	CampaignID      string
	AcceptedCount   int
	QuarantineCount int
	PerParameter    []ParameterQuality
}

// DeviceBreakdown holds per-device stats for a campaign.
type DeviceBreakdown struct {
	PseudoDeviceID string
	DeviceClass    string
	AcceptanceRate float64
	ReadingCount   int
	LastSeen       *string
}

// TemporalBucket holds reading counts for a time bucket.
type TemporalBucket struct {
	Bucket string // ISO timestamp
	Count  int
}

// EnrollmentFunnel holds enrollment stages.
type EnrollmentFunnel struct {
	Enrolled     int
	Active       int
	Contributing int
}

// ScitizenReadingStats holds aggregated reading statistics for a scitizen.
type ScitizenReadingStats struct {
	Volume      int
	QualityRate float64
	Consistency float64
	Diversity   int
}
