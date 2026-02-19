package reading

import "context"

// Repository defines the interface for reading data operations.
type Repository interface {
	Persist(ctx context.Context, input PersistReadingInput) (*Reading, error)
	Quarantine(ctx context.Context, id string, reason string) error
	Query(ctx context.Context, input QueryReadingsInput) ([]Reading, error)
	QuarantineByWindow(ctx context.Context, input QuarantineByWindowInput) (int64, error)
	GetCampaignQuality(ctx context.Context, campaignID string) (*QualityMetrics, error)
	Shutdown()
}
