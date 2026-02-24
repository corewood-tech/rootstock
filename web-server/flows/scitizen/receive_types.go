package scitizen

// RegisterInput is what callers send to ScitzenRegistrationFlow.
type RegisterInput struct {
	UserID     string
	TOSVersion string
}

// BrowseInput filters published campaigns for browsing.
type BrowseInput struct {
	Longitude  *float64
	Latitude   *float64
	RadiusKm   *float64
	SensorType *string
	Limit      int
	Offset     int
}

// SearchInput is full-text campaign search.
type SearchInput struct {
	Query  string
	Limit  int
	Offset int
}

// EnrollDeviceInput enrolls a device in a campaign.
type EnrollDeviceInput struct {
	ScitizenID     string
	DeviceID       string
	CampaignID     string
	ConsentVersion string
	ConsentScope   string
}

// GetNotificationsInput filters notifications.
type GetNotificationsInput struct {
	UserID     string
	TypeFilter *string
	Limit      int
	Offset     int
}
