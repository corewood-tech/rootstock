package device

import "time"

// Device is the device record returned by device ops.
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

// DeviceCapabilities holds what matters for eligibility checks.
type DeviceCapabilities struct {
	Class           string
	Tier            int
	Sensors         []string
	FirmwareVersion string
}

// EnrollmentCode is the enrollment code record.
type EnrollmentCode struct {
	Code      string
	DeviceID  string
	ExpiresAt time.Time
	Used      bool
}
