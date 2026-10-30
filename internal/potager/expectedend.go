package potager

// MaturityFrom represents the step from which days to maturity are counted.
type MaturityFrom string

const (
	MaturityFromSowing     MaturityFrom = "sowing"
	MaturityFromTransplant MaturityFrom = "transplant"
)

// ExpectedEndSource tells where the source of expected end date.
type ExpectedEndSource string

const (
	ExpectedEndSourceComputed ExpectedEndSource = "computed"
	ExpectedEndSourceStored   ExpectedEndSource = "stored"
	ExpectedEndSourceNone     ExpectedEndSource = "none"
	// ExpectedEndSourceNeedsTransplant The screen asks for the transplant date.
	ExpectedEndSourceNeedsTransplant ExpectedEndSource = "needs transplant"
)
