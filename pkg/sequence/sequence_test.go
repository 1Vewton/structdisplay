package sequence

import (
	"reflect"
	"testing"
)

// TestStruct defines the data struct for testing
type TestStruct struct {
	Field1 int
	Field2 int
	Field3 int
}

// TestGetSequence tests the sequence getting
func TestGetSequence(t *testing.T) {
	t.Parallel()
	newFieldSequence := NewFieldSequence()
	testSequence := []string{
		"Field1",
		"Field2",
		"Field3",
	}
	newFieldSequence.AddSequenceByTypeName(
		"TestStruct",
		testSequence,
	)
	tempTestStruct := TestStruct{
		Field1: 1,
		Field2: 2,
		Field3: 3,
	}
	testStructType := reflect.TypeOf(tempTestStruct)
	fetchedSequence := newFieldSequence.GetSequence(testStructType)
	if fetchedSequence == nil {
		t.Fatal(
			"the sequence should be fetched successfully",
		)
	}
	for i, fieldName := range fetchedSequence {
		if len(testSequence) <= i {
			t.Fatalf(
				"%d is higher or equals to the length (%d) of the expected sequence",
				i,
				len(testSequence),
			)
		}
		if fieldName != testSequence[i] {
			t.Fatalf(
				"on index %d, the original field name %s is not equals to the field name fetched %s",
				i,
				testSequence[i],
				fieldName,
			)
		}
	}
}
