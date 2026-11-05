package potager

import "fmt"

// Purpose is why the crop is cultivated.
type Purpose string

const (
	// PurposeFood is meant to be harvested.
	PurposeFood Purpose = "food"
	// PurposeCover is meant to be terminated.
	PurposeCover Purpose = "cover"
)

// Demand is a crop's nitrogen needs.
type Demand string

const (
	DemandLight    Demand = "light"
	DemandModerate Demand = "moderate"
	DemandHeavy    Demand = "heavy"
)

// Effect is the impact of a crop on the soil. Computed in NetEffect.
type Effect string

const (
	EffectDepleting         Effect = "depleting"
	EffectSlightlyDepleting Effect = "slightly depleting"
	EffectNeutral           Effect = "neutral"
	EffectEnriching         Effect = "enriching"
)

// NetEffect calculates the net effect of a crop on the soil. It panics if provided invalid demand or purpose.
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

	case purpose == PurposeFood && (demand == DemandLight || demand == DemandModerate):
		return EffectSlightlyDepleting

	case purpose == PurposeFood:
		panic(fmt.Sprintf("NetEffect: unknown demand %q", demand))

	default:
		panic(fmt.Sprintf("NetEffect: unknown purpose %q", purpose))
	}
}
