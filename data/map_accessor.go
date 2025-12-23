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

func (m *mapAccessor) Raw() any {
	return m.value.Interface()
}

func (m *mapAccessor) GetValue(path string) (*Value, error) {
	if path == "" {
		return nil, fmt.Errorf("empty path")
	}

	fieldName, nextPath := cutPath(path)

	fieldAcc, err := m.GetField(fieldName)
	if err != nil {
		return nil, err
	}

	return fieldAcc.GetValue(nextPath)
}

// GetField returns field by name (map key)
func (m *mapAccessor) GetField(name string) (Accessor, error) {
	// Convert name to key value
	keyVal := reflect.ValueOf(name)

	// Get value from map
	val := m.value.MapIndex(keyVal)
	if !val.IsValid() {
		return nil, fmt.Errorf("key %s not found in map", name)
	}

	return NewAccessor(val), nil
}

func (m *mapAccessor) Accessors() []ObjectAccessor {
	accessors := []ObjectAccessor{m}
	return accessors
}
