package device

import (
	"context"

	certops "rootstock/web-server/ops/cert"
	deviceops "rootstock/web-server/ops/device"
)

// RegisterDeviceFlow orchestrates device enrollment:
// RedeemEnrollmentCode → IssueCert → UpdateCertSerial → UpdateDeviceStatus(active).
type RegisterDeviceFlow struct {
	deviceOps *deviceops.Ops
	certOps   *certops.Ops
}

// NewRegisterDeviceFlow creates the flow with its required ops.
func NewRegisterDeviceFlow(deviceOps *deviceops.Ops, certOps *certops.Ops) *RegisterDeviceFlow {
	return &RegisterDeviceFlow{deviceOps: deviceOps, certOps: certOps}
}

// Run registers a device using an enrollment code and CSR.
func (f *RegisterDeviceFlow) Run(ctx context.Context, input RegisterDeviceInput) (*RegisterDeviceResult, error) {
	// 1. Redeem code → get device ID
	code, err := f.deviceOps.RedeemEnrollmentCode(ctx, input.EnrollmentCode)
	if err != nil {
		return nil, err
	}

	// 2. Issue cert (CN = code.DeviceID, CSR from device)
	issued, err := f.certOps.IssueCert(ctx, certops.IssueCertInput{
		DeviceID: code.DeviceID,
		CSR:      input.CSR,
	})
	if err != nil {
		return nil, err
	}

	// 3. Record cert serial on device
	if err := f.deviceOps.UpdateCertSerial(ctx, code.DeviceID, issued.Serial); err != nil {
		return nil, err
	}

	// 4. Activate
	if err := f.deviceOps.UpdateDeviceStatus(ctx, code.DeviceID, "active"); err != nil {
		return nil, err
	}

	return &RegisterDeviceResult{
		DeviceID:  code.DeviceID,
		CertPEM:   issued.CertPEM,
		Serial:    issued.Serial,
		NotBefore: issued.NotBefore,
		NotAfter:  issued.NotAfter,
	}, nil
}
