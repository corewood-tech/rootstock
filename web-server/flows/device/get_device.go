package device

import (
	"context"

	deviceops "rootstock/web-server/ops/device"
)

// GetDeviceFlow retrieves a device by ID.
type GetDeviceFlow struct {
	deviceOps *deviceops.Ops
}

// NewGetDeviceFlow creates the flow with its required ops.
func NewGetDeviceFlow(deviceOps *deviceops.Ops) *GetDeviceFlow {
	return &GetDeviceFlow{deviceOps: deviceOps}
}

// Run retrieves a device by ID.
func (f *GetDeviceFlow) Run(ctx context.Context, input GetDeviceInput) (*Device, error) {
	result, err := f.deviceOps.GetDevice(ctx, input.DeviceID)
	if err != nil {
		return nil, err
	}
	return fromOpsDevice(result), nil
}

func fromOpsDevice(d *deviceops.Device) *Device {
	return &Device{
		ID:              d.ID,
		OwnerID:         d.OwnerID,
		Status:          d.Status,
		Class:           d.Class,
		FirmwareVersion: d.FirmwareVersion,
		Tier:            d.Tier,
		Sensors:         d.Sensors,
		CertSerial:      d.CertSerial,
		CreatedAt:       d.CreatedAt,
	}
}
