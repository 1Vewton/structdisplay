package templatedefault

import (
	"fmt"
	"strings"
)

const tab string = "    "

// DefaultTemplate defines the default template
type DefaultTemplate struct {
}

// NewDefaultTemplate creates new default template
func NewDefaultTemplate() *DefaultTemplate {
	return &DefaultTemplate{}
}

// DisplayField defines the default function that displays a field
func (template *DefaultTemplate) DisplayField(
	layer int,
	fieldName string,
	value any,
) string {
	resultString := strings.Repeat(
		tab,
		layer,
	)
	resultString += fmt.Sprintf(
		"- %s: %v \n",
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
	resultString := strings.Repeat(
		tab,
		max(layer-1, 0),
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
