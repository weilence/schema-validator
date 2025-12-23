package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapAccessor_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		m       any
		path    string
		want    string
		wantErr bool
	}{
		{"simple", map[string]string{"k": "v"}, "k", "v", false},
		{"nested concrete", map[string]map[string]int{"inner": {"n": 3}}, "inner.n", "3", false},
		{"missing", map[string]string{"a": "b"}, "nope", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := New(tt.m)
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
