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

// Reading is the reading record returned by IngestReadingFlow.
type Reading struct {
	ID               string
	DeviceID         string
	CampaignID       string
	Values           []ReadingValue
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
	Values          map[string]float64
	Timestamp       time.Time
	Geolocation     *string
	FirmwareVersion string
	IngestedAt      time.Time
	Status          string
}
