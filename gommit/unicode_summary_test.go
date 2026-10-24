package gommit

import (
	"strings"
	"testing"
)

func TestIsValidSummaryLengthCountsUnicodeRunes(t *testing.T) {
	message := strings.Repeat("é", 25)

	if !isValidSummaryLength(25, message) {
		t.Fatal("25 Unicode characters should fit a summary limit of 25")
	}

	if isValidSummaryLength(24, message) {
		t.Fatal("25 Unicode characters should exceed a summary limit of 24")
	}
}
