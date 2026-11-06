package ignore

import (
	"reflect"
	"testing"
)

// TestStruct defines the data struct for testing
type TestStruct struct {
	Field1 int
}

// TestStruct2 defines the second data struct for testing
type TestStruct2 struct {
	Field2 int
}

// TestInterface checks if ignored slice implements the ignore interface
func TestInterface(t *testing.T) {
	t.Parallel()
	var ignoredSliceInterface any = NewTermsIgnoredSlice()
	_, ok := ignoredSliceInterface.(TermsIgnored)
	if !ok {
		t.Fatal(
			"the ignored slice does not implements ignored terms interface",
		)
	}
}

// TestIgnoreSlice tests functions for ignore slice
func TestIgnoreSlice(t *testing.T) {
	t.Parallel()
	testData := TestStruct{
		Field1: 114514,
	}
	testData2 := TestStruct2{
		Field2: 1919810,
	}
	newTermsIgnoredSlice := NewTermsIgnoredSlice()
	testDataType := reflect.TypeOf(testData)
	newTermsIgnoredSlice.AddTypeByType(
		testDataType,
	)
	isTestDataTypeIgnored := newTermsIgnoredSlice.CheckIsTypeIgnoredByType(
		testDataType,
	)
	if !isTestDataTypeIgnored {
		t.Fatal(
			"the test struct should be ignored",
		)
	}
	testData2Type := reflect.TypeOf(testData2)
	isTestData2Ignored := newTermsIgnoredSlice.CheckIsTypeIgnoredByType(
		testData2Type,
	)
	if isTestData2Ignored {
		t.Fatal(
			"test struct 2 should not get ignored",
		)
	}
	newTermsIgnoredSlice.AddField(
		"abc",
		"efg",
	)
	isABCIgnored := newTermsIgnoredSlice.IsFieldIgnored(
		"abc",
	)
	if !isABCIgnored {
		t.Fatal(
			"abc should get ignored",
		)
	}
	isCotIgnored := newTermsIgnoredSlice.IsFieldIgnored(
		"cot",
	)
	if isCotIgnored {
		t.Fatal(
			"Cot should not get ignored",
		)
	}
}
