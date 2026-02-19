package device

import "context"

// Repository defines the interface for device registry operations.
type Repository interface {
	Create(ctx context.Context, input CreateDeviceInput) (*Device, error)
	Get(ctx context.Context, id string) (*Device, error)
	GetCapabilities(ctx context.Context, id string) (*DeviceCapabilities, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	QueryByClass(ctx context.Context, input QueryByClassInput) ([]Device, error)
	GenerateEnrollmentCode(ctx context.Context, input GenerateCodeInput) (*EnrollmentCode, error)
	RedeemEnrollmentCode(ctx context.Context, code string) (*EnrollmentCode, error)
	EnrollInCampaign(ctx context.Context, deviceID string, campaignID string) error
	Shutdown()
}
