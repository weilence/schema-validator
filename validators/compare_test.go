package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/schema"
)

func TestCompareValidators(t *testing.T) {
	r := NewRegistry()
	registerCompare(r)

	tests := []struct {
		name     string
		ruleName string
		value    any
		params   []any
		wantErr  bool
	}{
		// eq
		{"eq valid", "eq", "hello", []any{"hello"}, false},
		{"eq invalid", "eq", "hello", []any{"world"}, true},
		// eq_ignore_case
		{"eq_ignore_case valid", "eq_ignore_case", "Hello", []any{"hello"}, false},
		{"eq_ignore_case invalid", "eq_ignore_case", "Hello", []any{"world"}, true},
		// gt
		{"gt valid", "gt", 10, []any{"5"}, false},
		{"gt invalid", "gt", 5, []any{"10"}, true},
		// gte
		{"gte valid", "gte", 10, []any{"10"}, false},
		{"gte invalid", "gte", 5, []any{"10"}, true},
		// lt
		{"lt valid", "lt", 5, []any{"10"}, false},
		{"lt invalid", "lt", 10, []any{"5"}, true},
		// lte
		{"lte valid", "lte", 10, []any{"10"}, false},
		{"lte invalid", "lte", 15, []any{"10"}, true},
		// ne
		{"ne valid", "ne", "hello", []any{"world"}, false},
		{"ne invalid", "ne", "hello", []any{"hello"}, true},
		// ne_ignore_case
		{"ne_ignore_case valid", "ne_ignore_case", "Hello", []any{"world"}, false},
		{"ne_ignore_case invalid", "ne_ignore_case", "Hello", []any{"hello"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := schema.NewObject().
				AddField("test", schema.NewField().AddValidator(r.NewValidator(tt.ruleName, tt.params...)))
			ctx := schema.NewContext(s, data.New(map[string]any{"test": tt.value}))
			err := s.Validate(ctx)
			assert.NoError(t, err)
			assert.Equal(t, ctx.Errors().HasErrorCode(tt.ruleName), tt.wantErr)
		})
	}
}
