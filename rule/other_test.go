package rule

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/schema"
)

func TestOtherValidators(t *testing.T) {
	r := NewRegistry()
	registerOther(r)

	tests := []struct {
		name     string
		ruleName string
		value    any
		params   []any
		wantErr  bool
	}{
		// dir
		{"dir valid", "dir", "/tmp", nil, false},
		{"dir invalid", "dir", "/nonexistent", nil, true},
		// dirpath
		{"dirpath valid", "dirpath", "/tmp/test", nil, false},
		{"dirpath invalid", "dirpath", "invalid", nil, true},
		// file
		{"file valid", "file", "/etc/hosts", nil, false},
		{"file invalid", "file", "/nonexistent", nil, true},
		// filepath
		{"filepath valid", "filepath", "/tmp/test.txt", nil, false},
		{"filepath invalid", "filepath", "invalid", nil, true},
		// image
		{"image valid", "image", "test.jpg", nil, false},
		{"image invalid", "image", "test.txt", nil, true},
		// isdefault
		{"isdefault valid", "isdefault", "", nil, false},
		{"isdefault invalid", "isdefault", "value", nil, true},
		// len
		{"len valid", "len", "hello", []any{5}, false},
		{"len invalid", "len", "hello", []any{3}, true},
		// max
		{"max valid", "max", 5, []any{10}, false},
		{"max invalid", "max", 15, []any{10}, true},
		// min
		{"min valid", "min", 10, []any{5}, false},
		{"min invalid", "min", 3, []any{5}, true},
		// oneof
		{"oneof valid", "oneof", "a", []any{[]string{"a", "b", "c"}}, false},
		{"oneof invalid", "oneof", "d", []any{[]string{"a", "b", "c"}}, true},
		// required
		{"required valid", "required", "value", nil, false},
		{"required valid 2", "required", ptr(0), nil, false},
		{"required invalid", "required", "", nil, true},
		{"required invalid 2", "required", 0, nil, true},
		// unique (placeholder)
		{"unique valid", "unique", "value", nil, false},
		{"unique invalid", "unique", "value", nil, false}, // always pass for now
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

// TestMinValidatorWithPointerTypes tests min validator with pointer types
func TestMinValidatorWithPointerTypes(t *testing.T) {
	r := NewRegistry()
	registerOther(r)

	tests := []struct {
		name     string
		ruleName string
		value    any
		params   []any
		wantErr  bool
	}{
		{"min valid *int", "min", ptr(10), []any{5}, false},
		{"min invalid *int", "min", ptr(3), []any{5}, true},
		{"min equal *int", "min", ptr(5), []any{5}, false},

		{"min valid *float64", "min", ptr(float64(10.5)), []any{5}, false},
		{"min invalid *float64", "min", ptr(float64(3.5)), []any{5}, true},

		{"min nil *int", "min", (*int)(nil), []any{5}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := schema.NewObject().
				AddField("test", schema.NewField().AddValidator(r.NewValidator(tt.ruleName, tt.params...)))
			ctx := schema.NewContext(s, data.New(map[string]any{"test": tt.value}))
			err := s.Validate(ctx)
			assert.NoError(t, err)
			assert.Equal(t, ctx.Errors().HasErrorCode(tt.ruleName), tt.wantErr,
				fmt.Sprintf("test case '%s' failed: value=%v, params=%v", tt.name, tt.value, tt.params))
		})
	}
}

// TestMaxValidatorWithPointerTypes tests max validator with pointer types
func TestMaxValidatorWithPointerTypes(t *testing.T) {
	r := NewRegistry()
	registerOther(r)

	tests := []struct {
		name     string
		ruleName string
		value    any
		params   []any
		wantErr  bool
	}{
		{"max valid *int", "max", ptr(5), []any{10}, false},
		{"max invalid *int", "max", ptr(15), []any{10}, true},
		{"max equal *int", "max", ptr(10), []any{10}, false},

		{"max valid *float64", "max", ptr(float64(5.5)), []any{10}, false},
		{"max invalid *float64", "max", ptr(float64(15.5)), []any{10}, true},

		{"max nil *int", "max", (*int)(nil), []any{10}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := schema.NewObject().
				AddField("test", schema.NewField().AddValidator(r.NewValidator(tt.ruleName, tt.params...)))
			ctx := schema.NewContext(s, data.New(map[string]any{"test": tt.value}))
			err := s.Validate(ctx)
			assert.NoError(t, err)
			assert.Equal(t, ctx.Errors().HasErrorCode(tt.ruleName), tt.wantErr,
				fmt.Sprintf("test case '%s' failed: value=%v, params=%v", tt.name, tt.value, tt.params))
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
