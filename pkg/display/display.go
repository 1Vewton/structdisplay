package display

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/1Vewton/structdisplay/pkg/template"
)

// Converter defines the converter for displaying string
type Converter[T template.Template] struct {
	recursionLimit *int
	templateUsed   T
}

// NewConverter creates new converter
func NewConverter[T template.Template]() *Converter[T] {
	return &Converter[T]{}
}

// WithRecursionLimit sets recursion limit
func (converter *Converter[T]) WithRecursionLimit(
	recursionLimit int,
) *Converter[T] {
	converter.recursionLimit = &recursionLimit
	return converter
}

// ConvertAnyData converts any data to string
func (converter *Converter[T]) ConvertAnyData(
	title string,
	baseString string,
	data any,
	layer int,
) (string, error) {
	resultString := baseString
	if data == nil {
		return resultString, errors.New(
			"you cannot input a nil data",
		)
	}
	if converter.recursionLimit != nil {
		if layer >= *converter.recursionLimit {
			return resultString, fmt.Errorf(
				"the recursion limit %d is reached",
				*converter.recursionLimit,
			)
		}
	}
	// Checks if this is a data struct
	dType := reflect.TypeOf(data)
	dValue := reflect.ValueOf(data)
	if dType.Kind() != reflect.Struct {
		resultString += converter.templateUsed.DisplayField(
			layer,
			title,
			dValue,
		)
	} else {
		resultString += converter.templateUsed.DisplayTitle(
			layer,
			title,
		)
	}
	resultString += converter.templateUsed.DisplayTail(
		layer,
	)
	return resultString, nil
}
