package data

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue_GetValue(t *testing.T) {
	v := &Value{rval: reflect.ValueOf(123)}
	_, err := v.GetValue("x")
	assert.Error(t, err)

	ret, err := v.GetValue("")
	assert.NoError(t, err)
	assert.Equal(t, v, ret)
}

func TestValue_String(t *testing.T) {
	tests := []struct {
		name string
		val  reflect.Value
		want string
	}{
		{"string", reflect.ValueOf("hello"), "hello"},
		{"int", reflect.ValueOf(42), "42"},
		{"uint", reflect.ValueOf(uint(7)), "7"},
		{"float", reflect.ValueOf(1.5), "1.5"},
		{"bool", reflect.ValueOf(true), "true"},
		{"nilptr", reflect.ValueOf((*int)(nil)), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Value{rval: tt.val}
			assert.Equal(t, tt.want, v.String())
		})
	}
}

func TestValue_Int(t *testing.T) {
	tests := []struct {
		name    string
		val     reflect.Value
		want    int64
		wantErr bool
	}{
		{"int", reflect.ValueOf(int64(42)), 42, false},
		{"uint", reflect.ValueOf(uint(10)), 10, false},
		{"float", reflect.ValueOf(float32(3.0)), 3, false},
		{"string", reflect.ValueOf("123"), 123, false},
		{"badstring", reflect.ValueOf("x"), 0, true},
		{"nilptr", reflect.ValueOf((*int)(nil)), 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Value{rval: tt.val}
			got, err := v.Int()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestValue_Float(t *testing.T) {
	tests := []struct {
		name    string
		val     reflect.Value
		want    float64
		wantErr bool
	}{
		{"float64", reflect.ValueOf(1.23), 1.23, false},
		{"int", reflect.ValueOf(int(2)), 2.0, false},
		{"uint", reflect.ValueOf(uint(3)), 3.0, false},
		{"string", reflect.ValueOf("4.56"), 4.56, false},
		{"badstring", reflect.ValueOf("nope"), 0, true},
		{"nilptr", reflect.ValueOf((*float64)(nil)), 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Value{rval: tt.val}
			got, err := v.Float()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestValue_Bool(t *testing.T) {
	tests := []struct {
		name    string
		val     reflect.Value
		want    bool
		wantErr bool
	}{
		{"true", reflect.ValueOf(true), true, false},
		{"stringTrue", reflect.ValueOf("true"), true, false},
		{"stringFalse", reflect.ValueOf("false"), false, false},
		{"badstring", reflect.ValueOf("yes"), false, true},
		{"nilptr", reflect.ValueOf((*bool)(nil)), false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Value{rval: tt.val}
			got, err := v.Bool()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestValue_Any(t *testing.T) {
	vnil := &Value{rval: reflect.ValueOf((*int)(nil))}
	assert.Nil(t, vnil.Any())

	v := &Value{rval: reflect.ValueOf(7)}
	assert.Equal(t, 7, v.Any())
}
