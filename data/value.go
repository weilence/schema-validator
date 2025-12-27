package data

import (
	"errors"
	"reflect"

	"github.com/spf13/cast"
)

type Value struct {
	rval reflect.Value
}

func NewValue(v any) *Value {
	rv := reflect.ValueOf(v)
	return &Value{rval: rv}
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
	return p.rval.CanInt()
}

func (p *Value) IsUint() bool {
	return p.rval.CanUint()
}

func (p *Value) IsFloat() bool {
	return p.rval.CanFloat()
}

func (p *Value) IsSliceOrArray() bool {
	kind := p.Kind()
	return kind == reflect.Slice || kind == reflect.Array
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
