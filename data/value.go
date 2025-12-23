package data

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type Value struct {
	rval reflect.Value
}

func (p *Value) Raw() any {
	return p.rval.Interface()
}

func (p *Value) GetValue(path string) (*Value, error) {
	if path != "" {
		return nil, errors.New("cannot traverse further on primitive value")
	}

	return p, nil
}

// String returns string representation
func (p *Value) String() string {
	val := p.rval
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
func (p *Value) Int() (int64, error) {
	val := p.rval
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
func (p *Value) Float() (float64, error) {
	val := p.rval
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
func (p *Value) Bool() (bool, error) {
	val := p.rval
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

func (p *Value) Any() any {
	val := p.rval
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}
	return val.Interface()
}
