package device

import (
	"context"

	deviceops "rootstock/web-server/ops/device"
)

// ReinstateDeviceFlow marks a revoked device as active again.
type ReinstateDeviceFlow struct {
	deviceOps *deviceops.Ops
}

// NewReinstateDeviceFlow creates the flow with its required ops.
func NewReinstateDeviceFlow(deviceOps *deviceops.Ops) *ReinstateDeviceFlow {
	return &ReinstateDeviceFlow{deviceOps: deviceOps}
}

// Run reinstates a revoked device.
func (f *ReinstateDeviceFlow) Run(ctx context.Context, input ReinstateDeviceInput) error {
	return f.deviceOps.UpdateDeviceStatus(ctx, input.DeviceID, "active")
}
