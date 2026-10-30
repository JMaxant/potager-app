package potager

// MaturityFrom represents the step from which days to maturity are counted.
type MaturityFrom string

const (
	MaturityFromSowing     MaturityFrom = "sowing"
	MaturityFromTransplant MaturityFrom = "transplant"
)
