package data

import (
	"errors"
	"reflect"
)

type sliceAccessor struct {
	value reflect.Value
}

func newSliceAccessor(v reflect.Value) *sliceAccessor {
	return &sliceAccessor{value: v}
}

// Kind returns KindArray
func (s *sliceAccessor) Kind() DataKind {
	return KindArray
}

// IsNil checks if the slice/array is nil
func (s *sliceAccessor) IsNil() bool {
	if !s.value.IsValid() {
		return true
	}
	// Arrays are never nil, but slices can be
	if s.value.Kind() == reflect.Slice {
		return s.value.IsNil()
	}
	return false
}

// AsField returns error
func (s *sliceAccessor) AsField() (FieldAccessor, error) {
	return nil, errors.New("array is not a field")
}

// AsObject returns error
func (s *sliceAccessor) AsObject() (ObjectAccessor, error) {
	return nil, errors.New("array is not an object")
}

// AsArray returns itself
func (s *sliceAccessor) AsArray() (ArrayAccessor, error) {
	return s, nil
}

// GetIndex returns element at index
func (s *sliceAccessor) GetIndex(idx int) (Accessor, bool) {
	if s.IsNil() {
		return nil, false
	}

	if idx < 0 || idx >= s.value.Len() {
		return nil, false
	}

	elem := s.value.Index(idx)
	return NewValue(elem.Interface()), true
}

// Len returns array length
func (s *sliceAccessor) Len() int {
	if s.IsNil() {
		return 0
	}
	return s.value.Len()
}

// Iterate calls fn for each element
func (s *sliceAccessor) Iterate(fn func(idx int, elem Accessor) error) error {
	if s.IsNil() {
		return nil
	}

	for i := 0; i < s.value.Len(); i++ {
		elem := s.value.Index(i)
		accessor := NewValue(elem.Interface())
		if err := fn(i, accessor); err != nil {
			return err
		}
	}

	return nil
}
