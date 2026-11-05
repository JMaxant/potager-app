package potager

// Plot is a bounded growing area of the garden. A 0 width or length is an unspecified value.
type Plot struct {
	LengthCm int
	WidthCm  int
}

// AreaCm2 returns the plot's area in cm², or 0 when either dimension is
// missing. Converting to m² for display is left to web.
func (p Plot) AreaCm2() int {
	return p.LengthCm * p.WidthCm
}
