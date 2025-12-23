package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructAccessor_TableDriven(t *testing.T) {
	type Inner struct{ Age int }
	type Outer struct{ Inner }
	type OuterPtr struct{ *Inner }

	tests := []struct {
		name    string
		data    any
		path    string
		want    string
		wantErr bool
	}{
		{"simple field ptr", &struct{ Name string }{Name: "alice"}, "Name", "alice", false},
		{"embedded value", Outer{Inner: Inner{Age: 30}}, "Age", "30", false},
		{"embedded ptr", OuterPtr{Inner: &Inner{Age: 42}}, "Age", "42", false},
		{"nested field", struct{ Addr struct{ City string } }{Addr: struct{ City string }{City: "Beijing"}}, "Addr.City", "Beijing", false},
		{"not exist", struct{ A int }{A: 1}, "Nope", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := New(tt.data)
			v, err := acc.GetValue(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tt.want, v.String())
		})
	}
}
