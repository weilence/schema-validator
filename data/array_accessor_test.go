package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayAccessor_GetValue_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		arr     any
		path    string
		want    string
		wantErr bool
	}{
		{"string index", []string{"zero", "one", "two"}, "[2]", "two", false},
		{"int index", []int{7, 8, 9}, "[0]", "7", false},
		{"pointer to struct elem", []*struct{ V int }{{V: 1}, {V: 2}}, "[1]", "<struct { V int } Value>", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := New(tt.arr)
			v, err := acc.GetValue(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, v.String())
		})
	}
}

func TestArrayAccessor_Errors_TableDriven(t *testing.T) {
	acc := New([]int{1, 2})

	tests := []struct {
		name    string
		path    string
		wantErr string
	}{
		{"plain invalid", "x", "invalid array index: x"},
		{"scan invalid", "[x]", "invalid array index in scan: [x]"},
		{"oob", "[5]", "index 5 out of bounds"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := acc.GetValue(tt.path)
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}

func TestArrayAccessor_GetIndex_Len_Iterate(t *testing.T) {
	arr := []int{10, 20}
	aa := New(arr).(*ArrayAccessor)

	assert.Equal(t, 2, aa.Len())

	elemAcc, err := aa.GetIndex(1)
	assert.NoError(t, err)
	assert.IsType(t, &Value{}, elemAcc)

	p := elemAcc.(*Value)
	v, err := p.Int()
	assert.NoError(t, err)
	assert.Equal(t, int64(20), v)

	sum := 0
	err = aa.Iterate(func(i int, e Accessor) error {
		pv, err := e.GetValue("")
		if err != nil {
			return err
		}

		vv, err := pv.Int()
		if err != nil {
			return err
		}
		sum += int(vv)

		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 30, sum)

	_, err = aa.GetIndex(5)
	assert.Error(t, err)
}
