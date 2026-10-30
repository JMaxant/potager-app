package potager

// Planting represents one act of planting or sowing, with any purpose (food or cover).
// A zero date or int mean absent (e.g: a bought tomato plant).
// Maturity days fields come from the catalogue.
type Planting struct {
	Purpose              Purpose
	MaturityFrom         MaturityFrom
	SownAt               Date
	TransplantedAt       Date
	ExpectedEndAt        Date
	DaysToMaturity       int
	MaturityDaysOverride int
}
