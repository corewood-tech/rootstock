package score

import "time"

// Contribution is the scitizen's contribution profile returned by GetContribution.
type Contribution struct {
	ScitizenID  string
	Volume      int
	QualityRate float64
	Consistency float64
	Diversity   int
	Total       float64
	UpdatedAt   time.Time
	Badges      []Badge
}

// Badge is an awarded badge.
type Badge struct {
	ID        string
	BadgeType string
	AwardedAt time.Time
}

// UpdateContributionScoreResult is the result of the UpdateContributionScore flow.
type UpdateContributionScoreResult struct {
	Score         Score
	BadgesAwarded []string
	SweepEntries  int
}

// Score is a contribution score snapshot.
type Score struct {
	ScitizenID  string
	Volume      int
	QualityRate float64
	Consistency float64
	Diversity   int
	Total       float64
	UpdatedAt   time.Time
}

// LeaderboardEntry is a single row in the leaderboard.
type LeaderboardEntry struct {
	Rank          int
	ScitizenID    string
	Score         float64
	BadgeCount    int
	CampaignCount int
}

// LeaderboardResult holds the full leaderboard response.
type LeaderboardResult struct {
	Entries   []LeaderboardEntry
	Total     int
	Requester *LeaderboardEntry
}
