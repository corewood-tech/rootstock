package scitizen

// CreateProfileInput is what callers send to CreateProfile.
type CreateProfileInput struct {
	UserID     string
	TOSVersion string
}

// UpdateOnboardingInput updates onboarding flags.
type UpdateOnboardingInput struct {
	UserID           string
	DeviceRegistered *bool
	CampaignEnrolled *bool
	FirstReading     *bool
}

// BrowseInput filters published campaigns.
type BrowseInput struct {
	Longitude  *float64
	Latitude   *float64
	RadiusKm   *float64
	SensorType *string
	Limit      int
	Offset     int
}

// SearchInput is full-text search across campaigns.
type SearchInput struct {
	Query  string
	Limit  int
	Offset int
}

// GetNotificationsInput filters notifications.
type GetNotificationsInput struct {
	UserID     string
	TypeFilter *string
	Limit      int
	Offset     int
}
