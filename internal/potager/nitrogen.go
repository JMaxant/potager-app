package potager

// SoilEffect is what a planting leaves behind for whatever follows it. It is
// deliberately never stored: it is a function of intent, appetite and botany,
// and any of the three can be corrected in the catalogue after the fact.
type SoilEffect string

const (
	EffectEnriching         SoilEffect = "enriching"
	EffectNeutral           SoilEffect = "neutral"
	EffectSlightlyDepleting SoilEffect = "slightly_depleting"
	EffectDepleting         SoilEffect = "depleting"
)

// NetEffect reports what a planting leaves in the soil.
//
// The distinction that matters is between what a plant *takes up* and what it
// *exports*. A cover crop is terminated before it fruits, so everything it took
// up is returned to the bed when it is terminated: mustard is a hungry plant
// with a neutral net effect. A legume grown for food is the mirror image — it
// fixes nitrogen and then carries most of it away in the seed, so a harvested
// pea crop enriches the bed far less than its nodules suggest.
func NetEffect(purpose Purpose, demand NitrogenDemand, nitrogenFixing bool) SoilEffect {
	switch {
	case purpose == PurposeCover && nitrogenFixing:
		// Fixed nitrogen, none of it exported: the reason to sow vesce or
		// trèfle ahead of a hungry crop.
		return EffectEnriching
	case purpose == PurposeCover:
		// A catch crop. It holds nutrients over winter that would otherwise
		// leach, and gives them back when it is terminated.
		return EffectNeutral
	case nitrogenFixing:
		// Grown for food: the nitrogen leaves in the harvest.
		return EffectNeutral
	case demand == DemandHeavy:
		return EffectDepleting
	default:
		return EffectSlightlyDepleting
	}
}

// NetEffectOf reports the net effect of a planting, taking the most demanding
// component of a mix. A mix is judged by its hungriest member.
func NetEffectOf(p Planting) SoilEffect {
	worst := EffectEnriching
	for _, pv := range p.Varieties {
		s := pv.Variety.Species
		if e := NetEffect(p.Purpose, s.NitrogenDemand, s.Family.NitrogenFixing); severity(e) > severity(worst) {
			worst = e
		}
	}
	return worst
}

func severity(e SoilEffect) int {
	switch e {
	case EffectEnriching:
		return 0
	case EffectNeutral:
		return 1
	case EffectSlightlyDepleting:
		return 2
	case EffectDepleting:
		return 3
	}
	return 0
}
