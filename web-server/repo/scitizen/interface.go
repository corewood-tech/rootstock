package scitizen

import "context"

// Repository defines the interface for scitizen data operations.
// Graph node: 0x24 (ScitizenRepository)
type Repository interface {
	CreateProfile(ctx context.Context, input CreateProfileInput) (*Profile, error)
	GetProfile(ctx context.Context, userID string) (*Profile, error)
	UpdateOnboarding(ctx context.Context, input UpdateOnboardingInput) error
	GetDashboard(ctx context.Context, userID string) (*Dashboard, error)
	GetContributions(ctx context.Context, userID string) ([]ReadingHistory, error)
	BrowseCampaigns(ctx context.Context, input BrowseInput) ([]CampaignSummary, int, error)
	GetCampaignDetail(ctx context.Context, campaignID string) (*CampaignDetail, error)
	SearchCampaigns(ctx context.Context, input SearchInput) ([]CampaignSummary, int, error)
	GetDevices(ctx context.Context, ownerID string) ([]DeviceSummary, error)
	GetDeviceDetail(ctx context.Context, deviceID string) (*DeviceDetail, error)
	GetNotifications(ctx context.Context, input GetNotificationsInput) ([]Notification, int, int, error)
	Shutdown()
}
