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

// Dashboard is the dashboard data returned by CampaignDashboardFlow.
type Dashboard struct {
	CampaignID      string
	AcceptedCount   int
	QuarantineCount int
}
