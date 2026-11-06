package ignore

import (
	"reflect"
	"sync"
)

// TermsIgnoredSlice defines the terms ignored that stored in a slice
type TermsIgnoredSlice struct {
	lock   sync.RWMutex
	fields []string
	types  []string
}

// NewTermsIgnoredSlice defines new terms ignored slice
func NewTermsIgnoredSlice() *TermsIgnoredSlice {
	return &TermsIgnoredSlice{
		fields: []string{},
		types:  []string{},
	}
}

// AddField adds new field to it
func (termsSlice *TermsIgnoredSlice) AddField(
	fieldNames ...string,
) {
	termsSlice.lock.Lock()
	defer termsSlice.lock.Unlock()
	for _, fieldName := range fieldNames {
		termsSlice.fields = append(termsSlice.fields, fieldName)
	}
}

// IsFieldIgnored checks whether a field is ignored
func (termsSlice *TermsIgnoredSlice) IsFieldIgnored(
	fieldName string,
) bool {
	termsSlice.lock.RLock()
	defer termsSlice.lock.RUnlock()
	for _, fieldNamei := range termsSlice.fields {
		if fieldNamei == fieldName {
			return true
		}
	}
	return false
}

// AddType adds new type into the list so that this type can get ignored
func (termsSlice *TermsIgnoredSlice) AddType(
	typeNames ...string,
) {
	termsSlice.lock.Lock()
	defer termsSlice.lock.Unlock()
	for _, typeName := range typeNames {
		termsSlice.types = append(termsSlice.types, typeName)
	}
}

// AddTypeByType adds new type by reflect type
func (termsSlice *TermsIgnoredSlice) AddTypeByType(
	types ...reflect.Type,
) {
	convertedTypes := []string{}
	for _, rawType := range types {
		convertedTypes = append(convertedTypes, rawType.Name())
	}
	termsSlice.AddType(convertedTypes...)
}

// IsTypeIgnored checks if certain type is ignored
func (termsSlice *TermsIgnoredSlice) IsTypeIgnored(
	typeName string,
) bool {
	termsSlice.lock.RLock()
	defer termsSlice.lock.RUnlock()
	for _, ignoredType := range termsSlice.types {
		if ignoredType == typeName {
			return true
		}
	}
	return false
}

// CheckIsTypeIgnoredByType checks if certain type is ignored by reflect type
func (termsSlice *TermsIgnoredSlice) CheckIsTypeIgnoredByType(
	reflectType reflect.Type,
) bool {
	return termsSlice.IsTypeIgnored(
		reflectType.Name(),
	)
}
