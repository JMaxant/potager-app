package potager

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
	// ExpectedEndSourceNone means there is no date such as a cover crop with no date entered, or a food crop with no days to maturity.
	ExpectedEndSourceNone ExpectedEndSource = "none"
	// ExpectedEndSourceNeedsTransplant means there is no date because the transplant date is missing; unlike None, the user can supply it.
	ExpectedEndSourceNeedsTransplant ExpectedEndSource = "needs transplant"
)

// ExpectedEnd returns the expected end date of a planting, and the source of that date.
func ExpectedEnd(p Planting) (Date, ExpectedEndSource) {
	var endDate Date
	var source ExpectedEndSource
	switch {
	case p.Purpose == PurposeCover:
		return Date{}, ExpectedEndSourceNone
	case p.Purpose == PurposeFood && p.DaysToMaturity != 0:
		return p.TransplantedAt.addDays(p.DaysToMaturity), ExpectedEndSourceComputed

	case p.Purpose == PurposeFood && p.DaysToMaturity == 0:
	}
	if p.DaysToMaturity != 0 {
		endDate = p.TransplantedAt.addDays(p.DaysToMaturity)
		source = ExpectedEndSourceComputed
	}
	return endDate, source
}
