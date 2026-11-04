package process

import (
	"fmt"
	"testing"
)

// testTestFieldProcessNilMethod
func testTestFieldProcessNilMethod(
	layer int,
) string {
	return fmt.Sprintf(
		"- %d",
		layer,
	)
}

// TestProcessNil tests the nil processing
func TestProcessNil(t *testing.T) {
	t.Parallel()
	const testField = "testField"
	set := NewNilToolset()
	set.AddProcessMethod(
		testField,
		testTestFieldProcessNilMethod,
	)
	expectedResult := testTestFieldProcessNilMethod(
		114514,
	)
	result := set.ProcessNil(
		testField,
		114514,
	)
	if result == nil {
		t.Fatal(
			"expected the result not to be nil",
		)
	}
	if *result != expectedResult {
		t.Fatalf(
			"expected %s, got %s",
			expectedResult,
			*result,
		)
	}
	result2 := set.ProcessNil(
		"nilField",
		114514,
	)
	if result2 != nil {
		t.Fatal(
			"expected the result 2 to be nil",
		)
	}
}
