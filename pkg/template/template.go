package template

// Template defines template for displaying
type Template interface {
	DisplayField(
		layer int,
		fieldName string,
		value any,
	) string
	DisplayTitle(
		layer int,
		title string,
	) string
	DisplayTail(
		layer int,
	) string
}
