package display

import (
	"testing"
)

// TestConverterSetting tests the display function
func TestConverterSetting(t *testing.T) {
	t.Parallel()
	converter := NewConverter().WithRecursionLimit(
		99,
	)
	if converter.recursionLimit == nil {
		t.Fatal(
			"recursion limit setting failed",
		)
	}
	if *converter.recursionLimit != 99 {
		t.Fatalf(
			"expected the recursion limit to be %d, got %d",
			99,
			*converter.recursionLimit,
		)
	}
}
