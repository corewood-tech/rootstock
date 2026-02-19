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
