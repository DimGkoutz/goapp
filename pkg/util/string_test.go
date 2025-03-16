package util

import (
	"regexp"
	"testing"
)

func TestRandString(t *testing.T) {
	t.Run("ValidHexOutput", func(t *testing.T) {
		// Test different lengths
		lengths := []int{1, 5, 10, 20}

		// Regular expression to match hex strings
		hexPattern := regexp.MustCompile("^[0-9a-f]+$")

		for _, length := range lengths {
			result := RandString(length)

			// Verify the length
			if actualLen := len(result); actualLen != length {
				t.Errorf("RandString(%d) returned string of length %d, expected length %d",
					length, actualLen, length)
			}

			// Verify it contains only hex characters
			if !hexPattern.MatchString(result) {
				t.Errorf("RandString(%d) returned non-hex characters: %s",
					length, result)
			}
		}
	})

	t.Run("CharacterDistribution", func(t *testing.T) {
		// Count occurrences of each character in a long string
		length := 10000
		result := RandString(length)

		counts := make(map[byte]int)
		for i := 0; i < len(result); i++ {
			counts[result[i]]++
		}

		// Check that each hex character appears at least once
		for _, c := range "0123456789abcdef" {
			if counts[byte(c)] == 0 {
				t.Errorf("Character '%c' does not appear in the output", c)
			}
		}

		// Verify the randomness quality by checking that character distribution falls within acceptable range
		expected := length / 16 // 16 hex characters
		for c, count := range counts {
			if count < int(float64(expected)*0.7) || count > int(float64(expected)*1.3) {
				t.Errorf("Character '%c' appears %d times, expected around %d", c, count, expected)
			}
		}
	})
}

func BenchmarkRandString(b *testing.B) {
	// Benchmark different sizes of strings
	benchmarks := []struct {
		name string
		size int
	}{
		{"Small", 10},
		{"Medium", 100},
		{"Large", 1000},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				RandString(bm.size)
			}
		})
	}
}
