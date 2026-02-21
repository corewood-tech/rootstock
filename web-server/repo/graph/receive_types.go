package graph

import "time"

// TransitionInput requests a state transition for a campaign.
type TransitionInput struct {
	CampaignRef string
	EventName   string
}

// UpdateBaselineInput provides a new reading value to incorporate
// into the rolling statistics via Welford's algorithm.
type UpdateBaselineInput struct {
	CampaignRef   string
	ParameterName string
	Value         float64
}

// CheckAnomalyInput provides a reading value to check against baselines.
type CheckAnomalyInput struct {
	CampaignRef   string
	ParameterName string
	Value         float64
}

// EnrollmentInput creates a device-campaign enrollment edge.
type EnrollmentInput struct {
	DeviceRef   string
	OwnerRef    string
	CampaignRef string
	EnrolledAt  time.Time
}
