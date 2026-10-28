package operations

import (
	"fmt"
	"reflect"
)

// FetchDataResult defines the result after fetching data in a field
type FetchDataResult struct {
	Value reflect.Value
	Type  reflect.Type
	IsNil bool
}

// FetchData fetches the info of a variable.
// It will try to find the value and kind pointed if data is a pointer.
func FetchData(
	data any,
	steps int,
	maximumSteps *int,
) (*FetchDataResult, error) {
	if maximumSteps != nil {
		if steps >= *maximumSteps {
			return nil, fmt.Errorf(
				"%d reaches the recursion maximum limitation",
				steps,
			)
		}
	}
	dType := reflect.TypeOf(data)
	dValue := reflect.ValueOf(data)
	isNil := data == nil
	if dType.Kind() == reflect.Pointer {
		if dValue.CanInterface() {
			return FetchData(
				dValue.Elem().Interface(),
				steps+1,
				maximumSteps,
			)
		}
		return nil, fmt.Errorf(
			"type %s cannot be interfaced",
			dType.Name(),
		)
	}
	return &FetchDataResult{
		Value: dValue,
		Type:  dType,
		IsNil: isNil,
	}, nil
}

// Equals checks if two results are the same
func (result *FetchDataResult) Equals(
	another *FetchDataResult,
) bool {
	return result.IsNil == another.IsNil &&
		result.Value.Equal(another.Value) &&
		result.Type.Kind() == another.Type.Kind()
}
