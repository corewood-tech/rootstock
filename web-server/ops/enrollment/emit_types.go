package enrollment

import "time"

// Enrollment is the enrollment record returned by ops.
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
