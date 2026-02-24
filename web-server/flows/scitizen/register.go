package scitizen

import (
	"context"

	scitizenops "rootstock/web-server/ops/scitizen"
)

// ScitzenRegistrationFlow handles scitizen account registration.
// Graph node: 0x2f — implements FR-011 (0x1), FR-080 (0x11)
type ScitzenRegistrationFlow struct {
	scitizenOps *scitizenops.Ops
}

// NewScitzenRegistrationFlow creates the flow with its required ops.
func NewScitzenRegistrationFlow(scitizenOps *scitizenops.Ops) *ScitzenRegistrationFlow {
	return &ScitzenRegistrationFlow{scitizenOps: scitizenOps}
}

// Run registers a scitizen and creates their profile with ToS acceptance.
func (f *ScitzenRegistrationFlow) Run(ctx context.Context, input RegisterInput) (*Profile, error) {
	result, err := f.scitizenOps.CreateProfile(ctx, scitizenops.CreateProfileInput{
		UserID:     input.UserID,
		TOSVersion: input.TOSVersion,
	})
	if err != nil {
		return nil, err
	}

	return &Profile{
		UserID:           result.UserID,
		TOSAccepted:      result.TOSAccepted,
		DeviceRegistered: result.DeviceRegistered,
		CampaignEnrolled: result.CampaignEnrolled,
		FirstReading:     result.FirstReading,
	}, nil
}
