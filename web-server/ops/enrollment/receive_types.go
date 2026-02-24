package enrollment

// EnrollInput is what callers send to Enroll.
type EnrollInput struct {
	DeviceID       string
	CampaignID     string
	ScitizenID     string
	ConsentVersion string
	ConsentScope   string
}

// CreateNotificationInput creates a notification.
type CreateNotificationInput struct {
	UserID       string
	Type         string
	Message      string
	ResourceLink *string
}
