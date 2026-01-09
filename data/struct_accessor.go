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
	derefV := v
	for derefV.Kind() == reflect.Pointer {
		derefV = derefV.Elem()
	}

	var embedValues []*structAccessor
	for i := 0; i < derefV.NumField(); i++ {
		field := derefV.Type().Field(i)
		if field.Anonymous {
			embedField := derefV.Field(i)
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

func (s *structAccessor) deref() reflect.Value {
	v := s.value
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	return v
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
		return NewValueAccessor(s.value), nil
	}

	fieldName, nextPath := cutPath(path)

	fieldAcc, err := s.GetField(fieldName)
	if err != nil {
		return nil, err
	}

	return fieldAcc.GetValue(nextPath)
}

func (s *structAccessor) GetField(name string) (Accessor, error) {
	v := s.deref()
	field := v.FieldByName(name)
	if field != (reflect.Value{}) {
		return NewAccessor(field), nil
	}

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
