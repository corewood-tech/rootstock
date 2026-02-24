package scitizen

// CreateProfileInput is what the op sends to create a scitizen profile.
type CreateProfileInput struct {
	UserID     string
	TOSVersion string
}

// UpdateOnboardingInput updates onboarding state flags.
type UpdateOnboardingInput struct {
	UserID            string
	DeviceRegistered  *bool
	CampaignEnrolled  *bool
	FirstReading      *bool
}

// BrowseInput filters published campaigns for scitizen browsing.
type BrowseInput struct {
	Longitude  *float64
	Latitude   *float64
	RadiusKm   *float64
	SensorType *string
	Limit      int
	Offset     int
}

// SearchInput is full-text search across published campaigns.
type SearchInput struct {
	Query  string
	Limit  int
	Offset int
}

// GetNotificationsInput filters notifications for a user.
type GetNotificationsInput struct {
	UserID     string
	TypeFilter *string
	Limit      int
	Offset     int
}
