package sequence

import (
	"reflect"
	"slices"
	"sync"
)

// FieldSequence defines the sequence of field displaying
type FieldSequence struct {
	lock      sync.RWMutex
	sequences map[string][]string
}

// NewFieldSequence creates new field sequence
func NewFieldSequence() *FieldSequence {
	return &FieldSequence{
		sequences: map[string][]string{},
	}
}

// AddSequenceByTypeName adds sequence by type name of the struct
func (fieldSeq *FieldSequence) AddSequenceByTypeName(
	typeName string,
	seq []string,
) {
	fieldSeq.lock.Lock()
	defer fieldSeq.lock.Unlock()
	fieldSeq.sequences[typeName] = slices.Clone(seq)
}

// AddSequence adds sequence for a type
func (fieldSeq *FieldSequence) AddSequence(
	typeInfo reflect.Type,
	seq []string,
) {
	typeName := typeInfo.Name()
	fieldSeq.AddSequenceByTypeName(
		typeName,
		seq,
	)
}

// GetSequenceByTypeName gets the sequence for a type by type name.
// Returns nil if no such sequence of a type exists
func (fieldSeq *FieldSequence) GetSequenceByTypeName(
	typeName string,
) []string {
	sequenceFetched, exists := fieldSeq.sequences[typeName]
	if !exists {
		return nil
	}
	return slices.Clone(sequenceFetched)
}

// GetSequence gets the sequence for a type.
// Returns nil if no such sequence of a type exists
func (fieldSeq *FieldSequence) GetSequence(
	typeInfo reflect.Type,
) []string {
	typeName := typeInfo.Name()
	return fieldSeq.GetSequenceByTypeName(
		typeName,
	)
}
