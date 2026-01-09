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

func (s *ArrayAccessor) deref() reflect.Value {
	v := s.value
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	return v
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
		return NewValueAccessor(s.value), nil
	}

	part, nextPath := cutPath(path)
	elemAcc, err := s.GetField(part)
	if err != nil {
		return nil, err
	}

	return elemAcc.GetValue(nextPath)
}

func (s *ArrayAccessor) GetIndex(idx int) (Accessor, error) {
	v := s.deref()
	if idx < 0 || idx >= v.Len() {
		return nil, fmt.Errorf("index %d out of bounds", idx)
	}

	elem := v.Index(idx)
	return NewAccessor(elem), nil
}

func (s *ArrayAccessor) Len() int {
	return s.deref().Len()
}

func (s *ArrayAccessor) Iterate(fn func(idx int, elem Accessor) error) error {
	v := s.deref()
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		accessor := NewAccessor(elem)
		if err := fn(i, accessor); err != nil {
			return err
		}
	}

	return nil
}
