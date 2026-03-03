package campaign

import "time"

// Campaign is the campaign record returned by campaign flows.
type Campaign struct {
	ID          string
	OrgID       string
	Status      string
	WindowStart *time.Time
	WindowEnd   *time.Time
	CreatedBy   string
	CreatedAt   time.Time
}

// ParameterQualityItem holds per-parameter quality metrics.
type ParameterQualityItem struct {
	ParameterName    string
	AcceptedCount    int
	QuarantinedCount int
}

// DeviceBreakdownItem holds per-device stats.
type DeviceBreakdownItem struct {
	PseudoDeviceID string
	DeviceClass    string
	AcceptanceRate float64
	ReadingCount   int
	LastSeen       *string
}

// EnrollmentFunnelItem holds enrollment stage counts.
type EnrollmentFunnelItem struct {
	Enrolled     int
	Active       int
	Contributing int
}

// TemporalBucketItem holds reading counts for a time bucket.
type TemporalBucketItem struct {
	Bucket string
	Count  int
}

// Dashboard is the dashboard data returned by CampaignDashboardFlow.
type Dashboard struct {
	CampaignID       string
	AcceptedCount    int
	QuarantineCount  int
	ParameterQuality []ParameterQualityItem
	DeviceBreakdown  []DeviceBreakdownItem
	EnrollmentFunnel EnrollmentFunnelItem
	TemporalCoverage []TemporalBucketItem
}
