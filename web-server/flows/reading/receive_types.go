package reading

import "time"

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

// ExportDataInput is what callers send to ExportDataFlow.
type ExportDataInput struct {
	CampaignID string
	Secret     string
	Limit      int
	Offset     int
}
