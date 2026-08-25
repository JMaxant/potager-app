package potager_test

import (
	"testing"

	"potager-app/internal/potager"
)

func TestExpectedEnd(t *testing.T) {
	solanaceae := potager.Family{Name: "Solanacées", RotationIntervalYears: 4}
	apiaceae := potager.Family{Name: "Apiacées", RotationIntervalYears: 3}

	tomate := timedVariety("Marmande", solanaceae, potager.MaturityFromTransplant, 80)
	carotte := timedVariety("Nantaise", apiaceae, potager.MaturityFromSowing, 80)
	phacelie := potager.Variety{Name: "Phacélie", Species: potager.Species{
		Name: "Phacélie", MaturityFrom: potager.MaturityFromSowing,
	}}

	sown := potager.NewDate(2026, 3, 15)
	planted := potager.NewDate(2026, 5, 10)

	t.Run("a transplanted crop counts from the transplant, not the sowing", func(t *testing.T) {
		p := potager.Planting{
			Purpose:        potager.PurposeFood,
			SownAt:         &sown,
			TransplantedAt: &planted,
			Varieties:      []potager.PlantingVariety{{Variety: tomate}},
		}
		got, ok := potager.ExpectedEnd(p)
		if !ok {
			t.Fatal("expected a computed end date")
		}
		if want := potager.NewDate(2026, 7, 29); !got.Equal(want.Time) {
			t.Errorf("got %s, want %s", got.Format("2006-01-02"), want.Format("2006-01-02"))
		}
	})

	t.Run("a direct-sown crop counts from the sowing", func(t *testing.T) {
		p := potager.Planting{
			Purpose:   potager.PurposeFood,
			SownAt:    &sown,
			Varieties: []potager.PlantingVariety{{Variety: carotte}},
		}
		got, ok := potager.ExpectedEnd(p)
		if !ok {
			t.Fatal("expected a computed end date")
		}
		if want := potager.NewDate(2026, 6, 3); !got.Equal(want.Time) {
			t.Errorf("got %s, want %s", got.Format("2006-01-02"), want.Format("2006-01-02"))
		}
	})

	t.Run("a cover crop has no computed end", func(t *testing.T) {
		p := potager.Planting{
			Purpose:   potager.PurposeCover,
			SownAt:    &sown,
			Varieties: []potager.PlantingVariety{{Variety: phacelie}},
		}
		if _, ok := potager.ExpectedEnd(p); ok {
			t.Error("expected no computed end for a cover crop")
		}
	})

	t.Run("a stored end date always wins", func(t *testing.T) {
		override := potager.NewDate(2026, 6, 1)
		p := potager.Planting{
			Purpose:       potager.PurposeCover,
			SownAt:        &sown,
			ExpectedEndAt: &override,
			Varieties:     []potager.PlantingVariety{{Variety: phacelie}},
		}
		got, ok := potager.ExpectedEnd(p)
		if !ok || !got.Equal(override.Time) {
			t.Errorf("got %v (ok=%v), want the stored override %s", got, ok, override.Format("2006-01-02"))
		}
	})

	t.Run("a mix ends when its slowest member ends", func(t *testing.T) {
		fast := timedVariety("Radis", apiaceae, potager.MaturityFromSowing, 25)
		p := potager.Planting{
			Purpose:   potager.PurposeFood,
			SownAt:    &sown,
			Varieties: []potager.PlantingVariety{{Variety: fast}, {Variety: carotte}},
		}
		got, _ := potager.ExpectedEnd(p)
		if want := potager.NewDate(2026, 6, 3); !got.Equal(want.Time) {
			t.Errorf("got %s, want %s", got.Format("2006-01-02"), want.Format("2006-01-02"))
		}
	})
}

func TestSeasonIsTheYearItLeavesTheGround(t *testing.T) {
	sown := potager.NewDate(2025, 10, 12)
	lifted := potager.NewDate(2026, 7, 4)

	ail := potager.Planting{
		Purpose:   potager.PurposeFood,
		SownAt:    &sown,
		EndedAt:   &lifted,
		Varieties: []potager.PlantingVariety{{Variety: timedVariety("Blanc de Lomagne", potager.Family{Name: "Alliacées"}, potager.MaturityFromSowing, 240)}},
	}

	got, ok := potager.Season(ail)
	if !ok {
		t.Fatal("expected a season")
	}
	if got != 2026 {
		t.Errorf("garlic sown in Oct 2025 and lifted in Jul 2026 belongs to season %d, want 2026", got)
	}
}

func timedVariety(name string, f potager.Family, from potager.MaturityFrom, days int) potager.Variety {
	return potager.Variety{
		Name:           name,
		DaysToMaturity: &days,
		Species:        potager.Species{Name: name, Family: f, MaturityFrom: from, NitrogenDemand: potager.DemandModerate},
	}
}
