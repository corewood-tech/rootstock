package mqtt

import (
	"context"

	mqttrepo "rootstock/web-server/repo/mqtt"
)

// Ops holds MQTT operations. Each method is one op.
type Ops struct {
	repo mqttrepo.Repository
}

// NewOps creates MQTT ops backed by the given repository.
func NewOps(repo mqttrepo.Repository) *Ops {
	return &Ops{repo: repo}
}

// PushDeviceConfig publishes a retained config payload to a device's config topic.
func (o *Ops) PushDeviceConfig(ctx context.Context, input PushDeviceConfigInput) error {
	return o.repo.PushDeviceConfig(ctx, mqttrepo.PushConfigInput{
		DeviceID: input.DeviceID,
		Payload:  input.Payload,
	})
}
