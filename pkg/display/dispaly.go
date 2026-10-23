package display

import (
	"errors"
	"reflect"
)

// DisplayAnyData dispalys any data in string form
func DisplayAnyData(
	baseString string,
	data any,
	layer int,
) (string, error) {
	if data == nil {
		return "", errors.New(
			"you cannot input a nil data",
		)
	}
	// Checks if this is a data struct
	dType := reflect.TypeOf(data)
	if dType.Kind() != reflect.Struct {

	}
	reflect.ValueOf(data)
	return "", nil
}
