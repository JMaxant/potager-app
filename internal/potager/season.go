package potager

// Season returns the season a planting belongs to: the year it ended.
func Season(endedAt Date) int {
	return endedAt.t.Year()
}
