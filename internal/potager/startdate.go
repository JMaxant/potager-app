package potager

// StartDate returns the earliest date of SownAt and TransplantedAt, or a zero Date{} if neither were provided.
func StartDate(p Planting) Date {
	if p.SownAt.IsZero() {
		return p.TransplantedAt
	}

	return p.SownAt
}
