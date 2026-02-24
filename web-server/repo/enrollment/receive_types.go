package enrollment

// EnrollInput is what the op sends to create an enrollment.
type EnrollInput struct {
	DeviceID     string
	CampaignID   string
	ScitizenID   string
	ConsentVersion string
	ConsentScope   string
}

// CreateNotificationInput creates a new notification.
type CreateNotificationInput struct {
	UserID       string
	Type         string
	Message      string
	ResourceLink *string
}
