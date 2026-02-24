package enrollment

import "time"

// Enrollment is the campaign_enrollments record.
type Enrollment struct {
	ID          string
	DeviceID    string
	CampaignID  string
	ScitizenID  string
	Status      string
	EnrolledAt  time.Time
	WithdrawnAt *time.Time
}

// NotificationPreference is a per-type notification preference.
type NotificationPreference struct {
	Type  string
	InApp bool
	Email bool
}
