package data

import (
	"errors"
	"reflect"
)

// Value wraps any Go value and provides unified access
type Value struct {
	raw  interface{}
	rval reflect.Value
	kind DataKind
}

// NewValue creates a new value wrapper
func NewValue(v interface{}) *Value {
	if v == nil {
		return &Value{
			raw:  nil,
			rval: reflect.Value{},
			kind: KindPrimitive,
		}
	}

	rval := reflect.ValueOf(v)
	return &Value{
		raw:  v,
		rval: rval,
		kind: detectKind(rval),
	}
}

// Raw returns the underlying raw value
func (v *Value) Raw() interface{} {
	return v.raw
}

func detectKind(v reflect.Value) DataKind {
	// Handle pointers by dereferencing
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return KindPrimitive
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		return KindObject
	case reflect.Map:
		return KindObject
	case reflect.Slice, reflect.Array:
		return KindArray
	default:
		return KindPrimitive
	}
}

// Kind returns the kind of data
func (v *Value) Kind() DataKind {
	return v.kind
}

// IsNil checks if the underlying data is nil
func (v *Value) IsNil() bool {
	if v.raw == nil {
		return true
	}
	if !v.rval.IsValid() {
		return true
	}

	// Check if the value itself is nil (for pointers, slices, maps, etc.)
	switch v.rval.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.rval.IsNil()
	}

	return false
}

// AsField returns field accessor (for primitives)
func (v *Value) AsField() (FieldAccessor, error) {
	if v.kind != KindPrimitive {
		return nil, errors.New("not a primitive field")
	}
	return &primitiveAccessor{value: v}, nil
}

// AsObject returns object accessor (for structs/maps)
func (v *Value) AsObject() (ObjectAccessor, error) {
	if v.kind != KindObject {
		return nil, errors.New("not an object")
	}

	rval := v.rval
	for rval.Kind() == reflect.Ptr {
		if rval.IsNil() {
			return nil, errors.New("nil pointer")
		}
		rval = rval.Elem()
	}

	if rval.Kind() == reflect.Struct {
		return newStructAccessor(rval), nil
	}
	return newMapAccessor(rval), nil
}

// AsArray returns array accessor (for slices/arrays)
func (v *Value) AsArray() (ArrayAccessor, error) {
	if v.kind != KindArray {
		return nil, errors.New("not an array")
	}
	return newSliceAccessor(v.rval), nil
}
