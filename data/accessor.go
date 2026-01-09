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

	originalRv := rv

	derefRv := rv
	for derefRv.Kind() == reflect.Pointer {
		if derefRv.IsNil() {
			return NewValueAccessor(originalRv)
		}
		derefRv = derefRv.Elem()
	}

	switch derefRv.Kind() {
	case reflect.Slice, reflect.Array:
		return NewArrayAccessor(originalRv)
	case reflect.Map:
		return NewMapAccessor(originalRv)
	case reflect.Struct:
		return NewStructAccessor(originalRv)
	default:
		return NewValueAccessor(originalRv)
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
