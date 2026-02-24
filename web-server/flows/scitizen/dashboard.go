package scitizen

import (
	"context"

	scitizenops "rootstock/web-server/ops/scitizen"
)

// ScitzenDashboardFlow aggregates scitizen dashboard data.
// Graph node: 0x29 — implements US-003 (0x7), FR-081 (0x9), FR-040 (0xb), FR-095 (0xf)
type ScitzenDashboardFlow struct {
	scitizenOps *scitizenops.Ops
}

// NewScitzenDashboardFlow creates the flow with its required ops.
func NewScitzenDashboardFlow(scitizenOps *scitizenops.Ops) *ScitzenDashboardFlow {
	return &ScitzenDashboardFlow{scitizenOps: scitizenOps}
}

// Run returns aggregated dashboard data for the scitizen.
func (f *ScitzenDashboardFlow) Run(ctx context.Context, userID string) (*Dashboard, error) {
	result, err := f.scitizenOps.GetDashboard(ctx, userID)
	if err != nil {
		return nil, err
	}

	badges := make([]Badge, len(result.Badges))
	for i, b := range result.Badges {
		badges[i] = Badge{ID: b.ID, BadgeType: b.BadgeType, AwardedAt: b.AwardedAt}
	}
	enrollments := make([]Enrollment, len(result.Enrollments))
	for i, e := range result.Enrollments {
		enrollments[i] = Enrollment{
			ID: e.ID, DeviceID: e.DeviceID, CampaignID: e.CampaignID,
			Status: e.Status, EnrolledAt: e.EnrolledAt,
		}
	}

	return &Dashboard{
		ActiveEnrollments: result.ActiveEnrollments,
		TotalReadings:     result.TotalReadings,
		AcceptedReadings:  result.AcceptedReadings,
		ContributionScore: result.ContributionScore,
		Badges:            badges,
		Enrollments:       enrollments,
	}, nil
}
