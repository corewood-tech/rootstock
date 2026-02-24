package enrollment

import "context"

// Repository defines the interface for enrollment data operations.
// Graph node: 0x2c (EnrollmentRepository)
type Repository interface {
	Enroll(ctx context.Context, input EnrollInput) (*Enrollment, error)
	Withdraw(ctx context.Context, enrollmentID string) error
	GetByID(ctx context.Context, id string) (*Enrollment, error)
	GetByDeviceCampaign(ctx context.Context, deviceID, campaignID string) (*Enrollment, error)
	MarkRead(ctx context.Context, userID string, ids []string) (int, error)
	CreateNotification(ctx context.Context, input CreateNotificationInput) error
	GetPreferences(ctx context.Context, userID string) ([]NotificationPreference, error)
	UpdatePreferences(ctx context.Context, userID string, prefs []NotificationPreference) error
	Shutdown()
}
