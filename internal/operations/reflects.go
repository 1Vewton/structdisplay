package operations

import (
	"errors"
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
	isNil := data == nil
	if isNil {
		return &FetchDataResult{
			IsNil: isNil,
		}, nil
	}
	dType := reflect.TypeOf(data)
	dValue := reflect.ValueOf(data)
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

// GetFields gets the fields for a data struct.
// Returns nil if it is not a data struct
func GetFields(
	data any,
	maximumSteps *int,
) (map[string]*FetchDataResult, error) {
	result := map[string]*FetchDataResult{}
	dataInfo, err := FetchData(
		data,
		0,
		maximumSteps,
	)
	if err != nil {
		return nil, err
	}
	if dataInfo.IsNil {
		return nil, errors.New(
			"you should not pass a nil value in",
		)
	}
	if dataInfo.Type.Kind() != reflect.Struct {
		return nil, nil
	}
	fieldsNum := dataInfo.Type.NumField()
	for i := range fieldsNum {
		field := dataInfo.Type.Field(i)
		if !field.IsExported() {
			continue
		}
		fieldByName := dataInfo.Value.FieldByName(
			field.Name,
		)
		if !fieldByName.CanInterface() {
			return nil, fmt.Errorf(
				"data in %s field cannot be interfaced",
				field.Name,
			)
		}
		fieldData := fieldByName.Interface()
		fieldResult, err := FetchData(
			fieldData,
			0,
			maximumSteps,
		)
		if err != nil {
			return nil, err
		}
		result[field.Name] = fieldResult
	}
	return result, nil
}
