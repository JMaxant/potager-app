package potager

type Purpose string

const (
	PurposeFood  Purpose = "food"
	PurposeCover Purpose = "cover"
)

type Demand string

const (
	DemandLight    Demand = "light"
	DemandModerate Demand = "moderate"
	DemandHeavy    Demand = "heavy"
)

type Effect string

const (
	EffectDepleting         Effect = "depleting"
	EffectSlightlyDepleting Effect = "slightly depleting"
	EffectNeutral           Effect = "neutral"
	EffectEnriching         Effect = "enriching"
)

// NetEffect calculates the net effect of a crop on the soil.
func NetEffect(purpose Purpose, fixesNitrogen bool, demand Demand) Effect {
	switch {
	case purpose == PurposeCover && fixesNitrogen:
		return EffectEnriching

	case purpose == PurposeCover:
		return EffectNeutral

	case purpose == PurposeFood && fixesNitrogen:
		return EffectNeutral

	case purpose == PurposeFood && demand == DemandHeavy:
		return EffectDepleting

	case purpose == PurposeFood:
		return EffectSlightlyDepleting

	default:
		panic("todo msg")
	}
}
