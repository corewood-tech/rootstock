package graph

import "time"

// CampaignState represents the current lifecycle state of a campaign.
type CampaignState struct {
	CampaignRef string
	StateName   string
	EnteredAt   time.Time
}

// Transition describes a possible state transition from the current state.
type Transition struct {
	EventName   string
	TargetState string
	Guard       string
	SideEffect  string
}

// Baseline holds the rolling statistics for a campaign parameter.
type Baseline struct {
	CampaignRef      string
	ParameterName    string
	SampleCount      int64
	Mean             float64
	M2               float64
	Min              float64
	Max              float64
	StddevMultiplier float64
	HardMin          float64
	HardMax          float64
	LastUpdated      time.Time
}

// AnomalyResult indicates a reading fell outside baseline bounds.
type AnomalyResult struct {
	Reason     string
	Value      float64
	LowerBound float64
	UpperBound float64
	Mean       float64
	Stddev     float64
}

// Enrollment represents a device-campaign relationship.
type Enrollment struct {
	DeviceRef        string
	CampaignRef      string
	OwnerRef         string
	EnrolledAt       time.Time
	WithdrawnAt      *time.Time
	EnrollmentStatus string
}
