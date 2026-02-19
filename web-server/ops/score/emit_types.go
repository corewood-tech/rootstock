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
	ID        string
	BadgeType string
	AwardedAt time.Time
}
