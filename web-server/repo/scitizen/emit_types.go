package scitizen

import "time"

// Profile is the scitizen_profiles record.
type Profile struct {
	UserID           string
	TOSAccepted      bool
	TOSVersion       *string
	TOSAcceptedAt    *time.Time
	DeviceRegistered bool
	CampaignEnrolled bool
	FirstReading     bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
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

// ReadingHistory is per-device per-campaign reading stats.
type ReadingHistory struct {
	DeviceID   string
	CampaignID string
	Total      int
	Accepted   int
	Rejected   int
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

// CampaignDetail is the full campaign detail for enrollment decision.
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

// DeviceSummary is a device list entry for scitizen.
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

// DeviceDetail is full device info with enrollments and connection history.
type DeviceDetail struct {
	ID              string
	OwnerID         string
	Status          string
	Class           string
	FirmwareVersion string
	Tier            int
	Sensors         []string
	CertSerial      *string
	CreatedAt       time.Time
	Enrollments     []Enrollment
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
