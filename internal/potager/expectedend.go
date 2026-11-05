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
	// ExpectedEndSourceNone means there is no date such as a cover crop with no date entered, a food crop with no days to maturity or a store-bought MaturityFromSowing crop.
	ExpectedEndSourceNone ExpectedEndSource = "none"
	// ExpectedEndSourceNeedsTransplant means there is no date because the transplant date is missing; unlike None, the user can supply it.
	ExpectedEndSourceNeedsTransplant ExpectedEndSource = "needs transplant"
)
