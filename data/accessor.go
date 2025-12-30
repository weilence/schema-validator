package data

import (
	"reflect"
	"strings"
)

// DataKind represents the kind of data
type DataKind int

const (
	KindPrimitive DataKind = iota
	KindArray
	KindObject
)

// Accessor provides unified interface for accessing different data types
type Accessor interface {
	GetField(name string) (Accessor, error)
	GetValue(path string) (*Value, error)
	Raw() any
}

func New(data any) Accessor {
	rv := reflect.ValueOf(data)
	return NewAccessor(rv)
}

func NewAccessor(rv reflect.Value) Accessor {
	if rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		return NewArrayAccessor(rv)
	case reflect.Map:
		return NewMapAccessor(rv)
	case reflect.Struct, reflect.Pointer:
		return NewStructAccessor(rv)
	default:
		return &Value{rval: rv}
	}
}

type ObjectAccessor interface {
	Accessor
	Accessors() []ObjectAccessor
}

func cutPath(path string) (string, string) {
	before, after, _ := strings.Cut(path, ".")
	return before, after
}
