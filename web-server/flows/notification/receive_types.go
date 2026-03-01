package notification

// ListInput is what callers send to ListNotificationsFlow.
type ListInput struct {
	UserID     string
	TypeFilter *string
	Limit      int
	Offset     int
}

// MarkReadInput is what callers send to MarkReadFlow.
type MarkReadInput struct {
	UserID          string
	NotificationIDs []string
}

// UpdatePreferencesInput is what callers send to UpdatePreferencesFlow.
type UpdatePreferencesInput struct {
	UserID      string
	Preferences []PreferenceInput
}

// PreferenceInput is a single preference update.
type PreferenceInput struct {
	Type  string
	InApp bool
	Email bool
}
