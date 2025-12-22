package data

import (
	"errors"
	"reflect"
)

type mapAccessor struct {
	value reflect.Value
}

func newMapAccessor(v reflect.Value) *mapAccessor {
	return &mapAccessor{value: v}
}

// Kind returns KindObject
func (m *mapAccessor) Kind() DataKind {
	return KindObject
}

// IsNil checks if the map is nil
func (m *mapAccessor) IsNil() bool {
	return !m.value.IsValid() || m.value.IsNil()
}

// AsField returns error
func (m *mapAccessor) AsField() (FieldAccessor, error) {
	return nil, errors.New("map is not a field")
}

// AsObject returns itself
func (m *mapAccessor) AsObject() (ObjectAccessor, error) {
	return m, nil
}

// AsArray returns error
func (m *mapAccessor) AsArray() (ArrayAccessor, error) {
	return nil, errors.New("map is not an array")
}

// GetField returns field by name (map key)
func (m *mapAccessor) GetField(name string) (Accessor, bool) {
	if m.IsNil() {
		return nil, false
	}

	// Convert name to key value
	keyVal := reflect.ValueOf(name)

	// Get value from map
	val := m.value.MapIndex(keyVal)
	if !val.IsValid() {
		return nil, false
	}

	return NewValue(val.Interface()), true
}

// Fields returns all map keys as strings
func (m *mapAccessor) Fields() []string {
	if m.IsNil() {
		return []string{}
	}

	keys := m.value.MapKeys()
	result := make([]string, 0, len(keys))

	for _, key := range keys {
		// Convert key to string
		keyStr := ""
		if key.Kind() == reflect.String {
			keyStr = key.String()
		} else {
			keyStr = key.String() // fallback to default string representation
		}
		result = append(result, keyStr)
	}

	return result
}

// Len returns number of map entries
func (m *mapAccessor) Len() int {
	if m.IsNil() {
		return 0
	}
	return m.value.Len()
}
