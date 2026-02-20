package mqtt

// PushConfigInput is the input for PushDeviceConfig.
type PushConfigInput struct {
	DeviceID string
	Payload  []byte
}

// PublishInput is the input for PublishToDevice.
type PublishInput struct {
	Topic   string
	Payload []byte
	QoS     byte
}
