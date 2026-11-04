package potager

import (
	"testing"
	"time"
)

func TestExpectedEnd(t *testing.T) {
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
			gotDate, gotSource := ExpectedEnd(tt.planting)
			if gotDate != tt.wantDate || gotSource != tt.wantSource {
				t.Errorf("ExpectedEnd(%+v) = (%v, %q), want (%v, %q)",
					tt.planting, gotDate, gotSource, tt.wantDate, tt.wantSource)
			}
		})
	}
}

func TestExpectedEndPanics(t *testing.T) {
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
					t.Errorf("ExpectedEnd(%+v) did not panic (improve message)", tt.planting)
				}
			}()
			ExpectedEnd(tt.planting)
		})
	}
}
