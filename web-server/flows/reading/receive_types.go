package reading

import "time"

// IngestReadingInput is what callers send to IngestReadingFlow.
type IngestReadingInput struct {
	DeviceID        string
	CampaignID      string
	Values          map[string]float64 // parameter name -> value
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
