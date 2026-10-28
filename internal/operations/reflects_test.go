package operations

import (
	"reflect"
	"testing"
)

// TestGetField tests the field getting
func TestGetField(t *testing.T) {
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
