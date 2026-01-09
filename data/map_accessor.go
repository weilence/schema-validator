package data

import (
	"fmt"
	"reflect"
)

type mapAccessor struct {
	value reflect.Value
}

func NewMapAccessor(v reflect.Value) *mapAccessor {
	return &mapAccessor{value: v}
}

func (m *mapAccessor) deref() reflect.Value {
	v := m.value
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	return v
}

func (m *mapAccessor) Raw() any {
	return m.value.Interface()
}

func (m *mapAccessor) GetValue(path string) (*Value, error) {
	if path == "" {
		return NewValueAccessor(m.value), nil
	}

	fieldName, nextPath := cutPath(path)

	fieldAcc, err := m.GetField(fieldName)
	if err != nil {
		return nil, err
	}

	return fieldAcc.GetValue(nextPath)
}

func (m *mapAccessor) GetField(name string) (Accessor, error) {
	v := m.deref()
	keyVal := reflect.ValueOf(name)

	val := v.MapIndex(keyVal)
	if !val.IsValid() {
		return nil, fmt.Errorf("key %s not found in map", name)
	}

	return NewAccessor(val), nil
}

func (m *mapAccessor) Accessors() []ObjectAccessor {
	accessors := []ObjectAccessor{m}
	return accessors
}
