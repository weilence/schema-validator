package data

import (
	"errors"
	"fmt"
	"reflect"
)

type ArrayAccessor struct {
	value reflect.Value
}

func NewArrayAccessor(v reflect.Value) *ArrayAccessor {
	return &ArrayAccessor{value: v}
}

func (s *ArrayAccessor) GetField(name string) (Accessor, error) {
	var idx int
	n, err := fmt.Sscanf(name, "[%d]", &idx)
	if err != nil || n != 1 {
		return nil, errors.New("invalid array index in scan: " + name)
	}

	elemAcc, err := s.GetIndex(idx)
	if err != nil {
		return nil, fmt.Errorf("index %d out of bounds", idx)
	}

	return elemAcc, nil
}

func (s *ArrayAccessor) Raw() any {
	return s.value.Interface()
}

func (s *ArrayAccessor) GetValue(path string) (*Value, error) {
	if path == "" {
		return &Value{rval: s.value}, nil
	}

	part, nextPath := cutPath(path)
	elemAcc, err := s.GetField(part)
	if err != nil {
		return nil, err
	}

	return elemAcc.GetValue(nextPath)
}

// GetIndex returns element at index
func (s *ArrayAccessor) GetIndex(idx int) (Accessor, error) {
	if idx < 0 || idx >= s.value.Len() {
		return nil, fmt.Errorf("index %d out of bounds", idx)
	}

	elem := s.value.Index(idx)
	return NewAccessor(elem), nil
}

// Len returns array length
func (s *ArrayAccessor) Len() int {
	return s.value.Len()
}

// Iterate calls fn for each element
func (s *ArrayAccessor) Iterate(fn func(idx int, elem Accessor) error) error {
	for i := 0; i < s.value.Len(); i++ {
		elem := s.value.Index(i)
		accessor := NewAccessor(elem)
		if err := fn(i, accessor); err != nil {
			return err
		}
	}

	return nil
}
