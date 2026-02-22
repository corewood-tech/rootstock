package graph

import (
	"context"

	graphrepo "rootstock/web-server/repo/graph"
)

// Ops holds graph operations. Each method is one op.
type Ops struct {
	repo graphrepo.Repository
}

// NewOps creates graph ops backed by the given repository.
func NewOps(repo graphrepo.Repository) *Ops {
	return &Ops{repo: repo}
}

// InitCampaignState creates a campaign instance in the "draft" state.
func (o *Ops) InitCampaignState(ctx context.Context, campaignRef string) (*CampaignState, error) {
	result, err := o.repo.InitCampaignState(ctx, campaignRef)
	if err != nil {
		return nil, err
	}
	return fromRepoCampaignState(result), nil
}

// TransitionState attempts a state transition for a campaign.
func (o *Ops) TransitionState(ctx context.Context, input TransitionStateInput) (*CampaignState, error) {
	result, err := o.repo.TransitionState(ctx, toRepoTransitionInput(input))
	if err != nil {
		return nil, err
	}
	return fromRepoCampaignState(result), nil
}

// GetValidTransitions returns the set of valid transitions from the campaign's current state.
func (o *Ops) GetValidTransitions(ctx context.Context, campaignRef string) ([]Transition, error) {
	results, err := o.repo.GetValidTransitions(ctx, campaignRef)
	if err != nil {
		return nil, err
	}
	out := make([]Transition, len(results))
	for i, t := range results {
		out[i] = fromRepoTransition(&t)
	}
	return out, nil
}

// GetCurrentState returns the current lifecycle state for a campaign.
func (o *Ops) GetCurrentState(ctx context.Context, campaignRef string) (*CampaignState, error) {
	result, err := o.repo.GetCurrentState(ctx, campaignRef)
	if err != nil {
		return nil, err
	}
	return fromRepoCampaignState(result), nil
}

// UpdateBaseline applies a new reading value to the rolling statistics.
func (o *Ops) UpdateBaseline(ctx context.Context, input UpdateBaselineInput) (*Baseline, error) {
	result, err := o.repo.UpdateBaseline(ctx, toRepoUpdateBaselineInput(input))
	if err != nil {
		return nil, err
	}
	return fromRepoBaseline(result), nil
}

// CheckAnomaly evaluates a reading value against the baseline bounds.
// Returns nil if within bounds.
func (o *Ops) CheckAnomaly(ctx context.Context, input CheckAnomalyInput) (*AnomalyResult, error) {
	result, err := o.repo.CheckAnomaly(ctx, toRepoCheckAnomalyInput(input))
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return fromRepoAnomalyFlag(result), nil
}

// AddEnrollment creates a device-campaign enrollment edge in the graph.
func (o *Ops) AddEnrollment(ctx context.Context, input AddEnrollmentInput) error {
	return o.repo.AddEnrollment(ctx, toRepoEnrollmentInput(input))
}

// WithdrawEnrollment marks an enrollment as withdrawn.
func (o *Ops) WithdrawEnrollment(ctx context.Context, deviceRef, campaignRef string) error {
	return o.repo.WithdrawEnrollment(ctx, deviceRef, campaignRef)
}

// GetDeviceCampaigns returns all campaigns a device is enrolled in.
func (o *Ops) GetDeviceCampaigns(ctx context.Context, deviceRef string) ([]Enrollment, error) {
	results, err := o.repo.GetDeviceCampaigns(ctx, deviceRef)
	if err != nil {
		return nil, err
	}
	out := make([]Enrollment, len(results))
	for i, e := range results {
		out[i] = fromRepoEnrollment(&e)
	}
	return out, nil
}

// GetCampaignDevices returns all devices enrolled in a campaign.
func (o *Ops) GetCampaignDevices(ctx context.Context, campaignRef string) ([]Enrollment, error) {
	results, err := o.repo.GetCampaignDevices(ctx, campaignRef)
	if err != nil {
		return nil, err
	}
	out := make([]Enrollment, len(results))
	for i, e := range results {
		out[i] = fromRepoEnrollment(&e)
	}
	return out, nil
}

// --- toRepo converters ---

func toRepoTransitionInput(in TransitionStateInput) graphrepo.TransitionInput {
	return graphrepo.TransitionInput{
		CampaignRef: in.CampaignRef,
		EventName:   in.EventName,
	}
}

func toRepoUpdateBaselineInput(in UpdateBaselineInput) graphrepo.UpdateBaselineInput {
	return graphrepo.UpdateBaselineInput{
		CampaignRef:   in.CampaignRef,
		ParameterName: in.ParameterName,
		Value:         in.Value,
	}
}

func toRepoCheckAnomalyInput(in CheckAnomalyInput) graphrepo.CheckAnomalyInput {
	return graphrepo.CheckAnomalyInput{
		CampaignRef:   in.CampaignRef,
		ParameterName: in.ParameterName,
		Value:         in.Value,
	}
}

func toRepoEnrollmentInput(in AddEnrollmentInput) graphrepo.EnrollmentInput {
	return graphrepo.EnrollmentInput{
		DeviceRef:   in.DeviceRef,
		OwnerRef:    in.OwnerRef,
		CampaignRef: in.CampaignRef,
		EnrolledAt:  in.EnrolledAt,
	}
}

// --- fromRepo converters ---

func fromRepoCampaignState(r *graphrepo.CampaignInstanceState) *CampaignState {
	return &CampaignState{
		CampaignRef: r.CampaignRef,
		StateName:   r.StateName,
		EnteredAt:   r.EnteredAt,
	}
}

func fromRepoTransition(t *graphrepo.ValidTransition) Transition {
	return Transition{
		EventName:   t.EventName,
		TargetState: t.TargetState,
		Guard:       t.Guard,
		SideEffect:  t.SideEffect,
	}
}

func fromRepoBaseline(b *graphrepo.AnomalyBaseline) *Baseline {
	return &Baseline{
		CampaignRef:      b.CampaignRef,
		ParameterName:    b.ParameterName,
		SampleCount:      b.SampleCount,
		Mean:             b.Mean,
		M2:               b.M2,
		Min:              b.Min,
		Max:              b.Max,
		StddevMultiplier: b.StddevMultiplier,
		HardMin:          b.HardMin,
		HardMax:          b.HardMax,
		LastUpdated:      b.LastUpdated,
	}
}

func fromRepoAnomalyFlag(f *graphrepo.AnomalyFlag) *AnomalyResult {
	return &AnomalyResult{
		Reason:     f.Reason,
		Value:      f.Value,
		LowerBound: f.LowerBound,
		UpperBound: f.UpperBound,
		Mean:       f.Mean,
		Stddev:     f.Stddev,
	}
}

func fromRepoEnrollment(e *graphrepo.EnrollmentEdge) Enrollment {
	return Enrollment{
		DeviceRef:        e.DeviceRef,
		CampaignRef:      e.CampaignRef,
		OwnerRef:         e.OwnerRef,
		EnrolledAt:       e.EnrolledAt,
		WithdrawnAt:      e.WithdrawnAt,
		EnrollmentStatus: e.EnrollmentStatus,
	}
}
