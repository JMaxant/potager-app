package potager

import "testing"

func TestNetEffect(t *testing.T) {
	tests := []struct {
		name          string
		purpose       Purpose
		fixesNitrogen bool
		demand        Demand
		want          Effect
	}{
		{
			name:          "a nitrogen-fixing cover crop enriches the bed",
			purpose:       PurposeCover,
			fixesNitrogen: true,
			demand:        DemandModerate,
			want:          EffectEnriching,
		},
		{
			name:          "a non-nitrogen-fixing cover crop has a neutral effect",
			purpose:       PurposeCover,
			fixesNitrogen: false,
			demand:        DemandHeavy,
			want:          EffectNeutral,
		},
		{
			name:          "a nitrogen-fixing food crop has a neutral effect",
			purpose:       PurposeFood,
			fixesNitrogen: true,
			demand:        DemandHeavy,
			want:          EffectNeutral,
		},
		{
			name:          "a heavy demand non-nitrogen-fixing food crop depletes the bed",
			purpose:       PurposeFood,
			fixesNitrogen: false,
			demand:        DemandHeavy,
			want:          EffectDepleting,
		},
		{
			name:          "a moderate demand non-nitrogen-fixing food crop slightly depletes the bed",
			purpose:       PurposeFood,
			fixesNitrogen: false,
			demand:        DemandModerate,
			want:          EffectSlightlyDepleting,
		},
		{
			name:          "a light demand non-nitrogen-fixing food crop slightly depletes the bed",
			purpose:       PurposeFood,
			fixesNitrogen: false,
			demand:        DemandLight,
			want:          EffectSlightlyDepleting,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NetEffect(tt.purpose, tt.fixesNitrogen, tt.demand)
			if got != tt.want {
				t.Errorf("NetEffect(%v, %v, %v) = %v, want %v",
					tt.purpose, tt.fixesNitrogen, tt.demand, got, tt.want)
			}
		})
	}
}

func TestNetEffectPanics(t *testing.T) {
	tests := []struct {
		name          string
		purpose       Purpose
		fixesNitrogen bool
		demand        Demand
	}{
		{
			name:          "an unknown purpose panics",
			purpose:       Purpose("orchard"),
			fixesNitrogen: false,
			demand:        DemandHeavy,
		},
		{
			name:          "a food crop with an unknown demand panics",
			purpose:       PurposeFood,
			fixesNitrogen: false,
			demand:        "Huge",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				// If need be, move r := recover() out of the check, and and t.Log(r)
				if r := recover(); r == nil {
					t.Errorf("NetEffect(%q, %v, %q) did not panic",
						tt.purpose, tt.fixesNitrogen, tt.demand)
				}
			}()
			NetEffect(tt.purpose, tt.fixesNitrogen, tt.demand)
		})
	}
}
