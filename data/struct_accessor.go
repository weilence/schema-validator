package data

import (
	"fmt"
	"reflect"
)

type structAccessor struct {
	value reflect.Value

	embedValues []*structAccessor
}

func NewStructAccessor(v reflect.Value) *structAccessor {
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	var embedValues []*structAccessor
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if field.Anonymous {
			embedField := v.Field(i)
			if embedField.Kind() == reflect.Pointer && !embedField.IsNil() {
				embedField = embedField.Elem()
			}

			if embedField.Kind() == reflect.Struct {
				embedValues = append(embedValues, NewStructAccessor(embedField))
			}
		}
	}

	return &structAccessor{
		value:       v,
		embedValues: embedValues,
	}
}

func (s *structAccessor) Raw() any {
	if s.value.CanInterface() {
		return s.value.Interface()
	}

	return nil
}

// TODO: resolve conflict with embedded fields in same name
func (s *structAccessor) GetValue(path string) (*Value, error) {
	if path == "" {
		return &Value{rval: s.value}, nil
	}

	fieldName, nextPath := cutPath(path)

	fieldAcc, err := s.GetField(fieldName)
	if err != nil {
		return nil, err
	}

	return fieldAcc.GetValue(nextPath)
}

// GetField returns field by name (supports embedded fields)
func (s *structAccessor) GetField(name string) (Accessor, error) {
	v := s.value.FieldByName(name)
	if v != (reflect.Value{}) {
		return NewAccessor(v), nil
	}

	// Check embedded structs
	for _, embed := range s.embedValues {
		if fieldAcc, err := embed.GetField(name); err == nil {
			return fieldAcc, nil
		}
	}

	return nil, fmt.Errorf("field %s not found", name)
}

func (s *structAccessor) Accessors() []ObjectAccessor {
	accessors := []ObjectAccessor{s}
	for _, embed := range s.embedValues {
		accessors = append(accessors, embed.Accessors()...)
	}

	return accessors
}
