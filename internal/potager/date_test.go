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

func TestAddDays(t *testing.T) {
	tests := []struct {
		name string
		in   Date
		days int
		want Date
	}{
		{
			name: "80-day tomato transplanted on May 10th is ready on July 29th (AC 5)",
			in:   DateFromInt(2023, time.May, 10),
			days: 80,
			want: DateFromInt(2023, time.July, 29),
		},
		{
			name: "add 0 day.",
			in:   DateFromInt(2023, time.May, 10),
			days: 0,
			want: DateFromInt(2023, time.May, 10),
		},
		{
			name: "day addition changes year.",
			in:   DateFromInt(2023, time.December, 31),
			days: 1,
			want: DateFromInt(2024, time.January, 1),
		},
		{
			name: "leap year",
			in:   DateFromInt(2028, time.February, 28),
			days: 1,
			want: DateFromInt(2028, time.February, 29),
		},
		{
			name: "non-leap year",
			in:   DateFromInt(2028, time.February, 29),
			days: 1,
			want: DateFromInt(2028, time.March, 1),
		},
		{
			name: "substract day",
			in:   DateFromInt(2025, time.May, 10),
			days: -21,
			want: DateFromInt(2025, time.April, 19),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.addDays(tt.days)
			if got != tt.want {
				t.Errorf("AddDays(%v, %d) = %v, want %v", tt.in, tt.days, got, tt.want)
			}
		})
	}
}

func TestBefore(t *testing.T) {
	tests := []struct {
		name  string
		d     Date
		other Date
		want  bool
	}{
		{
			name:  "the day before is before",
			d:     DateFromInt(2026, time.July, 28),
			other: DateFromInt(2026, time.July, 29),
			want:  true,
		},
		{
			name:  "the day after is not before",
			d:     DateFromInt(2026, time.July, 30),
			other: DateFromInt(2026, time.July, 29),
			want:  false,
		},
		{
			name:  "the same day is not before",
			d:     DateFromInt(2026, time.July, 29),
			other: DateFromInt(2026, time.July, 29),
			want:  false,
		},
		{
			name:  "the last of the year is before the first day of the year",
			d:     DateFromInt(2025, time.December, 31),
			other: DateFromInt(2026, time.January, 1),
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.before(tt.other)
			if got != tt.want {
				t.Errorf("%v.Before(%v) = %v, want %v", tt.d, tt.other, got, tt.want)
			}
		})
	}
}

func TestDateIsZero(t *testing.T) {
	tests := []struct {
		name string
		date Date
		want bool
	}{
		{
			name: "the zero Date is an absent date",
			date: Date{},
			want: true,
		},
		{
			name: "the date is real and non-zero",
			date: DateFromInt(2026, time.May, 10),
			want: false,
		},
		{
			name: "empty time.Time is an absent date",
			date: DateFromTime(time.Time{}),
			want: true,
		},
		{
			name: "january 1st, 0001 is the zero Date",
			date: DateFromInt(1, time.January, 1),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.date.IsZero()
			if got != tt.want {
				t.Errorf("%v.IsZero() = %v, want %v", tt.date, got, tt.want)
			}
		})
	}
}
