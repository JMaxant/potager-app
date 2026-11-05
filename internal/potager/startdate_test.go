package potager

import (
	"testing"
	"time"
)

func TestStartDate(t *testing.T) {
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
			got := StartDate(tt.planting)
			if got != tt.want {
				t.Errorf("StartDate(%+v) = %v, want %v", tt.planting, got, tt.want)
			}
		})
	}
}
