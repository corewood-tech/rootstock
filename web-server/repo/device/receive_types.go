package device

// CreateDeviceInput is what the CreateDevice op sends to the repository.
type CreateDeviceInput struct {
	OwnerID         string
	Class           string
	FirmwareVersion string
	Tier            int
	Sensors         []string
}

// QueryByClassInput is what the QueryDevicesByClass op sends to the repository.
type QueryByClassInput struct {
	Class          string
	FirmwareMinGte string // firmware_version >= this
	FirmwareMaxLte string // firmware_version <= this
}

// GenerateCodeInput is what the GenerateEnrollmentCode op sends to the repository.
type GenerateCodeInput struct {
	DeviceID string
	Code     string
	TTL      int // seconds
}
