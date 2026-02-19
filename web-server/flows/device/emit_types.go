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
