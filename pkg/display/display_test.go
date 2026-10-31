package display

import (
	"os"
	"reflect"
	"testing"

	"github.com/1Vewton/structdisplay/pkg/sequence"
)

// TestStruct defines the data struct for testing
type TestStruct struct {
	Field1 string
	field2 int
	Field3 int
	Field4 float64
}

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

// TestDefaultConvert tests the converter for default template
func TestDefaultConvert(t *testing.T) {
	t.Parallel()
	const data1FilePath string = "./testdata/data1.md"
	testData := &TestStruct{
		Field1: "Acceed",
		field2: 114514,
		Field3: 114514,
		Field4: 1.14514,
	}
	data1, err := os.ReadFile(
		data1FilePath,
	)
	if err != nil {
		t.Fatal(err)
	}
	stringData1 := string(data1)
	newSequence := sequence.NewFieldSequence()
	newSequence.AddSequence(
		reflect.TypeOf(*testData),
		[]string{
			"Field1",
			"Field3",
			"Field4",
		},
	)
	converter := NewConverter().WithRecursionLimit(
		99,
	).WithFieldSequence(
		newSequence,
	)
	result, err := converter.ConvertAnyData(
		"Data",
		"",
		testData,
		0,
	)
	if err != nil {
		t.Fatal(err)
	}
	if result != stringData1 {
		t.Logf(
			"expected:\n%sgot:\n%s",
			stringData1,
			result,
		)
	}
}
