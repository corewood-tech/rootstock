package device

import (
	"context"

	devicerepo "rootstock/web-server/repo/device"
)

// Ops holds device operations. Each method is one op.
type Ops struct {
	repo devicerepo.Repository
}

// NewOps creates device ops backed by the given repository.
func NewOps(repo devicerepo.Repository) *Ops {
	return &Ops{repo: repo}
}

// GenerateEnrollmentCode generates a one-time enrollment code for a device.
// Op #11: FR-013
func (o *Ops) GenerateEnrollmentCode(ctx context.Context, input GenerateCodeInput) (*EnrollmentCode, error) {
	result, err := o.repo.GenerateEnrollmentCode(ctx, toRepoGenerateCodeInput(input))
	if err != nil {
		return nil, err
	}
	return fromRepoEnrollmentCode(result), nil
}

// RedeemEnrollmentCode validates and marks a code as used.
// Op #12: FR-013
func (o *Ops) RedeemEnrollmentCode(ctx context.Context, code string) (*EnrollmentCode, error) {
	result, err := o.repo.RedeemEnrollmentCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return fromRepoEnrollmentCode(result), nil
}

// CreateDevice creates a device registry entry.
// Op #13: FR-016
func (o *Ops) CreateDevice(ctx context.Context, input CreateDeviceInput) (*Device, error) {
	result, err := o.repo.Create(ctx, toRepoCreateDeviceInput(input))
	if err != nil {
		return nil, err
	}
	return fromRepoDevice(result), nil
}

// GetDevice reads a device from the registry by ID.
// Op #14
func (o *Ops) GetDevice(ctx context.Context, id string) (*Device, error) {
	result, err := o.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return fromRepoDevice(result), nil
}

// GetDeviceCapabilities returns device class, tier, sensors, firmware for eligibility.
// Op #15: FR-019
func (o *Ops) GetDeviceCapabilities(ctx context.Context, id string) (*DeviceCapabilities, error) {
	result, err := o.repo.GetCapabilities(ctx, id)
	if err != nil {
		return nil, err
	}
	return &DeviceCapabilities{
		Class:           result.Class,
		Tier:            result.Tier,
		Sensors:         result.Sensors,
		FirmwareVersion: result.FirmwareVersion,
	}, nil
}

// UpdateDeviceStatus changes device status.
// Op #16: FR-030, FR-031, FR-033
func (o *Ops) UpdateDeviceStatus(ctx context.Context, id string, status string) error {
	return o.repo.UpdateStatus(ctx, id, status)
}

// QueryDevicesByClass batch queries devices by class and firmware range.
// Op #17: FR-031
func (o *Ops) QueryDevicesByClass(ctx context.Context, input QueryByClassInput) ([]Device, error) {
	results, err := o.repo.QueryByClass(ctx, toRepoQueryByClassInput(input))
	if err != nil {
		return nil, err
	}
	out := make([]Device, len(results))
	for i, r := range results {
		out[i] = *fromRepoDevice(&r)
	}
	return out, nil
}

// EnrollDeviceInCampaign adds a campaign association to a device.
// Op #18: FR-017, FR-018
func (o *Ops) EnrollDeviceInCampaign(ctx context.Context, deviceID string, campaignID string) error {
	return o.repo.EnrollInCampaign(ctx, deviceID, campaignID)
}

func toRepoCreateDeviceInput(in CreateDeviceInput) devicerepo.CreateDeviceInput {
	return devicerepo.CreateDeviceInput{
		OwnerID:         in.OwnerID,
		Class:           in.Class,
		FirmwareVersion: in.FirmwareVersion,
		Tier:            in.Tier,
		Sensors:         in.Sensors,
	}
}

func toRepoQueryByClassInput(in QueryByClassInput) devicerepo.QueryByClassInput {
	return devicerepo.QueryByClassInput{
		Class:          in.Class,
		FirmwareMinGte: in.FirmwareMinGte,
		FirmwareMaxLte: in.FirmwareMaxLte,
	}
}

func toRepoGenerateCodeInput(in GenerateCodeInput) devicerepo.GenerateCodeInput {
	return devicerepo.GenerateCodeInput{
		DeviceID: in.DeviceID,
		Code:     in.Code,
		TTL:      in.TTL,
	}
}

func fromRepoDevice(r *devicerepo.Device) *Device {
	return &Device{
		ID:              r.ID,
		OwnerID:         r.OwnerID,
		Status:          r.Status,
		Class:           r.Class,
		FirmwareVersion: r.FirmwareVersion,
		Tier:            r.Tier,
		Sensors:         r.Sensors,
		CertSerial:      r.CertSerial,
		CreatedAt:       r.CreatedAt,
	}
}

func fromRepoEnrollmentCode(r *devicerepo.EnrollmentCode) *EnrollmentCode {
	return &EnrollmentCode{
		Code:      r.Code,
		DeviceID:  r.DeviceID,
		ExpiresAt: r.ExpiresAt,
		Used:      r.Used,
	}
}
