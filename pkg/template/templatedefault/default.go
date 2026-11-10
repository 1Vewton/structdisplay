package templatedefault

import (
	"fmt"
	"reflect"
	"strings"
)

const tab string = "    "

// DefaultTemplate defines the default template
type DefaultTemplate struct {
	recursionLimit *int
}

// NewDefaultTemplate creates new default template
func NewDefaultTemplate(
	recursionLimit *int,
) *DefaultTemplate {
	return &DefaultTemplate{}
}

// DisplayField defines the default function that displays a field
func (template *DefaultTemplate) DisplayField(
	layer int,
	fieldName any,
	value any,
) string {
	if template.recursionLimit != nil {
		if layer > *template.recursionLimit {
			return ""
		}
	}
	typeOfValue := reflect.TypeOf(
		value,
	)
	usedLayer := max(layer, 0)
	resultString := strings.Repeat(
		tab,
		usedLayer,
	)
	// You can set display rules for field by mapping the type of field
	if typeOfValue.Kind() == reflect.Map {
		valueMap, ok := value.(map[any]any)
		if !ok {
			resultString += fmt.Sprintf(
				"- %s: %v \n",
				fieldName,
				value,
			)
		} else {
			resultString += fmt.Sprintf(
				"- %v: \n",
				fieldName,
			)
			for key, value := range valueMap {
				resultString += template.DisplayField(
					layer+1,
					key,
					value,
				)
			}
		}
		return resultString
	}
	resultString += fmt.Sprintf(
		"- %v: %v \n",
		fieldName,
		value,
	)
	return resultString
}

// DisplayTitle defines the default function that displays title of a field
func (template *DefaultTemplate) DisplayTitle(
	layer int,
	title string,
) string {
	usedLayer := max(layer, 0)
	resultString := strings.Repeat(
		tab,
		usedLayer,
	)
	resultString += fmt.Sprintf(
		"- %s: \n",
		title,
	)
	return resultString
}

// DisplayTail displayes the tail of the default template
func (template *DefaultTemplate) DisplayTail(
	layer int,
) string {
	return ""
}
