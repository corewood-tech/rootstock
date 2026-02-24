package scitizen

import "time"

// Profile is the scitizen profile returned by flows.
type Profile struct {
	UserID           string
	TOSAccepted      bool
	DeviceRegistered bool
	CampaignEnrolled bool
	FirstReading     bool
}

// OnboardingState is the onboarding progress.
type OnboardingState struct {
	DeviceRegistered     bool
	CampaignEnrolled     bool
	FirstReadingSubmitted bool
	TOSAccepted          bool
}

// Dashboard aggregates scitizen dashboard data.
type Dashboard struct {
	ActiveEnrollments int
	TotalReadings     int
	AcceptedReadings  int
	ContributionScore float64
	Badges            []Badge
	Enrollments       []Enrollment
}

// Badge represents an awarded badge.
type Badge struct {
	ID        string
	BadgeType string
	AwardedAt time.Time
}

// Enrollment is a campaign enrollment summary.
type Enrollment struct {
	ID         string
	DeviceID   string
	CampaignID string
	Status     string
	EnrolledAt time.Time
}

// CampaignSummary is a browsable campaign entry.
type CampaignSummary struct {
	ID              string
	Status          string
	WindowStart     *time.Time
	WindowEnd       *time.Time
	EnrollmentCount int
	RequiredSensors []string
	CreatedAt       time.Time
}

// CampaignDetail is the full campaign detail.
type CampaignDetail struct {
	CampaignID      string
	Status          string
	WindowStart     *time.Time
	WindowEnd       *time.Time
	Parameters      []Parameter
	Regions         []Region
	Eligibility     []EligibilityCriteria
	EnrollmentCount int
	ProgressPercent float64
}

// Parameter mirrors campaign parameter data.
type Parameter struct {
	Name      string
	Unit      string
	MinRange  *float64
	MaxRange  *float64
	Precision *int
}

// Region mirrors campaign region data.
type Region struct {
	GeoJSON string
}

// EligibilityCriteria mirrors campaign eligibility data.
type EligibilityCriteria struct {
	DeviceClass     string
	Tier            int
	RequiredSensors []string
	FirmwareMin     string
}

// EnrollResult is the result of an enrollment attempt.
type EnrollResult struct {
	Enrolled     bool
	Reason       string
	EnrollmentID string
}

// DeviceSummary is a device list entry.
type DeviceSummary struct {
	ID                string
	Status            string
	Class             string
	FirmwareVersion   string
	Tier              int
	Sensors           []string
	ActiveEnrollments int
	LastSeen          *time.Time
}

// DeviceDetail is full device info.
type DeviceDetail struct {
	ID                string
	OwnerID           string
	Status            string
	Class             string
	FirmwareVersion   string
	Tier              int
	Sensors           []string
	CertSerial        *string
	CreatedAt         time.Time
	Enrollments       []Enrollment
	ConnectionHistory []ConnectionEvent
}

// ConnectionEvent is a device connection history entry.
type ConnectionEvent struct {
	EventType string
	Timestamp time.Time
	Reason    *string
}

// Notification is an in-app notification record.
type Notification struct {
	ID           string
	UserID       string
	Type         string
	Message      string
	Read         bool
	ResourceLink *string
	CreatedAt    time.Time
}

// NotificationsResult is the result of a notifications query.
type NotificationsResult struct {
	Notifications []Notification
	UnreadCount   int
	Total         int
}

// ReadingHistory is per-device per-campaign reading stats.
type ReadingHistory struct {
	DeviceID   string
	CampaignID string
	Total      int
	Accepted   int
	Rejected   int
}

// ContributionsResult is the result of a contributions query.
type ContributionsResult struct {
	Histories         []ReadingHistory
	ContributionScore float64
	Badges            []Badge
}
