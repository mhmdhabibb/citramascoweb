package utils

import (
	"regexp"
	"testing"
)

func TestGenerateRandomNumberCode(t *testing.T) {
	lengths := []int{4, 6, 8, 10}
	numericRegex := regexp.MustCompile("^[0-9]+$")

	for _, length := range lengths {
		code, err := GenerateRandomNumberCode(length)
		if err != nil {
			t.Fatalf("failed to generate code for length %d: %v", length, err)
		}

		if len(code) != length {
			t.Errorf("expected length %d, got %d (code: %s)", length, len(code), code)
		}

		if !numericRegex.MatchString(code) {
			t.Errorf("expected only digits, got non-digit characters in code: %s", code)
		}
	}
}
