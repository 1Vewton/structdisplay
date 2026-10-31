package templatedefault

import (
	"testing"

	"github.com/1Vewton/structdisplay/pkg/template"
)

// TestDefaultDisplayFieldFunc tests the default function for displaying a field
func TestDefaultDisplayFieldFunc(t *testing.T) {
	t.Parallel()
	const expectedResult string = "- Data: 5 \n"
	defaultTemplate := NewDefaultTemplate()
	fetchedResult := defaultTemplate.DisplayField(
		0,
		"Data",
		5,
	)
	if expectedResult != fetchedResult {
		t.Errorf(
			"expected %s, got %s",
			expectedResult,
			fetchedResult,
		)
	}
	const expectedResultTab1 string = "    - Test: ggds \n"
	fetchedResult = defaultTemplate.DisplayField(
		1,
		"Test",
		"ggds",
	)
	if expectedResultTab1 != fetchedResult {
		t.Errorf(
			"expected %s, got %s",
			expectedResultTab1,
			fetchedResult,
		)
	}
	const expectedResultTab2 string = "        - Test: ggds \n"
	fetchedResult = defaultTemplate.DisplayField(
		2,
		"Test",
		"ggds",
	)
	if expectedResultTab2 != fetchedResult {
		t.Errorf(
			"expected %s, got %s",
			expectedResultTab2,
			fetchedResult,
		)
	}
}

// TestDefaultDisplayTitleFunc tests the default function for displaying title
func TestDefaultDisplayTitleFunc(t *testing.T) {
	t.Parallel()
	defaultTemplate := NewDefaultTemplate()
	const expectedResult1 string = "- abc: \n"
	result1 := defaultTemplate.DisplayTitle(
		0,
		"abc",
	)
	if result1 != expectedResult1 {
		t.Errorf(
			"expected %s, got %s",
			expectedResult1,
			result1,
		)
	}
	const expectedResult2 string = "        - abc: \n"
	result2 := defaultTemplate.DisplayTitle(
		2,
		"abc",
	)
	if result2 != expectedResult2 {
		t.Errorf(
			"expected %s, got %s",
			expectedResult2,
			result2,
		)
	}
}

// TestInterface tests whether it can adapt to the template interface
func TestInterface(t *testing.T) {
	t.Parallel()
	var defaultTemplate interface{} = NewDefaultTemplate()
	_, ok := defaultTemplate.(template.Template)
	if !ok {
		t.Error(
			"this cannot fit to the template interface",
		)
	}
}
