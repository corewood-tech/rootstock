package scitizen

import (
	"context"

	scitizenops "rootstock/web-server/ops/scitizen"
)

// OnboardingFlow tracks onboarding progress for scitizens.
// Graph node: 0x26 — implements US-001 (0xd), FR-104 (0x18)
type OnboardingFlow struct {
	scitizenOps *scitizenops.Ops
}

// NewOnboardingFlow creates the flow with its required ops.
func NewOnboardingFlow(scitizenOps *scitizenops.Ops) *OnboardingFlow {
	return &OnboardingFlow{scitizenOps: scitizenOps}
}

// Run returns the current onboarding state.
func (f *OnboardingFlow) Run(ctx context.Context, userID string) (*OnboardingState, error) {
	profile, err := f.scitizenOps.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	if profile == nil {
		return &OnboardingState{}, nil
	}

	return &OnboardingState{
		DeviceRegistered:     profile.DeviceRegistered,
		CampaignEnrolled:     profile.CampaignEnrolled,
		FirstReadingSubmitted: profile.FirstReading,
		TOSAccepted:          profile.TOSAccepted,
	}, nil
}
