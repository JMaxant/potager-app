package potager

import (
	"testing"
	"time"
)

func TestPlanting_OccupiedFrom(t *testing.T) {
	tests := []struct {
		name     string
		planting Planting
		want     Date
	}{
		{
			name: "A tomato sown in a tray on 15/03/2026 occupies its plot from its planting out on 10/05/2026",
			planting: Planting{
				MaturityFrom:   MaturityFromTransplant,
				SownAt:         DateFromInt(2026, time.March, 15),
				SownInPlace:    false,
				TransplantedAt: DateFromInt(2026, time.May, 10),
			},
			want: DateFromInt(2026, time.May, 10),
		},
		{
			name: "a tomato sown in a tray on 15/03/2026 but not yet planted does not occupy a plot",
			planting: Planting{
				MaturityFrom: MaturityFromTransplant,
				SownAt:       DateFromInt(2026, time.March, 15),
				SownInPlace:  false,
			},
			want: Date{},
		},
		{
			name: "fava beans sown in place on 15/11/2025 occupies a plot from its sowing date",
			planting: Planting{
				MaturityFrom: MaturityFromSowing,
				SownAt:       DateFromInt(2025, time.November, 15),
				SownInPlace:  true,
			},
			want: DateFromInt(2025, time.November, 15),
		},
		{
			name: "bought lettuce seedlings, transplanted on 01/05/2026 occupies a plot from its TransplantedAt date",
			planting: Planting{
				MaturityFrom:   MaturityFromSowing,
				TransplantedAt: DateFromInt(2026, time.May, 1),
				SownInPlace:    false,
			},
			want: DateFromInt(2026, time.May, 1),
		},
		{
			name: "lettuce sown, not yet transplanted, does not occupy a plot, although its species is sown",
			planting: Planting{
				MaturityFrom: MaturityFromSowing,
				SownAt:       DateFromInt(2026, time.May, 1),
				SownInPlace:  false,
			},
			want: Date{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.planting.OccupiedFrom()
			if got != tt.want {
				t.Errorf("Planting.OccupiedFrom() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanting_StartDate(t *testing.T) {
	tests := []struct {
		name     string
		planting Planting
		want     Date
	}{
		{
			name: "a tomato sown in a tray on 15/03/2026 and planted out on 10/05/2026 starts on 15/03/2026",
			planting: Planting{
				SownAt:         DateFromInt(2026, time.March, 15),
				TransplantedAt: DateFromInt(2026, time.May, 10),
			},
			want: DateFromInt(2026, time.March, 15),
		},
		{
			name: "A lettuce bought, planted only on 01/05/2026, starts on 01/05/2026",
			planting: Planting{
				TransplantedAt: DateFromInt(2026, time.May, 1),
			},
			want: DateFromInt(2026, time.May, 1),
		},
		{
			name: "Fava beans sown on 15/11/2026, no transplant, starts on 15/11/2026",
			planting: Planting{
				SownAt: DateFromInt(2026, time.November, 15),
			},
			want: DateFromInt(2026, time.November, 15),
		},
		{
			name: "A crop sown and transplanted on the same day returns that day",
			planting: Planting{
				SownAt:         DateFromInt(2026, time.March, 15),
				TransplantedAt: DateFromInt(2026, time.March, 15),
			},
			want: DateFromInt(2026, time.March, 15),
		},
		{
			name:     "A crop with or transplantation date provided return a zero Date",
			planting: Planting{},
			want:     Date{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.planting.StartDate()
			if got != tt.want {
				t.Errorf("Planting.StartDate() = Planting %+v,  %v, want %v", tt.planting, got, tt.want)
			}
		})
	}
}

func TestPlanting_ExpectedEnd(t *testing.T) {
	tests := []struct {
		name       string
		planting   Planting
		wantDate   Date
		wantSource ExpectedEndSource
	}{
		{
			name: "an 80-day tomato transplanted on 10/05/2026 is expected on 29/07/2026 (AC 5)",
			planting: Planting{
				Purpose:        PurposeFood,
				MaturityFrom:   MaturityFromTransplant,
				TransplantedAt: DateFromInt(2026, time.May, 10),
				DaysToMaturity: 80,
			},
			wantDate:   DateFromInt(2026, time.July, 29),
			wantSource: ExpectedEndSourceComputed,
		},
		{
			name: "an 80-day tomato on 10/05/2026 not transplanted expects ExpectedEndSourceNeedsTransplant",
			planting: Planting{
				Purpose:        PurposeFood,
				MaturityFrom:   MaturityFromTransplant,
				SownAt:         DateFromInt(2026, time.May, 10),
				DaysToMaturity: 80,
			},
			wantDate:   Date{},
			wantSource: ExpectedEndSourceNeedsTransplant,
		},
		{
			name: "a 30-day radish sown on 01/04/2026 is expected on 01/05/2026",
			planting: Planting{
				Purpose:        PurposeFood,
				MaturityFrom:   MaturityFromSowing,
				SownAt:         DateFromInt(2026, time.April, 1),
				DaysToMaturity: 30,
			},
			wantDate:   DateFromInt(2026, time.May, 1),
			wantSource: ExpectedEndSourceComputed,
		},
		{
			name: "an 80-day tomato transplanted on 10/05/2026 with a 90 days override is expected on 08/08/2026",
			planting: Planting{
				Purpose:              PurposeFood,
				MaturityFrom:         MaturityFromTransplant,
				TransplantedAt:       DateFromInt(2026, time.May, 10),
				DaysToMaturity:       80,
				MaturityDaysOverride: 90,
			},
			wantDate:   DateFromInt(2026, time.August, 8),
			wantSource: ExpectedEndSourceComputed,
		},
		{
			name: "a food crop with a stored ExpectedEndAt is expected as ExpectedEndSourceStored rather than ExpectedEndSourceComputed",
			planting: Planting{
				Purpose:        PurposeFood,
				MaturityFrom:   MaturityFromSowing,
				SownAt:         DateFromInt(2026, time.April, 1),
				DaysToMaturity: 80,
				ExpectedEndAt:  DateFromInt(2026, time.August, 8),
			},
			wantDate:   DateFromInt(2026, time.August, 8),
			wantSource: ExpectedEndSourceStored,
		},
		{
			name: "a cover crop sown with days to maturity but no stored end has no expected end (AC 6)",
			planting: Planting{
				Purpose:        PurposeCover,
				MaturityFrom:   MaturityFromSowing,
				SownAt:         DateFromInt(2026, time.April, 1),
				DaysToMaturity: 120,
			},
			wantDate:   Date{},
			wantSource: ExpectedEndSourceNone,
		},
		{
			name: "a Cover Crop with a stored ExpectedEndAt uses it (AC 6)",
			planting: Planting{
				Purpose:       PurposeCover,
				MaturityFrom:  MaturityFromSowing,
				ExpectedEndAt: DateFromInt(2026, time.August, 8),
			},
			wantDate:   DateFromInt(2026, time.August, 8),
			wantSource: ExpectedEndSourceStored,
		},
		{
			name: "a food crop without DaysToMaturity or override is expected as ExpectedEndSourceNone",
			planting: Planting{
				Purpose:      PurposeFood,
				MaturityFrom: MaturityFromSowing,
				SownAt:       DateFromInt(2026, time.April, 1),
			},
			wantDate:   Date{},
			wantSource: ExpectedEndSourceNone,
		},
		{
			name: "a sown species with no sowing date (bought) shows no expected date",
			planting: Planting{
				Purpose:        PurposeFood,
				MaturityFrom:   MaturityFromSowing,
				DaysToMaturity: 30,
				TransplantedAt: DateFromInt(2026, time.May, 1),
			},
			wantDate:   Date{},
			wantSource: ExpectedEndSourceNone,
		},
		{
			name: "a sown species with a transplant date (started in a pot) still counts from the sowing",
			planting: Planting{
				Purpose:        PurposeFood,
				MaturityFrom:   MaturityFromSowing,
				SownAt:         DateFromInt(2026, time.April, 1),
				TransplantedAt: DateFromInt(2026, time.April, 20),
				DaysToMaturity: 30,
			},
			wantDate:   DateFromInt(2026, time.May, 1),
			wantSource: ExpectedEndSourceComputed,
		},
		{
			name: "a sown species with an override days to maturity uses it",
			planting: Planting{
				Purpose:              PurposeFood,
				MaturityFrom:         MaturityFromSowing,
				SownAt:               DateFromInt(2026, time.April, 1),
				DaysToMaturity:       30,
				MaturityDaysOverride: 40,
			},
			wantDate:   DateFromInt(2026, time.May, 11),
			wantSource: ExpectedEndSourceComputed,
		},
		{
			name: "an unknown purpose with a stored expected end uses it, without panicking",
			planting: Planting{
				Purpose:        "fodder",
				SownAt:         DateFromInt(2026, time.April, 1),
				DaysToMaturity: 30,
				MaturityFrom:   MaturityFromSowing,
				ExpectedEndAt:  DateFromInt(2026, time.May, 15),
			},
			wantDate:   DateFromInt(2026, time.May, 15),
			wantSource: ExpectedEndSourceStored,
		},
		{
			name: "a cover crop without MaturityFrom or stored expected end has none, without panicking",
			planting: Planting{
				Purpose:        PurposeCover,
				SownAt:         DateFromInt(2026, time.April, 1),
				DaysToMaturity: 30,
			},
			wantDate:   Date{},
			wantSource: ExpectedEndSourceNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDate, gotSource := tt.planting.ExpectedEnd()
			if gotDate != tt.wantDate || gotSource != tt.wantSource {
				t.Errorf("Planting.ExpectedEnd() = Planting %+v : (%v, %q), want (%v, %q)",
					tt.planting, gotDate, gotSource, tt.wantDate, tt.wantSource)
			}
		})
	}
}

func TestPlanting_ExpectedEndPanics(t *testing.T) {
	tests := []struct {
		name     string
		planting Planting
	}{
		{
			name: "a food crop with an empty MaturityFrom panics",
			planting: Planting{
				Purpose:        PurposeFood,
				SownAt:         DateFromInt(2026, time.April, 1),
				DaysToMaturity: 30,
			},
		},
		{
			name: "a food crop with an unknown MaturityFrom panics",
			planting: Planting{
				Purpose:        PurposeFood,
				SownAt:         DateFromInt(2026, time.April, 1),
				DaysToMaturity: 30,
				MaturityFrom:   "transplanting",
			},
		},
		{
			name: "an unknown purpose without a stored expected end panics",
			planting: Planting{
				Purpose:        "fodder",
				SownAt:         DateFromInt(2026, time.April, 1),
				DaysToMaturity: 30,
				MaturityFrom:   MaturityFromSowing,
			},
		},
		{
			name: "a crop with an empty purpose panics",
			planting: Planting{
				SownAt:         DateFromInt(2026, time.April, 1),
				DaysToMaturity: 30,
				MaturityFrom:   MaturityFromSowing,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Planting.ExpectedEnd(): Planting %+v did not panic", tt.planting)
				}
			}()
			tt.planting.ExpectedEnd()
		})
	}
}

func TestPlanting_Status(t *testing.T) {
	tests := []struct {
		name     string
		planting Planting
		today    Date
		want     Status
	}{
		{
			name: "a tomato planted out on 10/05/2026, expected on 29/07/2026, is growing on 01/07/2026",
			planting: Planting{
				Purpose:        PurposeFood,
				MaturityFrom:   MaturityFromTransplant,
				SownInPlace:    false,
				TransplantedAt: DateFromInt(2026, time.May, 10),
				DaysToMaturity: 80,
			},
			today: DateFromInt(2026, time.July, 1),
			want:  StatusGrowing,
		},
		{
			name: "a crop harvested on 05/08/2026, is ended on 10/08/2026",
			planting: Planting{
				Purpose:        PurposeFood,
				MaturityFrom:   MaturityFromTransplant,
				SownInPlace:    false,
				TransplantedAt: DateFromInt(2026, time.May, 10),
				DaysToMaturity: 80,
				EndedAt:        DateFromInt(2026, time.August, 5),
			},
			today: DateFromInt(2026, time.August, 10),
			want:  StatusEnded,
		},
		{
			name: "a tomato sown on 15/03/2026, still in its tray, is Planned on 01/04/2026",
			planting: Planting{
				Purpose:        PurposeFood,
				MaturityFrom:   MaturityFromTransplant,
				SownAt:         DateFromInt(2026, time.March, 15),
				DaysToMaturity: 80,
			},
			today: DateFromInt(2026, time.April, 1),
			want:  StatusPlanned,
		},
		{
			name: "a crop whose transplant is expected on 20/05/2026, is Planned on 01/05/2026",
			planting: Planting{
				Purpose:        PurposeFood,
				MaturityFrom:   MaturityFromTransplant,
				TransplantedAt: DateFromInt(2026, time.May, 20),
				DaysToMaturity: 80,
			},
			today: DateFromInt(2026, time.May, 1),
			want:  StatusPlanned,
		},
		{
			name: "a crop past its expected end on 15/07/2026, not ended, is due on 01/08/2026",
			planting: Planting{
				Purpose:        PurposeFood,
				MaturityFrom:   MaturityFromTransplant,
				TransplantedAt: DateFromInt(2026, time.May, 1),
				ExpectedEndAt:  DateFromInt(2026, time.July, 15),
			},
			today: DateFromInt(2026, time.August, 1),
			want:  StatusDue,
		},
		{
			name: "a crop expected on 29/07/2026, is Growing on the 29/07/2026",
			planting: Planting{
				Purpose:        PurposeFood,
				MaturityFrom:   MaturityFromTransplant,
				TransplantedAt: DateFromInt(2026, time.May, 1),
				ExpectedEndAt:  DateFromInt(2026, time.July, 29),
			},
			today: DateFromInt(2026, time.July, 29),
			want:  StatusGrowing,
		},
		{
			name: "a cover crop without an expected end, sown in place on 15/11/2025, is still growing on 01/08/2026",
			planting: Planting{
				Purpose:      PurposeCover,
				MaturityFrom: MaturityFromSowing,
				SownAt:       DateFromInt(2025, time.November, 15),
				SownInPlace:  true,
			},
			today: DateFromInt(2026, time.August, 1),
			want:  StatusGrowing,
		},
		{
			name: "a crop whose ExpectedEnd is passed and EndedAt is provided, is Ended on its EndedAt",
			planting: Planting{
				Purpose:        PurposeFood,
				MaturityFrom:   MaturityFromTransplant,
				TransplantedAt: DateFromInt(2026, time.May, 1),
				ExpectedEndAt:  DateFromInt(2026, time.August, 1),
				EndedAt:        DateFromInt(2026, time.July, 20),
			},
			today: DateFromInt(2026, time.August, 10),
			want:  StatusEnded,
		},
		{
			name: "a crop planted on the same day is Growing",
			planting: Planting{
				Purpose:        PurposeFood,
				MaturityFrom:   MaturityFromTransplant,
				TransplantedAt: DateFromInt(2026, time.April, 1),
				DaysToMaturity: 80,
			},
			today: DateFromInt(2026, time.April, 1),
			want:  StatusGrowing,
		},
		{
			name: "a cover crop with an ExpectedEndAt provided is due",
			planting: Planting{
				Purpose:       PurposeCover,
				MaturityFrom:  MaturityFromSowing,
				SownAt:        DateFromInt(2026, time.April, 1),
				ExpectedEndAt: DateFromInt(2026, time.July, 20),
				SownInPlace:   true,
			},
			today: DateFromInt(2026, time.July, 21),
			want:  StatusDue,
		},
		{
			name: "a food crop without ExpectedEndAt or DaysToMaturity is Growing months later",
			planting: Planting{
				Purpose:      PurposeFood,
				MaturityFrom: MaturityFromSowing,
				SownAt:       DateFromInt(2026, time.April, 1),
				SownInPlace:  true,
			},
			today: DateFromInt(2026, time.August, 1),
			want:  StatusGrowing,
		},
		{
			name: "a bought lettuce planted out on 01/04/2026 is growing, with no expected end",
			planting: Planting{
				Purpose:        PurposeFood,
				MaturityFrom:   MaturityFromSowing,
				TransplantedAt: DateFromInt(2026, time.April, 1),
				DaysToMaturity: 60,
			},
			today: DateFromInt(2026, time.May, 1),
			want:  StatusGrowing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.planting.Status(tt.today)
			if got != tt.want {
				t.Errorf("Planting.Status() = Planting %+v, today %v : %v, want %v", tt.planting, tt.today, got, tt.want)
			}
		})
	}
}

func TestPlanting_InPlotBy(t *testing.T) {
	tests := []struct {
		name     string
		planting Planting
		today    Date
		want     bool
	}{
		{
			name: "a crop is still in its tray",
			planting: Planting{
				Purpose:     PurposeFood,
				SownAt:      DateFromInt(2026, time.March, 1),
				SownInPlace: false,
			},
			today: DateFromInt(2026, time.April, 1),
			want:  false,
		},
		{
			name: "a crop was transplanted before the given day",
			planting: Planting{
				Purpose:        PurposeFood,
				MaturityFrom:   MaturityFromTransplant,
				TransplantedAt: DateFromInt(2026, time.April, 1),
			},
			today: DateFromInt(2026, time.May, 1),
			want:  true,
		},
		{
			name: "a crop was transplanted on the given day",
			planting: Planting{
				Purpose:        PurposeFood,
				MaturityFrom:   MaturityFromTransplant,
				TransplantedAt: DateFromInt(2026, time.April, 1),
			},
			today: DateFromInt(2026, time.April, 1),
			want:  true,
		},
		{
			name: "a crop to be transplanted after the given day is not in its plot yet",
			planting: Planting{
				Purpose:        PurposeFood,
				MaturityFrom:   MaturityFromTransplant,
				TransplantedAt: DateFromInt(2026, time.June, 1),
			},
			today: DateFromInt(2026, time.May, 1),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.planting.InPlotBy(tt.today)
			if got != tt.want {
				t.Errorf("Planting.InPlotBy() = Planting %+v, today %v %v, want %v", tt.planting, tt.today, got, tt.want)
			}
		})
	}
}
