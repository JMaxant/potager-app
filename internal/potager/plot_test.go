package potager

import "testing"

func TestPlot_AreaCm2(t *testing.T) {
	tests := []struct {
		name string
		p    Plot
		want int
	}{
		{
			name: "a 3m x 1.2m plot is 3.6m²",
			p: Plot{
				LengthCm: 300,
				WidthCm:  120,
			},
			want: 36000,
		},
		{
			name: "a plot with only length has not set area",
			p: Plot{
				LengthCm: 100,
				WidthCm:  0,
			},
			want: 0,
		},
		{
			name: "a plot with only width has not set area",
			p: Plot{
				LengthCm: 0,
				WidthCm:  100,
			},
			want: 0,
		},
		{
			name: "a plot with no length or width has not set area",
			p: Plot{
				LengthCm: 0,
				WidthCm:  0,
			},
			want: 0,
		},
		{
			name: "a plot of 100m per 100m is 100 000 000 cm² (max value)",
			p: Plot{
				LengthCm: 10000,
				WidthCm:  10000,
			},
			want: 100000000,
		},
		{
			name: "a plot of 1cm per 1cm is 1cm² (min value)",
			p: Plot{
				LengthCm: 1,
				WidthCm:  1,
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.AreaCm2()
			if got != tt.want {
				t.Errorf("Plot.AreaCm2() = Plot %+v, got %v, want %v", tt.p, got, tt.want)
			}
		})
	}
}
