package data

import (
	"errors"
	"reflect"

	"github.com/spf13/cast"
)

type Value struct {
	rval reflect.Value
}

func NewValueAccessor(rv reflect.Value) *Value {
	return &Value{rval: rv}
}

func NewValue(v any) *Value {
	rv := reflect.ValueOf(v)
	return NewValueAccessor(rv)
}

func (p *Value) GetField(name string) (Accessor, error) {
	if name != "" {
		return nil, errors.New("cannot get field from primitive value")
	}

	return p, nil
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
	return p.rval.Kind()
}

func (p *Value) Len() int {
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

func (p *Value) IsString() bool {
	return p.rval.Kind() == reflect.String
}

func (p *Value) IsSliceOrArray() bool {
	kind := p.Kind()
	return kind == reflect.Slice || kind == reflect.Array
}

// String returns string representation
func (p *Value) String() string {
	return cast.Must[string](cast.ToStringE(p.rval.Interface()))
}

// Int returns int64 value
func (p *Value) Int() int64 {
	return cast.Must[int64](cast.ToInt64E(p.rval.Interface()))
}

// Int returns int64 value
func (p *Value) IntE() (int64, error) {
	return cast.ToInt64E(p.rval.Interface())
}

// Float returns float64 value
func (p *Value) Float() float64 {
	return cast.Must[float64](cast.ToFloat64E(p.rval.Interface()))
}

// Bool returns bool value
func (p *Value) Bool() bool {
	return cast.Must[bool](cast.ToBoolE(p.rval.Interface()))
}

func (p *Value) Uint() uint64 {
	return cast.Must[uint64](cast.ToUint64E(p.rval.Interface()))
}

func (p *Value) IsNilOrZero() bool {
	v := p.rval

	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Map,
		reflect.Pointer,
		reflect.Slice:
		return v.IsNil()
	}

	return v.IsZero()
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
