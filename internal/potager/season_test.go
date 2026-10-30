package potager

import (
	"testing"
	"time"
)

func TestSeason(t *testing.T) {
	tests := []struct {
		name    string
		endedAt Date
		want    int
	}{
		{
			name:    "garlic sown in October 2025 and lifted in July 2026 is season 2026 (AC 8)",
			endedAt: DateFromInt(2026, time.July, 4),
			want:    2026,
		},
		{
			name:    "crop harvested in 2026 is in season 2026",
			endedAt: DateFromInt(2026, time.January, 1),
			want:    2026,
		},
		{
			name:    "crop harvested in 2025 is in season 2025",
			endedAt: DateFromInt(2025, time.December, 31),
			want:    2025,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Season(tt.endedAt)
			if got != tt.want {
				t.Errorf("Season(%v) = %d, want %d", tt.endedAt, got, tt.want)
			}
		})
	}
}
