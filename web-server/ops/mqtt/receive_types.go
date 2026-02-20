package mqtt

// PushDeviceConfigInput is what callers send to PushDeviceConfig.
type PushDeviceConfigInput struct {
	DeviceID string
	Payload  []byte
}
