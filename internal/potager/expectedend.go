package potager

import "fmt"

// MaturityFrom represents the date from which days to maturity are counted.
// There is no fallback: a transplant species without a transplant date has no expected end.
type MaturityFrom string

const (
	MaturityFromSowing     MaturityFrom = "sowing"
	MaturityFromTransplant MaturityFrom = "transplant"
)

// ExpectedEndSource tells what the source of expected end date is.
type ExpectedEndSource string

const (
	ExpectedEndSourceComputed ExpectedEndSource = "computed"
	ExpectedEndSourceStored   ExpectedEndSource = "stored"
	// ExpectedEndSourceNone means there is no date such as a cover crop with no date entered, a food crop with no days to maturity or a store-bought MaturityFromSowing crop.
	ExpectedEndSourceNone ExpectedEndSource = "none"
	// ExpectedEndSourceNeedsTransplant means there is no date because the transplant date is missing; unlike None, the user can supply it.
	ExpectedEndSourceNeedsTransplant ExpectedEndSource = "needs transplant"
)

// ExpectedEnd returns the expected end date of a planting, and the source of that date. It panics if Purpose is unknown of MaturityFrom is unknown for a food crop.
func ExpectedEnd(p Planting) (Date, ExpectedEndSource) {
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

	switch p.MaturityFrom {
	case MaturityFromSowing:
	case MaturityFromTransplant:
	default:
		panic(fmt.Sprintf("ExpectedEnd: unknown MaturityFrom %q", p.MaturityFrom))
	}

	switch {
	case p.MaturityFrom == MaturityFromTransplant && p.TransplantedAt.IsZero():
		return Date{}, ExpectedEndSourceNeedsTransplant

	case p.MaturityFrom == MaturityFromTransplant && p.MaturityDaysOverride != 0:
		return p.TransplantedAt.addDays(p.MaturityDaysOverride), ExpectedEndSourceComputed

	case p.MaturityFrom == MaturityFromTransplant && p.DaysToMaturity != 0:
		return p.TransplantedAt.addDays(p.DaysToMaturity), ExpectedEndSourceComputed

	case p.MaturityFrom == MaturityFromSowing && p.SownAt.IsZero():
		return Date{}, ExpectedEndSourceNone

	case p.MaturityFrom == MaturityFromSowing && p.MaturityDaysOverride != 0:
		return p.SownAt.addDays(p.MaturityDaysOverride), ExpectedEndSourceComputed

	case p.MaturityFrom == MaturityFromSowing && p.DaysToMaturity != 0:
		return p.SownAt.addDays(p.DaysToMaturity), ExpectedEndSourceComputed

	default:
		return Date{}, ExpectedEndSourceNone
	}
}
