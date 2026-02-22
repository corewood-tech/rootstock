package graph

import "time"

// InitCampaignStateInput is what callers send to InitCampaignState.
type InitCampaignStateInput struct {
	CampaignRef string
}

// TransitionStateInput is what callers send to TransitionState.
type TransitionStateInput struct {
	CampaignRef string
	EventName   string
}

// UpdateBaselineInput is what callers send to UpdateBaseline.
type UpdateBaselineInput struct {
	CampaignRef   string
	ParameterName string
	Value         float64
}

// CheckAnomalyInput is what callers send to CheckAnomaly.
type CheckAnomalyInput struct {
	CampaignRef   string
	ParameterName string
	Value         float64
}

// AddEnrollmentInput is what callers send to AddEnrollment.
type AddEnrollmentInput struct {
	DeviceRef   string
	OwnerRef    string
	CampaignRef string
	EnrolledAt  time.Time
}
