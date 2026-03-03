package reading

import "context"

// Repository defines the interface for reading data operations.
type Repository interface {
	Persist(ctx context.Context, input PersistReadingInput) (*Reading, error)
	Quarantine(ctx context.Context, id string, reason string) error
	QuarantineValue(ctx context.Context, readingValueID string, reason string) error
	Query(ctx context.Context, input QueryReadingsInput) ([]Reading, error)
	QuarantineByWindow(ctx context.Context, input QuarantineByWindowInput) (int64, error)
	GetCampaignQuality(ctx context.Context, campaignID string) (*QualityMetrics, error)
	GetCampaignDeviceBreakdown(ctx context.Context, campaignID string, hmacSecret string) ([]DeviceBreakdown, error)
	GetCampaignTemporalCoverage(ctx context.Context, campaignID string) ([]TemporalBucket, error)
	GetEnrollmentFunnel(ctx context.Context, campaignID string) (*EnrollmentFunnel, error)
	GetScitizenReadingStats(ctx context.Context, scitizenID string) (*ScitizenReadingStats, error)
	Shutdown()
}
