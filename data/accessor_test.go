package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccessor_Combinations_TableDriven(t *testing.T) {
	type T1 struct{ Items []map[string]string }
	type T2 struct{ M map[string]int }
	type T3 struct{ Arr []string }

	tests := []struct {
		name    string
		root    any
		path    string
		want    string
		wantErr bool
	}{
		{"array->map", []map[string]int{{"a": 1}, {"b": 2}}, "[1].b", "2", false},
		{"map->array", map[string][]int{"arr": {5, 6}}, "arr.[0]", "5", false},
		{"struct->array->map", T1{Items: []map[string]string{{"x": "v0"}}}, "Items.[0].x", "v0", false},
		{"array->struct->map", []T2{{M: map[string]int{"a": 10}}}, "[0].M.a", "10", false},
		{"map->struct->array", map[string]T3{"k": {Arr: []string{"p", "q"}}}, "k.Arr.[1]", "q", false},
		// missing key/field
		{"missing map key", map[string]int{"a": 1}, "b", "", true},
		{"missing struct field", struct{ A int }{A: 1}, "Nope", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := New(tt.root)
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
