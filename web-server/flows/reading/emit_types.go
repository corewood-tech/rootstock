package reading

import "time"

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

// ExportDataResult is the result of ExportDataFlow.
type ExportDataResult struct {
	Readings []ExportedReading
}

// ExportedReading is a pseudonymized reading for export.
type ExportedReading struct {
	PseudoDeviceID  string
	CampaignID      string
	Value           float64
	Timestamp       time.Time
	Geolocation     *string
	FirmwareVersion string
	IngestedAt      time.Time
	Status          string
}
