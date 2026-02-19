package score

// UpsertScoreInput is what callers send to UpdateScore.
type UpsertScoreInput struct {
	ScitizenID  string
	Volume      int
	QualityRate float64
	Consistency float64
	Diversity   int
	Total       float64
}

// GrantSweepstakesInput is what callers send to GrantSweepstakes.
type GrantSweepstakesInput struct {
	ScitizenID       string
	Entries          int
	MilestoneTrigger string
}
