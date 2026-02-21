package graph

import "context"

// Repository defines the interface for runtime graph operations
// backed by Dgraph. Covers campaign state machines, anomaly baselines,
// and device-campaign relationship traversals.
type Repository interface {
	// --- Campaign State Machine ---

	// GetCurrentState returns the current lifecycle state for a campaign.
	GetCurrentState(ctx context.Context, campaignRef string) (*CampaignInstanceState, error)

	// TransitionState attempts a state transition. Returns the new state
	// or an error if the transition is invalid (no matching edge or guard fails).
	TransitionState(ctx context.Context, input TransitionInput) (*CampaignInstanceState, error)

	// GetValidTransitions returns the set of valid transitions from the
	// campaign's current state.
	GetValidTransitions(ctx context.Context, campaignRef string) ([]ValidTransition, error)

	// InitCampaignState creates a CampaignInstance node in the "draft" state.
	InitCampaignState(ctx context.Context, campaignRef string) (*CampaignInstanceState, error)

	// --- Anomaly Detection ---

	// GetBaseline returns the rolling statistics for a campaign parameter.
	GetBaseline(ctx context.Context, campaignRef string, parameterName string) (*AnomalyBaseline, error)

	// UpdateBaseline applies Welford's online algorithm to update rolling
	// statistics with a new accepted reading value.
	UpdateBaseline(ctx context.Context, input UpdateBaselineInput) (*AnomalyBaseline, error)

	// CheckAnomaly evaluates a reading value against the baseline bounds.
	// Returns nil if within bounds, or an AnomalyFlag if outside.
	CheckAnomaly(ctx context.Context, input CheckAnomalyInput) (*AnomalyFlag, error)

	// --- Device-Campaign Relationships ---

	// AddEnrollment creates a device→campaign enrollment edge.
	AddEnrollment(ctx context.Context, input EnrollmentInput) error

	// WithdrawEnrollment marks an enrollment as withdrawn.
	WithdrawEnrollment(ctx context.Context, deviceRef string, campaignRef string) error

	// GetDeviceCampaigns returns all campaigns a device is enrolled in.
	GetDeviceCampaigns(ctx context.Context, deviceRef string) ([]EnrollmentEdge, error)

	// GetCampaignDevices returns all devices enrolled in a campaign.
	GetCampaignDevices(ctx context.Context, campaignRef string) ([]EnrollmentEdge, error)

	// GetSharedDeviceCampaigns returns campaigns that share devices with
	// the given campaign (for cross-campaign analysis).
	GetSharedDeviceCampaigns(ctx context.Context, campaignRef string) ([]string, error)

	// Shutdown gracefully terminates the graph repository.
	Shutdown()
}
