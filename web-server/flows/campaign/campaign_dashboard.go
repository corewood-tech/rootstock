package campaign

import (
	"context"

	readingops "rootstock/web-server/ops/reading"
)

// DashboardFlow returns enriched campaign dashboard metrics.
type DashboardFlow struct {
	readingOps *readingops.Ops
	hmacSecret string
}

// NewDashboardFlow creates the flow with its required ops.
func NewDashboardFlow(readingOps *readingops.Ops, hmacSecret string) *DashboardFlow {
	return &DashboardFlow{readingOps: readingOps, hmacSecret: hmacSecret}
}

// Run fetches enriched campaign dashboard metrics.
func (f *DashboardFlow) Run(ctx context.Context, campaignID string) (*Dashboard, error) {
	metrics, err := f.readingOps.GetCampaignQuality(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	d := &Dashboard{
		CampaignID:      metrics.CampaignID,
		AcceptedCount:   metrics.AcceptedCount,
		QuarantineCount: metrics.QuarantineCount,
	}

	// Per-parameter quality
	for _, pq := range metrics.PerParameter {
		d.ParameterQuality = append(d.ParameterQuality, ParameterQualityItem{
			ParameterName:    pq.ParameterName,
			AcceptedCount:    pq.AcceptedCount,
			QuarantinedCount: pq.QuarantinedCount,
		})
	}

	// Device breakdown
	devices, err := f.readingOps.GetCampaignDeviceBreakdown(ctx, campaignID, f.hmacSecret)
	if err != nil {
		return nil, err
	}
	for _, db := range devices {
		d.DeviceBreakdown = append(d.DeviceBreakdown, DeviceBreakdownItem{
			PseudoDeviceID: db.PseudoDeviceID,
			DeviceClass:    db.DeviceClass,
			AcceptanceRate: db.AcceptanceRate,
			ReadingCount:   db.ReadingCount,
			LastSeen:       db.LastSeen,
		})
	}

	// Enrollment funnel
	funnel, err := f.readingOps.GetEnrollmentFunnel(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	d.EnrollmentFunnel = EnrollmentFunnelItem{
		Enrolled:     funnel.Enrolled,
		Active:       funnel.Active,
		Contributing: funnel.Contributing,
	}

	// Temporal coverage
	temporal, err := f.readingOps.GetCampaignTemporalCoverage(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	for _, t := range temporal {
		d.TemporalCoverage = append(d.TemporalCoverage, TemporalBucketItem{
			Bucket: t.Bucket,
			Count:  t.Count,
		})
	}

	return d, nil
}
