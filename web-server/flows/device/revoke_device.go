package device

import (
	"context"

	deviceops "rootstock/web-server/ops/device"
)

// RevokeDeviceFlow marks a device as revoked.
type RevokeDeviceFlow struct {
	deviceOps *deviceops.Ops
}

// NewRevokeDeviceFlow creates the flow with its required ops.
func NewRevokeDeviceFlow(deviceOps *deviceops.Ops) *RevokeDeviceFlow {
	return &RevokeDeviceFlow{deviceOps: deviceOps}
}

// Run revokes a device.
func (f *RevokeDeviceFlow) Run(ctx context.Context, input RevokeDeviceInput) error {
	return f.deviceOps.UpdateDeviceStatus(ctx, input.DeviceID, "revoked")
}
