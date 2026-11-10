package display

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/1Vewton/structdisplay/pkg/ignore"
	"github.com/1Vewton/structdisplay/pkg/process"
	"github.com/1Vewton/structdisplay/pkg/sequence"
)

// TestStruct defines the data struct for testing
type TestStruct struct {
	Field1 string
	field2 int
	Field3 int
	Field4 float64
}

// TestStructMid defines the test struct
type TestStructMid struct {
}

// TestStruct2 defines the test struct
type TestStruct2 struct {
	Field1       string
	Field2       *TestStruct
	FieldSpecial *TestStructMid
}

// TestStructWithIgnore defines the test struct with ignore
type TestStructWithIgnore struct {
	FieldBBC        string
	FieldA          int
	FieldTestStruct *TestStruct
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
	newSequence := sequence.NewFieldSequence().AddSequence(
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

// TestNestedStruct tests the display function on nested data struct
func TestNestedStruct(t *testing.T) {
	t.Parallel()
	// Defines data
	testStruct := &TestStruct2{
		Field1: "114514",
		Field2: &TestStruct{
			Field1: "114514",
			field2: 114514,
			Field3: 1919810,
			Field4: 3.14,
		},
		FieldSpecial: nil,
	}
	// Process nil method
	testProcessNilMethod := func(layer int) string {
		const tab string = "    "
		tabs := strings.Repeat(tab, layer)
		return tabs + "- FieldSpecial does not exists"
	}
	nilToolSet := process.NewNilToolset().AddProcessMethod(
		"FieldSpecial",
		testProcessNilMethod,
	)
	// Converter
	testConverter := NewConverter().WithProcessNilToolset(
		nilToolSet,
	).WithRecursionLimit(
		99,
	)
	_, err := testConverter.ConvertAnyData(
		"testConvert",
		"",
		testStruct,
		0,
	)
	if err != nil {
		t.Fatal(err)
	}
}

// TestIgnore tests the ignoring when performing display with ignore terms
func TestIgnore(t *testing.T) {
	t.Parallel()
	const expectedResult string = `- testIgnoreData: 
    - FieldA: 1 
`
	ignoredSlice := ignore.NewTermsIgnoredSlice()
	ignoredSlice.AddField(
		"FieldBBC",
	)
	subData := &TestStruct{
		Field1: "114514",
		field2: 7,
		Field3: 1919810,
		Field4: 3.14,
	}
	testIgnoreData := &TestStructWithIgnore{
		FieldBBC:        "114514",
		FieldA:          1,
		FieldTestStruct: subData,
	}
	ignoredSlice.AddTypeByType(
		reflect.TypeOf(
			*subData,
		),
	)
	testConverter := NewConverter().WithIgnoredTerms(
		ignoredSlice,
	)
	displayedResult, err := testConverter.ConvertAnyData(
		"testIgnoreData",
		"",
		testIgnoreData,
		0,
	)
	if err != nil {
		t.Fatal(err)
	}
	if displayedResult != expectedResult {
		t.Fatalf(
			"expected:\n%s\ngot:\n%s\n",
			expectedResult,
			displayedResult,
		)
	}
}
