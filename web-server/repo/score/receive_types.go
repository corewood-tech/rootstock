package score

// UpsertScoreInput is what the UpdateScore op sends to the repository.
type UpsertScoreInput struct {
	ScitizenID  string
	Volume      int
	QualityRate float64
	Consistency float64
	Diversity   int
	Total       float64
}

// GrantSweepstakesInput is what the GrantSweepstakes op sends.
type GrantSweepstakesInput struct {
	ScitizenID       string
	Entries          int
	MilestoneTrigger string
}
