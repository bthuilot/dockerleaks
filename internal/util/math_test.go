package util

import "testing"

func TestMinInt(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{
			name: "a is less than b",
			a:    1,
			b:    2,
			want: 1,
		},
		{
			name: "a is greater than b",
			a:    2,
			b:    1,
			want: 1,
		},
		{
			name: "a is equal to b",
			a:    1,
			b:    1,
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinInt(tt.a, tt.b); got != tt.want {
				t.Errorf("MinInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
