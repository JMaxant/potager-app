package potager_test

import (
	"testing"

	"potager-app/internal/potager"
)

func TestNetEffect(t *testing.T) {
	tests := []struct {
		name    string
		purpose potager.Purpose
		demand  potager.NitrogenDemand
		fixing  bool
		want    potager.SoilEffect
	}{
		{
			name:    "vesce as a cover crop fixes nitrogen and exports none",
			purpose: potager.PurposeCover, demand: potager.DemandLight, fixing: true,
			want: potager.EffectEnriching,
		},
		{
			name:    "moutarde as a cover crop is hungry but returns everything on incorporation",
			purpose: potager.PurposeCover, demand: potager.DemandLight, fixing: false,
			want: potager.EffectNeutral,
		},
		{
			name:    "a heavy feeder sown as a cover crop is still neutral: it is never harvested",
			purpose: potager.PurposeCover, demand: potager.DemandHeavy, fixing: false,
			want: potager.EffectNeutral,
		},
		{
			name:    "peas grown for food carry their fixed nitrogen away in the seed",
			purpose: potager.PurposeFood, demand: potager.DemandLight, fixing: true,
			want: potager.EffectNeutral,
		},
		{
			name:    "tomatoes deplete the bed",
			purpose: potager.PurposeFood, demand: potager.DemandHeavy, fixing: false,
			want: potager.EffectDepleting,
		},
		{
			name:    "carrots take a little",
			purpose: potager.PurposeFood, demand: potager.DemandModerate, fixing: false,
			want: potager.EffectSlightlyDepleting,
		},
		{
			name:    "radishes take a little",
			purpose: potager.PurposeFood, demand: potager.DemandLight, fixing: false,
			want: potager.EffectSlightlyDepleting,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := potager.NetEffect(tt.purpose, tt.demand, tt.fixing); got != tt.want {
				t.Errorf("NetEffect(%q, %q, fixing=%v) = %q, want %q",
					tt.purpose, tt.demand, tt.fixing, got, tt.want)
			}
		})
	}
}

func TestNetEffectOfMix(t *testing.T) {
	fabaceae := potager.Family{Name: "Fabacées", NitrogenFixing: true}
	brassicaceae := potager.Family{Name: "Brassicacées"}

	trefle := variety("Trèfle incarnat", fabaceae, potager.DemandLight)
	radisFourrager := variety("Radis fourrager", brassicaceae, potager.DemandLight)

	t.Run("a pure legume cover crop enriches", func(t *testing.T) {
		p := potager.Planting{Purpose: potager.PurposeCover, Varieties: []potager.PlantingVariety{{Variety: trefle}}}
		if got := potager.NetEffectOf(p); got != potager.EffectEnriching {
			t.Errorf("got %q, want %q", got, potager.EffectEnriching)
		}
	})

	t.Run("a mix is judged by its least generous member", func(t *testing.T) {
		p := potager.Planting{Purpose: potager.PurposeCover, Varieties: []potager.PlantingVariety{
			{Variety: trefle}, {Variety: radisFourrager},
		}}
		if got := potager.NetEffectOf(p); got != potager.EffectNeutral {
			t.Errorf("got %q, want %q", got, potager.EffectNeutral)
		}
	})
}

func variety(name string, f potager.Family, d potager.NitrogenDemand) potager.Variety {
	return potager.Variety{
		Name:    name,
		Species: potager.Species{Name: name, Family: f, NitrogenDemand: d, MaturityFrom: potager.MaturityFromSowing},
	}
}
