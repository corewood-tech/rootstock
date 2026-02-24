package scitizen

import (
	"context"
	"fmt"

	enrollmentops "rootstock/web-server/ops/enrollment"
	scitizenops "rootstock/web-server/ops/scitizen"
)

// EnrollDeviceCampaignFlow handles device enrollment in a campaign.
// Graph node: 0x28 — implements FR-016 (0x4), FR-017 (0x1a), FR-019 (0x1b),
//   FR-013 (0x20), FR-064 (0x21)
type EnrollDeviceCampaignFlow struct {
	scitizenOps   *scitizenops.Ops
	enrollmentOps *enrollmentops.Ops
}

// NewEnrollDeviceCampaignFlow creates the flow with its required ops.
func NewEnrollDeviceCampaignFlow(scitizenOps *scitizenops.Ops, enrollmentOps *enrollmentops.Ops) *EnrollDeviceCampaignFlow {
	return &EnrollDeviceCampaignFlow{scitizenOps: scitizenOps, enrollmentOps: enrollmentOps}
}

// Run enrolls a device in a campaign after eligibility check and consent capture.
func (f *EnrollDeviceCampaignFlow) Run(ctx context.Context, input EnrollDeviceInput) (*EnrollResult, error) {
	// Check if already enrolled
	existing, err := f.enrollmentOps.GetByDeviceCampaign(ctx, input.DeviceID, input.CampaignID)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.Status == "active" {
		return &EnrollResult{Enrolled: false, Reason: "device already enrolled in this campaign"}, nil
	}

	// Get campaign detail for eligibility check
	detail, err := f.scitizenOps.GetCampaignDetail(ctx, input.CampaignID)
	if err != nil {
		return nil, fmt.Errorf("get campaign detail: %w", err)
	}
	if detail.Status != "published" {
		return &EnrollResult{Enrolled: false, Reason: "campaign is not published"}, nil
	}

	// Get device to check eligibility
	device, err := f.scitizenOps.GetDeviceDetail(ctx, input.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("get device: %w", err)
	}
	if device.Status != "active" {
		return &EnrollResult{Enrolled: false, Reason: "device is not active"}, nil
	}

	// Check eligibility criteria
	if len(detail.Eligibility) > 0 {
		eligible := false
		for _, e := range detail.Eligibility {
			if device.Class == e.DeviceClass && device.Tier >= e.Tier {
				eligible = true
				break
			}
		}
		if !eligible {
			return &EnrollResult{Enrolled: false, Reason: "device does not meet eligibility criteria"}, nil
		}
	}

	// Enroll with consent
	enrollment, err := f.enrollmentOps.Enroll(ctx, enrollmentops.EnrollInput{
		DeviceID:       input.DeviceID,
		CampaignID:     input.CampaignID,
		ScitizenID:     input.ScitizenID,
		ConsentVersion: input.ConsentVersion,
		ConsentScope:   input.ConsentScope,
	})
	if err != nil {
		return nil, fmt.Errorf("enroll device: %w", err)
	}

	// Update onboarding state (best-effort)
	t := true
	_ = f.scitizenOps.UpdateOnboarding(ctx, scitizenops.UpdateOnboardingInput{
		UserID:           input.ScitizenID,
		CampaignEnrolled: &t,
	})

	return &EnrollResult{
		Enrolled:     true,
		Reason:       "enrolled successfully",
		EnrollmentID: enrollment.ID,
	}, nil
}
