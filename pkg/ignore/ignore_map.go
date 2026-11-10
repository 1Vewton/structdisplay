package ignore

import (
	"reflect"
	"sync"
)

// TermsIgnoredMap stores terms ignored in map
type TermsIgnoredMap struct {
	lock   sync.RWMutex
	fields map[string]int
	types  map[string]int
}

// NewTermsIgnoredMap creates new terms ignored map
func NewTermsIgnoredMap() *TermsIgnoredMap {
	return &TermsIgnoredMap{
		fields: map[string]int{},
		types:  map[string]int{},
	}
}

// AddField adds new fields that will be ignored
func (termsMap *TermsIgnoredMap) AddField(
	fieldNames ...string,
) {
	termsMap.lock.Lock()
	defer termsMap.lock.Unlock()
	for _, fieldName := range fieldNames {
		termsMap.fields[fieldName] = 1
	}
}

// IsFieldIgnored checks if certain field is ignored
func (termsMap *TermsIgnoredMap) IsFieldIgnored(
	fieldName string,
) bool {
	termsMap.lock.RLock()
	defer termsMap.lock.RUnlock()
	_, exists := termsMap.fields[fieldName]
	return exists
}

// AddType adds new type that will get ignored
func (termsMap *TermsIgnoredMap) AddType(
	typeNames ...string,
) {
	termsMap.lock.Lock()
	defer termsMap.lock.Unlock()
	for _, typeName := range typeNames {
		termsMap.types[typeName] = 1
	}
}

// AddTypeByType adds types by reflect type
func (termsMap *TermsIgnoredMap) AddTypeByType(
	types ...reflect.Type,
) {
	stringTypes := []string{}
	for _, typeSelected := range types {
		stringTypes = append(stringTypes, typeSelected.Name())
	}
	termsMap.AddType(stringTypes...)
}

// IsTypeIgnored checks if certain type is ignore
func (termsMap *TermsIgnoredMap) IsTypeIgnored(
	typeName string,
) bool {
	termsMap.lock.RLock()
	defer termsMap.lock.RUnlock()
	_, exists := termsMap.types[typeName]
	return exists
}

// CheckIsTypeIgnoredByType checks if certain reflect type is ignore
func (termsMap *TermsIgnoredMap) CheckIsTypeIgnoredByType(
	reflectType reflect.Type,
) bool {
	return termsMap.IsTypeIgnored(
		reflectType.Name(),
	)
}
