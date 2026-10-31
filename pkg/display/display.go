package display

import (
	"fmt"
	"reflect"

	"github.com/1Vewton/structdisplay/internal/operations"
	"github.com/1Vewton/structdisplay/pkg/sequence"
	"github.com/1Vewton/structdisplay/pkg/template"
	"github.com/1Vewton/structdisplay/pkg/template/templatedefault"
)

// Converter defines the converter for displaying string
type Converter[T template.Template] struct {
	recursionLimit *int
	templateUsed   T
	sequenceUsed   *sequence.FieldSequence
}

// NewConverter creates new converter
func NewConverter() *Converter[template.Template] {
	result := &Converter[template.Template]{}
	result.WithTemplate(
		templatedefault.NewDefaultTemplate(),
	)
	return result
}

// WithFieldSequence sets the field sequence for converter
func (converter *Converter[T]) WithFieldSequence(
	sequenceUsed *sequence.FieldSequence,
) *Converter[T] {
	converter.sequenceUsed = sequenceUsed
	return converter
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
	fetchedData, err := operations.FetchData(
		data,
		layer,
		converter.recursionLimit,
	)
	if err != nil {
		return resultString, err
	}
	if !fetchedData.Value.CanInterface() {
		return resultString, fmt.Errorf(
			"%s data cannot get interfaced",
			title,
		)
	}
	if fetchedData.Type.Kind() == reflect.Struct {
		// Process data struct
		resultString += converter.templateUsed.DisplayTitle(
			layer,
			title,
		)
		fields, err := operations.GetFields(
			fetchedData.Value.Interface(),
			converter.recursionLimit,
		)
		if err != nil {
			return resultString, err
		}
		datasForAssembling := map[string]string{}
		for fieldName, field := range fields {
			if !field.Value.CanInterface() {
				return resultString, fmt.Errorf(
					"%s cannot get interfaced",
					fieldName,
				)
			}
			subString, err := converter.ConvertAnyData(
				fieldName,
				"",
				field.Value.Interface(),
				layer+1,
			)
			if err != nil {
				return resultString, err
			}
			datasForAssembling[fieldName] = subString
		}
		if converter.sequenceUsed != nil {
			fieldSequence := converter.sequenceUsed.GetSequence(
				fetchedData.Type,
			)
			for _, field := range fieldSequence {
				fieldString, exists := datasForAssembling[field]
				if !exists {
					continue
				}
				resultString += fieldString
				delete(datasForAssembling, field)
			}
		}
		for _, rest := range datasForAssembling {
			resultString += rest
		}
		resultString += converter.templateUsed.DisplayTail(
			layer,
		)
		return resultString, nil
	}
	resultString += converter.templateUsed.DisplayField(
		layer,
		title,
		fetchedData.Value.Interface(),
	)
	return resultString, nil
}
