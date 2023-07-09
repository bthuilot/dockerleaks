package common

import "math"

// CalculateShannonEntropy will calculate the shannon entropy of a string
func CalculateShannonEntropy(s string) (entropy float64) {
	if len(s) == 0 {
		return
	}

	freq := make(map[rune]float64)
	for _, r := range s {
		freq[r]++
	}

	// normalize frequencies and calculate entropy
	for _, f := range freq {
		f /= float64(len(s))
		entropy += f * math.Log2(f)
	}

	entropy = -entropy
	return
}
