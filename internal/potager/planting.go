package potager

import "fmt"

// Planting represents one act of planting or sowing, with any purpose (food or cover).
// A zero date or int means absent (e.g: a bought tomato plant).
// DaysToMaturity and MaturityFrom come from the catalogue.
// MaturityDaysOverride is an user input.
// A false SownInPlace is a real value, rather that a Zero Value (i.e : it was either sown in place or it wasn't).
type Planting struct {
	Purpose              Purpose
	MaturityFrom         MaturityFrom
	SownAt               Date
	SownInPlace          bool
	TransplantedAt       Date
	ExpectedEndAt        Date
	DaysToMaturity       int
	MaturityDaysOverride int
}

// OccupiedFrom returns the first day the planting holds its plot, or a zero Date while it waits in its tray (ADR 0008).
func (p Planting) OccupiedFrom() Date {
	if p.SownInPlace {
		return p.SownAt
	}

	return p.TransplantedAt
}

// StartDate returns the earliest date of SownAt and TransplantedAt, or a zero Date{} if neither were provided.
// A crop is always sown before it is transplanted, hence SownAt takes priority if it is provided.
func (p Planting) StartDate() Date {
	if p.SownAt.IsZero() {
		return p.TransplantedAt
	}

	return p.SownAt
}

// ExpectedEnd returns the expected end date of a planting, and the source of that date. It panics if Purpose is unknown or MaturityFrom is unknown for a food crop.
func (p Planting) ExpectedEnd() (Date, ExpectedEndSource) {
	// If a Planting provides an ExpectedEndAt, ExpectedEndAt wins.
	if !p.ExpectedEndAt.IsZero() {
		return p.ExpectedEndAt, ExpectedEndSourceStored
	}

	switch p.Purpose {
	case PurposeCover:
		return Date{}, ExpectedEndSourceNone
	case PurposeFood:

	default:
		panic(fmt.Sprintf("ExpectedEnd: unknown Purpose %q", p.Purpose))
	}

	var start Date
	switch p.MaturityFrom {
	case MaturityFromSowing:
		start = p.SownAt
	case MaturityFromTransplant:
		start = p.TransplantedAt
		if start.IsZero() {
			return Date{}, ExpectedEndSourceNeedsTransplant
		}
	default:
		panic(fmt.Sprintf("ExpectedEnd: unknown MaturityFrom %q", p.MaturityFrom))
	}

	days := p.DaysToMaturity
	if p.MaturityDaysOverride != 0 {
		days = p.MaturityDaysOverride
	}

	if start.IsZero() || days == 0 {
		return Date{}, ExpectedEndSourceNone
	}

	return start.addDays(days), ExpectedEndSourceComputed
}
