package display

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/1Vewton/structdisplay/pkg/template"
	"github.com/1Vewton/structdisplay/pkg/template/templatedefault"
)

// Converter defines the converter for displaying string
type Converter[T template.Template] struct {
	recursionLimit *int
	templateUsed   T
}

// NewConverter creates new converter
func NewConverter() *Converter[template.Template] {
	result := &Converter[template.Template]{}
	result.WithTemplate(
		templatedefault.NewDefaultTemplate(),
	)
	return result
}

// WithTemplate defines template for converter
func (converter *Converter[T]) WithTemplate(
	temp T,
) *Converter[T] {
	converter.templateUsed = temp
	return converter
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
	resultString += converter.templateUsed.DisplayTitle(
		layer,
		title,
	)
	if dType.Kind() != reflect.Struct {
		resultString += converter.templateUsed.DisplayField(
			layer,
			title,
			dValue,
		)
	} else {
		// Iterating the fields
		numFields := dType.NumField()
		for i := range numFields {
			field := dType.Field(i)
			if !field.IsExported() {
				continue
			}
			fieldByName := dValue.FieldByName(
				field.Name,
			)
			isValid := fieldByName.IsValid()
			if !isValid {
				return resultString, fmt.Errorf(
					"%s field is not valid",
					field.Name,
				)
			}
			canInterface := fieldByName.CanInterface()
			if !canInterface {
				return resultString, fmt.Errorf(
					"%s field cannot be queried",
					field.Name,
				)
			}
			if field.Type.Kind() == reflect.Struct {
				subString, err := converter.ConvertAnyData(
					field.Name,
					resultString,
					fieldByName.Interface(),
					layer+1,
				)
				resultString = subString
				if err != nil {
					return resultString, err
				}
			} else {
				resultString += converter.templateUsed.DisplayField(
					layer,
					field.Name,
					fieldByName.Interface(),
				)
			}
		}
	}
	resultString += converter.templateUsed.DisplayTail(
		layer,
	)
	return resultString, nil
}
