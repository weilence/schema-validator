package data

import (
	"errors"
	"reflect"

	"github.com/spf13/cast"
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

func (p *Value) Kind() reflect.Kind {
	if p.rval.Kind() == reflect.Interface {
		return p.rval.Elem().Kind()
	}

	return p.rval.Kind()
}

func (p *Value) Len() int {
	if p.rval.Kind() == reflect.Interface {
		return p.rval.Elem().Len()
	}

	return p.rval.Len()
}

func (p *Value) IsInt() bool {
	kind := p.Kind()
	return kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64
}

func (p *Value) IsUint() bool {
	kind := p.Kind()
	return kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64
}

func (p *Value) IsFloat() bool {
	kind := p.Kind()
	return kind == reflect.Float32 || kind == reflect.Float64
}

// String returns string representation
func (p *Value) String() string {
	v, err := cast.ToStringE(p.rval.Interface())
	if err != nil {
		return p.rval.String()
	}

	return v
}

// Int returns int64 value
func (p *Value) Int() (int64, error) {
	return cast.ToInt64E(p.rval.Interface())
}

// Float returns float64 value
func (p *Value) Float() (float64, error) {
	return cast.ToFloat64E(p.rval.Interface())
}

// Bool returns bool value
func (p *Value) Bool() (bool, error) {
	return cast.ToBoolE(p.rval.Interface())
}

func (p *Value) Any() any {
	val := p.rval
	for val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}
	return val.Interface()
}
