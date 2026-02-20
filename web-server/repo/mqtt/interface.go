package mqtt

import "context"

// Repository wraps the embedded MQTT broker's inline client for server-side
// publish operations. The broker itself is managed by the server package;
// this repo only exposes the publish surface as a clean architecture boundary.
type Repository interface {
	// PushDeviceConfig publishes a retained message to rootstock/{deviceID}/config.
	PushDeviceConfig(ctx context.Context, input PushConfigInput) error

	// PublishToDevice publishes a non-retained message to a device-specific topic.
	PublishToDevice(ctx context.Context, input PublishInput) error

	Shutdown()
}
