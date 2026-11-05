package potager

// Season returns the season a planting belongs to: the year it ended.
// Only meaningful for an ended planting: a zero Date gives 1.
func Season(endedAt Date) int {
	return endedAt.t.Year()
}
