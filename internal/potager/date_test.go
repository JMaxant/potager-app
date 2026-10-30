package potager

import (
	"testing"
	"time"
)

func TestDateFromTime(t *testing.T) {
	paris, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		in   time.Time
		want Date
	}{
		{
			name: "late evening in winter stays on the same day in Paris",
			in:   time.Date(2026, time.March, 28, 22, 30, 0, 0, time.UTC), // 23:30 CET
			want: DateFromInt(2026, time.March, 28),
		},
		{
			name: "UTC timezone difference ends up in next day in Europe/Paris",
			in:   time.Date(2026, time.March, 28, 23, 30, 0, 0, time.UTC),
			want: DateFromInt(2026, time.March, 29),
		},
		{
			name: "UTC timezone difference + DST ends up in next day in Europe/Paris",
			in:   time.Date(2026, time.July, 15, 22, 30, 0, 0, time.UTC),
			want: DateFromInt(2026, time.July, 16),
		},
		{
			name: "DST switch in Europe/Paris",
			in:   time.Date(2026, time.March, 29, 1, 30, 0, 0, time.UTC),
			want: DateFromInt(2026, time.March, 29),
		},
		{
			name: "just after midnight in Europe/Paris is still that day in Europe/Paris",
			in:   time.Date(2025, time.October, 12, 0, 30, 0, 0, paris),
			want: DateFromInt(2025, time.October, 12),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DateFromTime(tt.in)
			if got != tt.want {
				t.Errorf("DateFromTime(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}
