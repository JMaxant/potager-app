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
				TransplantedAt: Date{},
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
				TransplantedAt: DateFromInt(2026, time.April, 01),
				DaysToMaturity: 30,
			},
			wantDate:   DateFromInt(2026, time.May, 01),
			wantSource: ExpectedEndSourceComputed,
		},
		{},
		// TODO(user): cases 2 to 10 of the plan in docs/step-1-t3-notes.md
		// (section "TestExpectedEnd plan"). One case per sentence of the rule;
		// list only the fields that matter, the others stay zero, i.e. absent.
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
