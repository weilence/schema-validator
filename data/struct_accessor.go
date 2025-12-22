package data

import (
	"errors"
	"reflect"

	reflectutil "github.com/weilence/schema-validator/reflect"
)

type structAccessor struct {
	value      reflect.Value
	structInfo *reflectutil.StructInfo
}

func newStructAccessor(v reflect.Value) *structAccessor {
	return &structAccessor{
		value:      v,
		structInfo: reflectutil.GetStructInfo(v.Type()),
	}
}

// Kind returns KindObject
func (s *structAccessor) Kind() DataKind {
	return KindObject
}

// IsNil always returns false for structs
func (s *structAccessor) IsNil() bool {
	return !s.value.IsValid()
}

// AsField returns error
func (s *structAccessor) AsField() (FieldAccessor, error) {
	return nil, errors.New("struct is not a field")
}

// AsObject returns itself
func (s *structAccessor) AsObject() (ObjectAccessor, error) {
	return s, nil
}

// AsArray returns error
func (s *structAccessor) AsArray() (ArrayAccessor, error) {
	return nil, errors.New("struct is not an array")
}

// GetField returns field by name (supports embedded fields)
func (s *structAccessor) GetField(name string) (Accessor, bool) {
	// Use cached struct info to find field (including embedded)
	fieldInfo, ok := s.structInfo.GetField(name)
	if !ok {
		return nil, false
	}

	// Access field value using reflection (handles embedded & private)
	fieldValue := fieldInfo.GetValue(s.value)

	if !fieldValue.IsValid() {
		return nil, false
	}

	return NewValue(fieldValue.Interface()), true
}

// Fields returns all field names
func (s *structAccessor) Fields() []string {
	return s.structInfo.FieldNames()
}

// Len returns number of fields
func (s *structAccessor) Len() int {
	return len(s.structInfo.FieldNames())
}
