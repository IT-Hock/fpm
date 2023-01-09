package utils

import "testing"

func TestClamp(t *testing.T) {
	tests := []struct {
		value, min, max, want int
	}{
		{0, 1, 10, 1},
		{5, 1, 10, 5},
		{15, 1, 10, 10},
	}
	for _, test := range tests {
		if got := Clamp(test.value, test.min, test.max); got != test.want {
			t.Errorf("Clamp(%d, %d, %d) = %d", test.value, test.min, test.max, got)
		}
	}
}
