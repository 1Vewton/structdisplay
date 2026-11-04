package process

import (
	"sync"
)

// NilMethod defines the method for processing
type NilMethod func(
	layer int,
) string

// NilToolset defines the set of tools for processing nil
type NilToolset struct {
	lock sync.RWMutex
	set  map[string]NilMethod
}

// NewNilToolset creates new toolset for processing nil value
func NewNilToolset() *NilToolset {
	return &NilToolset{
		set: map[string]NilMethod{},
	}
}

// AddProcessMethod adds process method for a field name
func (toolset *NilToolset) AddProcessMethod(
	fieldName string,
	method NilMethod,
) {
	toolset.lock.Lock()
	defer toolset.lock.Unlock()
	toolset.set[fieldName] = method
}

// ProcessNil processes nil
func (toolset *NilToolset) ProcessNil(
	fieldName string,
	layer int,
) *string {
	toolset.lock.RLock()
	defer toolset.lock.RUnlock()
	method, exists := toolset.set[fieldName]
	if !exists {
		return nil
	}
	result := method(layer)
	return &result
}
