package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/schema"
)

func TestStringValidators(t *testing.T) {
	r := NewRegistry()
	registerString(r)

	tests := []struct {
		name     string
		ruleName string
		value    string
		params   []any
		wantErr  bool
	}{
		// alpha
		{"alpha valid", "alpha", "abc", nil, false},
		{"alpha invalid", "alpha", "abc123", nil, true},
		// alphaspace
		{"alphaspace valid", "alphaspace", "hello world", nil, false},
		{"alphaspace invalid", "alphaspace", "hello123", nil, true},
		// alphanum
		{"alphanum valid", "alphanum", "abc123", nil, false},
		{"alphanum invalid", "alphanum", "abc-123", nil, true},
		// alphanumspace
		{"alphanumspace valid", "alphanumspace", "abc 123", nil, false},
		{"alphanumspace invalid", "alphanumspace", "abc-123", nil, true},
		// alphanumunicode
		{"alphanumunicode valid", "alphanumunicode", "abc123", nil, false},
		{"alphanumunicode invalid", "alphanumunicode", "abc-123", nil, true},
		// alphaunicode
		{"alphaunicode valid", "alphaunicode", "abc", nil, false},
		{"alphaunicode invalid", "alphaunicode", "abc123", nil, true},
		// ascii
		{"ascii valid", "ascii", "hello", nil, false},
		{"ascii invalid", "ascii", "héllo", nil, true},
		// boolean
		{"boolean valid", "boolean", "true", nil, false},
		{"boolean invalid", "boolean", "yes", nil, true},
		// contains
		{"contains valid", "contains", "hello world", []any{"world"}, false},
		{"contains invalid", "contains", "hello", []any{"world"}, true},
		// containsany
		{"containsany valid", "containsany", "hello", []any{"aeiou"}, false},
		{"containsany invalid", "containsany", "bcd", []any{"aeiou"}, true},
		// containsrune
		{"containsrune valid", "containsrune", "hello", []any{"e"}, false},
		{"containsrune invalid", "containsrune", "bcd", []any{"e"}, true},
		// endsnotwith
		{"endsnotwith valid", "endsnotwith", "hello", []any{"world"}, false},
		{"endsnotwith invalid", "endsnotwith", "hello world", []any{"world"}, true},
		// endswith
		{"endswith valid", "endswith", "hello world", []any{"world"}, false},
		{"endswith invalid", "endswith", "hello", []any{"world"}, true},
		// excludes
		{"excludes valid", "excludes", "hello", []any{"world"}, false},
		{"excludes invalid", "excludes", "hello world", []any{"world"}, true},
		// excludesall
		{"excludesall valid", "excludesall", "bcd", []any{"aeiou"}, false},
		{"excludesall invalid", "excludesall", "hello", []any{"aeiou"}, true},
		// excludesrune
		{"excludesrune valid", "excludesrune", "bcd", []any{"e"}, false},
		{"excludesrune invalid", "excludesrune", "hello", []any{"e"}, true},
		// lowercase
		{"lowercase valid", "lowercase", "hello", nil, false},
		{"lowercase invalid", "lowercase", "Hello", nil, true},
		// multibyte
		{"multibyte valid", "multibyte", "héllo", nil, false},
		{"multibyte invalid", "multibyte", "hello", nil, true},
		// number
		{"number valid", "number", "123", nil, false},
		{"number invalid", "number", "123.45", nil, true},
		// numeric
		{"numeric valid", "numeric", "123.45", nil, false},
		{"numeric invalid", "numeric", "abc", nil, true},
		// printascii
		{"printascii valid", "printascii", "hello", nil, false},
		{"printascii invalid", "printascii", "hello\n", nil, true},
		// startsnotwith
		{"startsnotwith valid", "startsnotwith", "world", []any{"hello"}, false},
		{"startsnotwith invalid", "startsnotwith", "hello world", []any{"hello"}, true},
		// startswith
		{"startswith valid", "startswith", "hello world", []any{"hello"}, false},
		{"startswith invalid", "startswith", "world", []any{"hello"}, true},
		// uppercase
		{"uppercase valid", "uppercase", "HELLO", nil, false},
		{"uppercase invalid", "uppercase", "Hello", nil, true},
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
