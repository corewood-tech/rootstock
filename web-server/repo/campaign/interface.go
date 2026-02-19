package campaign

import "context"

// Repository defines the interface for campaign data operations.
type Repository interface {
	Create(ctx context.Context, input CreateCampaignInput) (*Campaign, error)
	Publish(ctx context.Context, id string) error
	List(ctx context.Context, input ListCampaignsInput) ([]Campaign, error)
	GetRules(ctx context.Context, campaignID string) (*CampaignRules, error)
	GetEligibility(ctx context.Context, campaignID string) ([]EligibilityCriteria, error)
	Shutdown()
}
