package device

import (
	"context"

	certops "rootstock/web-server/ops/cert"
	deviceops "rootstock/web-server/ops/device"
)

// RenewCertFlow orchestrates certificate renewal:
// IssueCert → UpdateCertSerial.
type RenewCertFlow struct {
	deviceOps *deviceops.Ops
	certOps   *certops.Ops
}

// NewRenewCertFlow creates the flow with its required ops.
func NewRenewCertFlow(deviceOps *deviceops.Ops, certOps *certops.Ops) *RenewCertFlow {
	return &RenewCertFlow{deviceOps: deviceOps, certOps: certOps}
}

// Run renews a device certificate. Device ID comes from mTLS cert CN.
func (f *RenewCertFlow) Run(ctx context.Context, input RenewCertInput) (*RenewCertResult, error) {
	// 1. Issue new cert
	issued, err := f.certOps.IssueCert(ctx, certops.IssueCertInput{
		DeviceID: input.DeviceID,
		CSR:      input.CSR,
	})
	if err != nil {
		return nil, err
	}

	// 2. Update cert serial on device
	if err := f.deviceOps.UpdateCertSerial(ctx, input.DeviceID, issued.Serial); err != nil {
		return nil, err
	}

	return &RenewCertResult{
		CertPEM:   issued.CertPEM,
		Serial:    issued.Serial,
		NotBefore: issued.NotBefore,
		NotAfter:  issued.NotAfter,
	}, nil
}
