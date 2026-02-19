package device

// CreateDeviceInput is what callers send to CreateDevice.
type CreateDeviceInput struct {
	OwnerID         string
	Class           string
	FirmwareVersion string
	Tier            int
	Sensors         []string
}

// QueryByClassInput is what callers send to QueryDevicesByClass.
type QueryByClassInput struct {
	Class          string
	FirmwareMinGte string
	FirmwareMaxLte string
}

// GenerateCodeInput is what callers send to GenerateEnrollmentCode.
type GenerateCodeInput struct {
	DeviceID string
	Code     string
	TTL      int
}
