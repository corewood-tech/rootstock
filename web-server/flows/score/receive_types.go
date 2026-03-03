package score

// UpdateContributionScoreInput is the input to the UpdateContributionScore flow.
type UpdateContributionScoreInput struct {
	ScitizenID  string
	Volume      int
	QualityRate float64
	Consistency float64
	Diversity   int
	Total       float64
}

// RefreshScitizenScoreInput is the input to the RefreshScitizenScore flow.
type RefreshScitizenScoreInput struct {
	DeviceID string
}

// GetLeaderboardInput controls leaderboard queries.
type GetLeaderboardInput struct {
	CampaignID  *string
	TimePeriod  string
	Limit       int
	Offset      int
	RequesterID string
}
