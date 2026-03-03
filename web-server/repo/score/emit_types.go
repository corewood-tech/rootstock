package score

import "time"

// Score is the contribution score record.
type Score struct {
	ScitizenID  string
	Volume      int
	QualityRate float64
	Consistency float64
	Diversity   int
	Total       float64
	UpdatedAt   time.Time
}

// Badge is an awarded badge record.
type Badge struct {
	ID         string
	ScitizenID string
	BadgeType  string
	AwardedAt  time.Time
}

// SweepstakesEntry is a sweepstakes entry record.
type SweepstakesEntry struct {
	ID               string
	ScitizenID       string
	Entries          int
	MilestoneTrigger string
	GrantedAt        time.Time
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
