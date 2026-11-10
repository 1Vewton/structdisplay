package ignore

import (
	"reflect"
	"testing"
)

// TestStruct1 defines the struct for testing
type TestStruct1 struct {
	Field1 string
}

// TestStruct4 defines the struct for testing
type TestStruct4 struct {
	Field4 string
}

// TestMapInterface tests whether the ignore map implements the interface
func TestMapInterface(t *testing.T) {
	t.Parallel()
	var anyTerms any = NewTermsIgnoredMap()
	_, ok := anyTerms.(TermsIgnored)
	if !ok {
		t.Fatal(
			"the map does not implementes the ignore interface",
		)
	}
}

// TestIgnoreMap tests functions of ignore map
func TestIgnoreMap(t *testing.T) {
	t.Parallel()
	termsIgnoredMap := NewTermsIgnoredMap()
	newTestStruct1 := TestStruct1{
		Field1: "114514",
	}
	typeOfNewTestStruct1 := reflect.TypeOf(
		newTestStruct1,
	)
	termsIgnoredMap.AddTypeByType(
		typeOfNewTestStruct1,
	)
	isTypeTestStruct1Ignored := termsIgnoredMap.CheckIsTypeIgnoredByType(
		typeOfNewTestStruct1,
	)
	if !isTypeTestStruct1Ignored {
		t.Fatal(
			"expected type test struct 1 to get ignored. ",
		)
	}
	newTestStruct4 := TestStruct4{
		Field4: "114514",
	}
	typeOfTestStruct4 := reflect.TypeOf(
		newTestStruct4,
	)
	isTestStruct4Ignored := termsIgnoredMap.CheckIsTypeIgnoredByType(
		typeOfTestStruct4,
	)
	if isTestStruct4Ignored {
		t.Fatalf(
			"type %s is ignored",
			typeOfTestStruct4.Name(),
		)
	}
}
