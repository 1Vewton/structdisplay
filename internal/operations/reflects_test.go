package operations

import (
	"reflect"
	"testing"
)

// TestStructPre defines the 1st test struct
type TestStructPre struct {
	Prefield1 int
	pre       int
}

// TestStruct defines the struct for testing
type TestStruct struct {
	Field1    int
	Field2    int
	Field3    int
	preField4 int
	StructPre *TestStructPre
}

// TestFetchData tests the data fetching from pointer
func TestFetchData(t *testing.T) {
	t.Parallel()
	// Raw Data
	data1 := 114514
	result1, err := FetchData(data1, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result1.Type.Kind() != reflect.Int {
		t.Fatalf(
			"exxpected %s, got %s",
			result1.Type.Kind().String(),
			reflect.Int.String(),
		)
	}
	if !result1.Value.CanInterface() {
		t.Fatalf(
			"cannot interact with this result",
		)
	}
	interfaceResult1 := result1.Value.Interface()
	intData1, ok := interfaceResult1.(int)
	if !ok {
		t.Fatal(
			"the result 1 interfaced value ought to be int",
		)
	}
	if intData1 != 114514 {
		t.Errorf(
			"expected the int data fetched to be %d, got %d",
			114514,
			intData1,
		)
	}
	// 1 Pointer
	data2 := &data1
	result2, err := FetchData(data2, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result2.Type.Kind() != reflect.Int {
		t.Fatalf(
			"exxpected %s, got %s",
			result2.Type.Kind().String(),
			reflect.Int.String(),
		)
	}
	if !result2.Value.CanInterface() {
		t.Fatal(
			"the value of result2 should can be interfaced",
		)
	}
	interfaceResult2 := result2.Value.Interface()
	intData2, ok := interfaceResult2.(int)
	if !ok {
		t.Fatal(
			"the result2 interfaced value ought to give a int",
		)
	}
	if intData2 != 114514 {
		t.Fatalf(
			"expected result2 interfaced value to be %d, got %d",
			114514,
			intData2,
		)
	}
	// 2 Pointer
	data3 := &data2
	result3, err := FetchData(data3, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result3.Type.Kind() != reflect.Int {
		t.Fatalf(
			"exxpected %s, got %s",
			result3.Type.Kind().String(),
			reflect.Int.String(),
		)
	}
	if !result3.Value.CanInterface() {
		t.Fatal(
			"result3 value ought to be able to get interfaced",
		)
	}
	interfaceResult3 := result3.Value.Interface()
	intData3, ok := interfaceResult3.(int)
	if !ok {
		t.Fatal(
			"the interfaced value of result 3 ought to have type int",
		)
	}
	if intData3 != 114514 {
		t.Fatalf(
			"expected the interfaced data tobe %d, got %d",
			114514,
			intData3,
		)
	}
}

// TestFetchDataFromStruct tests the data fetching from pointer to a data struct
func TestFetchDataFromStruct(t *testing.T) {
	t.Parallel()
	testData := &TestStruct{
		Field1: 114514,
	}
	structResult1, err := FetchData(
		testData,
		0,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	if structResult1.IsNil {
		t.Fatal(
			"this struct ought not to be nil",
		)
	}
	if structResult1.Type.Kind() != reflect.Struct {
		t.Fatalf(
			"expected type to be %s, got %s",
			reflect.Struct.String(),
			structResult1.Type.String(),
		)
	}
	if !structResult1.Value.CanInterface() {
		t.Fatal(
			"the value of the struct result 1 expected to be able to interface",
		)
	}
	interfaceStructResult1 := structResult1.Value.Interface()
	structResultData1, ok := interfaceStructResult1.(TestStruct)
	if !ok {
		t.Fatal(
			"expected struct result interfaced to have TestStruct type",
		)
	}
	if structResultData1.Field1 != 114514 {
		t.Errorf(
			"expected the field1 of the struct result to be %d, got %d",
			114514,
			structResultData1.Field1,
		)
	}
}

// TestFetchDataFromNil tests fetching data from nil value
func TestFetchDataFromNil(t *testing.T) {
	t.Parallel()
	res, err := FetchData(
		nil,
		0,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsNil {
		t.Fatal(
			"the result is expected to find the data is nil",
		)
	}
}

// TestGetFields tests the function of getting fields from a data struct
func TestGetFields(t *testing.T) {
	t.Parallel()
	testData := &TestStruct{
		Field1:    1,
		Field2:    2,
		Field3:    3,
		preField4: 1,
		StructPre: &TestStructPre{
			Prefield1: 1,
			pre:       1,
		},
	}
	fieldsFetched, err := GetFields(
		testData,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	field1, exists := fieldsFetched["Field1"]
	if !exists {
		t.Fatal(
			"Field 1 expected to get fetched",
		)
	}
	if field1.Type.Kind() != reflect.Int {
		t.Fatalf(
			"expected type %s, got %s",
			field1.Type.Kind().String(),
			reflect.Int.String(),
		)
	}
	if !field1.Value.CanInterface() {
		t.Fatal(
			"cannot interface field field1",
		)
	}
	_, exists = fieldsFetched["preField4"]
	if exists {
		t.Fatal(
			"expected not to fetch a field that is a private field",
		)
	}
	structPreField, exists := fieldsFetched["StructPre"]
	if !exists {
		t.Fatal(
			"field StructPre expected to get fetched",
		)
	}
	if structPreField.Type.Kind() != reflect.Struct {
		t.Fatalf(
			"expected StructPre to have kind %s, got %s",
			reflect.Struct.String(),
			structPreField.Type.Kind().String(),
		)
	}
}
