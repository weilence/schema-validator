package validators

import (
	"testing"

	"github.com/weilence/schema-validator/schema"
)

// Test multi-parameter validator factory
func TestMultiParameterValidator(t *testing.T) {
	// Register a custom validator that takes multiple parameters
	registry := NewRegistry()

	// Example: between=10,20 validator
	registry.Register("between", func(ctx *schema.Context, params []string) error {
		if len(params) < 2 {
			return nil
		}

		min := parseIntHelper(params[0])
		max := parseIntHelper(params[1])

		field, err := ctx.Value()
		if err != nil {
			return nil
		}

		val, err := field.Int()
		if err != nil {
			return nil
		}

		if val < int64(min) || val > int64(max) {
			return schema.NewValidationError(ctx.Path(), "between", map[string]any{
				"min":    min,
				"max":    max,
				"actual": val,
			})
		}
		return nil
	})
}

// Test multi-parameter with three params
func TestThreeParameterValidator(t *testing.T) {
	registry := NewRegistry()

	// Example: enum=option1:option2:option3 validator
	registry.Register("enum", func(ctx *schema.Context, params []string) error {
		if len(params) == 0 {
			return nil
		}

		allowedValues := make(map[string]bool)
		for _, p := range params {
			allowedValues[p] = true
		}

		field, err := ctx.Value()
		if err != nil {
			return nil
		}

		val := field.String()
		if !allowedValues[val] {
			return schema.NewValidationError(ctx.Path(), "enum", map[string]any{
				"allowed": params,
				"actual":  val,
			})
		}
		return nil
	})

}

// Helper function
func parseIntHelper(s string) int {
	var result int
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			result = result*10 + int(s[i]-'0')
		}
	}
	return result
}
