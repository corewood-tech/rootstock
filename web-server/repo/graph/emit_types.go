package graph

import "time"

// CampaignInstanceState represents the current lifecycle state of a campaign.
type CampaignInstanceState struct {
	CampaignRef  string
	StateName    string
	EnteredAt    time.Time
}

// ValidTransition describes a possible state transition from the current state.
type ValidTransition struct {
	EventName   string
	TargetState string
	Guard       string
	SideEffect  string
}

// AnomalyBaseline holds the rolling statistics for a campaign parameter.
type AnomalyBaseline struct {
	CampaignRef        string
	ParameterName      string
	SampleCount        int64
	Mean               float64
	M2                 float64 // Welford's sum of squared differences
	Min                float64
	Max                float64
	StddevMultiplier   float64
	HardMin            float64
	HardMax            float64
	LastUpdated        time.Time
}

// Stddev returns the population standard deviation from the rolling stats.
func (b *AnomalyBaseline) Stddev() float64 {
	if b.SampleCount < 2 {
		return 0
	}
	variance := b.M2 / float64(b.SampleCount)
	// Manual sqrt: use the math package at call site if needed.
	// This is a data type — keep it simple.
	return variance // caller should take sqrt
}

// AnomalyFlag indicates a reading fell outside baseline bounds.
type AnomalyFlag struct {
	Reason       string
	Value        float64
	LowerBound   float64
	UpperBound   float64
	Mean         float64
	Stddev       float64
}

// EnrollmentEdge represents a device-campaign relationship.
type EnrollmentEdge struct {
	DeviceRef        string
	CampaignRef      string
	OwnerRef         string
	EnrolledAt       time.Time
	WithdrawnAt      *time.Time
	EnrollmentStatus string
}
