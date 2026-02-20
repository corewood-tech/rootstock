package device

import "time"

// Device is the device record returned by device flows.
type Device struct {
	ID              string
	OwnerID         string
	Status          string
	Class           string
	FirmwareVersion string
	Tier            int
	Sensors         []string
	CertSerial      *string
	CreatedAt       time.Time
}

// RegisterDeviceResult is the result of device enrollment.
type RegisterDeviceResult struct {
	DeviceID  string
	CertPEM   []byte
	Serial    string
	NotBefore time.Time
	NotAfter  time.Time
}

// RenewCertResult is the result of certificate renewal.
type RenewCertResult struct {
	CertPEM   []byte
	Serial    string
	NotBefore time.Time
	NotAfter  time.Time
}

// CACert is the CA certificate returned by GetCACertFlow.
type CACert struct {
	CertPEM []byte
}

// EnrollInCampaignResult is the result of the EnrollInCampaign flow.
type EnrollInCampaignResult struct {
	Enrolled bool
	Reason   string
}

// DeviceConfigPayload is the config pushed to a device after enrollment.
type DeviceConfigPayload struct {
	CampaignID string `json:"campaign_id"`
}
