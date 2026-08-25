package potager

// StartDate is the date a planting began occupying its plot, and the date
// maturity is counted from. Which of the two dates that is depends on the
// species: transplanted crops count from the transplant, direct-sown crops from
// the sowing. Either date falls back to the other when only one was recorded.
func StartDate(p Planting) (Date, bool) {
	from := MaturityFromSowing
	if len(p.Varieties) > 0 {
		from = p.Varieties[0].Variety.Species.MaturityFrom
	}

	first, second := p.SownAt, p.TransplantedAt
	if from == MaturityFromTransplant {
		first, second = p.TransplantedAt, p.SownAt
	}
	if first != nil {
		return *first, true
	}
	if second != nil {
		return *second, true
	}
	return Date{}, false
}

// ExpectedEnd reports when a planting is expected to leave its plot.
//
// A stored ExpectedEndAt always wins — the plan meets the weather, and the
// gardener is the authority on the outcome. Otherwise a food planting is
// computed from its maturity, and a cover crop is left open: its termination is
// driven backwards from when the bed is next needed, not forwards from sowing,
// so there is nothing to compute here.
func ExpectedEnd(p Planting) (Date, bool) {
	if p.ExpectedEndAt != nil {
		return *p.ExpectedEndAt, true
	}
	if p.Purpose == PurposeCover {
		return Date{}, false
	}

	days, ok := maturityDays(p)
	if !ok {
		return Date{}, false
	}
	start, ok := StartDate(p)
	if !ok {
		return Date{}, false
	}
	return start.AddDays(days), true
}

// maturityDays takes the slowest component of a planting: a mix is ready when
// its last member is.
func maturityDays(p Planting) (int, bool) {
	if p.MaturityDaysOverride != nil {
		return *p.MaturityDaysOverride, true
	}
	max, found := 0, false
	for _, pv := range p.Varieties {
		if d := pv.Variety.DaysToMaturity; d != nil && *d > max {
			max, found = *d, true
		}
	}
	return max, found
}

// TerminationLeadDays is how long before the bed is needed this planting must be
// terminated. The slowest component sets the pace: seigle in a mix is
// allelopathic and decomposes slowly, so it governs the whole sowing.
func TerminationLeadDays(p Planting) (int, bool) {
	max, found := 0, false
	for _, pv := range p.Varieties {
		if d := pv.Variety.Species.TerminationLeadDays; d != nil && *d > max {
			max, found = *d, true
		}
	}
	return max, found
}

// Season is the year a planting is counted under: the year it left the ground.
// A gardener remembers a season by what came out of it, which keeps garlic sown
// in October 2025 and lifted in July 2026 in one season rather than two.
func Season(p Planting) (int, bool) {
	if p.EndedAt != nil {
		return p.EndedAt.Year(), true
	}
	if end, ok := ExpectedEnd(p); ok {
		return end.Year(), true
	}
	return 0, false
}
