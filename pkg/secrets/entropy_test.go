package secrets

import "testing"

func TestCalculateShannonEntropy(t *testing.T) {

	var testCases = []struct {
		input        string
		expected     float64
		allowedError float64
	}{
		{"password", 2.75, 0.01},
		{"password123", 3.277, 0.01},
		{"password123!", 3.4183, 0.01},
		{"", 0, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			entropy := CalculateShannonEntropy(tc.input)
			if entropy < tc.expected-tc.allowedError || entropy > tc.expected+tc.allowedError {
				t.Errorf("Expected entropy to be %f, got %f", tc.expected, entropy)
			}
		})
	}
}
