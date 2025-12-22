package data

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type primitiveAccessor struct {
	value *Value
}

// Kind returns KindPrimitive
func (p *primitiveAccessor) Kind() DataKind {
	return KindPrimitive
}

// IsNil checks if the value is nil
func (p *primitiveAccessor) IsNil() bool {
	return p.value.IsNil()
}

// AsField returns itself
func (p *primitiveAccessor) AsField() (FieldAccessor, error) {
	return p, nil
}

// AsObject returns error
func (p *primitiveAccessor) AsObject() (ObjectAccessor, error) {
	return nil, errors.New("primitive is not an object")
}

// AsArray returns error
func (p *primitiveAccessor) AsArray() (ArrayAccessor, error) {
	return nil, errors.New("primitive is not an array")
}

// Value returns the underlying value
func (p *primitiveAccessor) Value() interface{} {
	if p.value.IsNil() {
		return nil
	}
	return p.value.raw
}

// String returns string representation
func (p *primitiveAccessor) String() string {
	if p.value.IsNil() {
		return ""
	}

	val := p.value.rval
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return ""
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.String:
		return val.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(val.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(val.Bool())
	default:
		return fmt.Sprintf("%v", val.Interface())
	}
}

// Int returns int64 value
func (p *primitiveAccessor) Int() (int64, error) {
	if p.value.IsNil() {
		return 0, errors.New("nil value")
	}

	val := p.value.rval
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return 0, errors.New("nil pointer")
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(val.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int64(val.Float()), nil
	case reflect.String:
		return strconv.ParseInt(val.String(), 10, 64)
	default:
		return 0, fmt.Errorf("cannot convert %v to int", val.Kind())
	}
}

// Float returns float64 value
func (p *primitiveAccessor) Float() (float64, error) {
	if p.value.IsNil() {
		return 0, errors.New("nil value")
	}

	val := p.value.rval
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return 0, errors.New("nil pointer")
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Float32, reflect.Float64:
		return val.Float(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(val.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(val.Uint()), nil
	case reflect.String:
		return strconv.ParseFloat(val.String(), 64)
	default:
		return 0, fmt.Errorf("cannot convert %v to float", val.Kind())
	}
}

// Bool returns bool value
func (p *primitiveAccessor) Bool() (bool, error) {
	if p.value.IsNil() {
		return false, errors.New("nil value")
	}

	val := p.value.rval
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return false, errors.New("nil pointer")
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Bool:
		return val.Bool(), nil
	case reflect.String:
		return strconv.ParseBool(val.String())
	default:
		return false, fmt.Errorf("cannot convert %v to bool", val.Kind())
	}
}
