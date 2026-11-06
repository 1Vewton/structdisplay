package ignore

import (
	"reflect"
)

// TermsIgnored defines the terms that are ignored
type TermsIgnored interface {
	IsFieldIgnored(
		fieldName string,
	) bool
	AddField(
		fieldNames ...string,
	)
	AddType(
		typeNames ...string,
	)
	AddTypeByType(
		types ...reflect.Type,
	)
	IsTypeIgnored(
		typeName string,
	) bool
	CheckIsTypeIgnoredByType(
		reflectType reflect.Type,
	) bool
}
