package reading

import "time"

// PersistReadingInput is what the PersistReading op sends to the repository.
type PersistReadingInput struct {
	DeviceID        string
	CampaignID      string
	Value           float64
	Timestamp       time.Time
	Geolocation     string // GeoJSON point, may be empty
	FirmwareVersion string
	CertSerial      string
}

// QueryReadingsInput is what the QueryReadings op sends to the repository.
type QueryReadingsInput struct {
	CampaignID string
	DeviceID   string
	Status     string
	Since      *time.Time
	Until      *time.Time
	Limit      int
}

// QuarantineByWindowInput is what the QuarantineByWindow op sends.
type QuarantineByWindowInput struct {
	DeviceIDs []string
	Since     time.Time
	Until     time.Time
	Reason    string
}
