package reading

import "time"

// PersistReadingInput is what callers send to PersistReading.
type PersistReadingInput struct {
	DeviceID        string
	CampaignID      string
	Value           float64
	Timestamp       time.Time
	Geolocation     string
	FirmwareVersion string
	CertSerial      string
}

// QueryReadingsInput is what callers send to QueryReadings.
type QueryReadingsInput struct {
	CampaignID string
	DeviceID   string
	Status     string
	Since      *time.Time
	Until      *time.Time
	Limit      int
	Offset     int
}

// QuarantineByWindowInput is what callers send to QuarantineByWindow.
type QuarantineByWindowInput struct {
	DeviceIDs []string
	Since     time.Time
	Until     time.Time
	Reason    string
}
