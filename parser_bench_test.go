package id

import (
	"fmt"
	"testing"
)

// parseAndFormatSlow simulates with fmt.Errorf that causes allocations.
func parseAndFormatSlow(id string) (string, error) {
	switch len(id) {
	case 32:
		for i := 0; i < len(id); i++ {
			c := id[i]
			if !isHexSlow(c) {
				return "", fmt.Errorf("invalid ID: non-hex character")
			}
		}
	case 36:
		// Similar logic but with fmt.Errorf
		return "", fmt.Errorf("invalid ID: simulated old version")
	default:
		return "", fmt.Errorf("invalid ID length")
	}
	return id, nil
}

// isHexSlow simulates range-checking.
func isHexSlow(c byte) bool {
	return (c >= '0' && c <= '9') ||
		(c >= 'a' && c <= 'f') ||
		(c >= 'A' && c <= 'F')
}

func BenchmarkValid(b *testing.B) {
	p := NewIDParser(NewNoOpCache())
	valid := TestIDs["valid-dashed"]

	b.Run("WithCache", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = p.Parse(valid)
		}
	})

	b.Run("WithoutCache", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = parseAndFormatSlow(valid)
		}
	})
}

func BenchmarkInvalid(b *testing.B) {
	p := NewIDParser(NewNoOpCache())
	invalidChar := TestIDs["invalid-char"]

	b.Run("WithCache", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = p.Parse(invalidChar)
		}
	})

	b.Run("WithoutCache", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = parseAndFormatSlow(invalidChar)
		}
	})
}
