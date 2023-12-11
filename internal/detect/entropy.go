package detect

import (
	"math"
)

// shannonEntropy returns the entropy of the given string.
func shannonEntropy(raw string) float64 {
	// Count characters
	counts := make(map[rune]int)

	for _, c := range raw {
		counts[c]++
	}

	// Compute entropy
	var entropy float64

	for _, c := range counts {
		freq := float64(c) / float64(len(raw))
		entropy -= freq * math.Log2(freq)
	}

	return entropy
}
