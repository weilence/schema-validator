package validators

import (
	"testing"

	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/schema"
)

// Test multi-parameter validator factory
func TestMultiParameterValidator(t *testing.T) {
	// Register a custom validator that takes multiple parameters
	// Example: between=10,20 validator
	Register("between", func(ctx *schema.Context, min any, max any) error {
		field := ctx.Value()

		minValue := data.NewValue(min)
		maxValue := data.NewValue(max)

		ok1, err := compareValue(GreaterThanOrEqual, field, minValue)
		if err != nil {
			return err
		}

		ok2, err := compareValue(LessThanOrEqual, field, maxValue)
		if err != nil {
			return err
		}

		if !ok1 || !ok2 {
			return schema.NewValidationError(ctx.Path(), "between", map[string]any{
				"min":    min,
				"max":    max,
				"actual": field.Raw(),
			})
		}
		return nil
	})

	// Create a context with a field value

	s := schema.NewFieldSchema().AddValidator(NewValidator("between", 10, 20))
	ctx := schema.NewContext(s, data.NewValue(15))
	// Validate the context
	err := s.Validate(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Test multi-parameter with three params
func TestThreeParameterValidator(t *testing.T) {
	// Example: enum=option1,option2,option3 validator
	Register("enum", func(ctx *schema.Context, params []string) error {
		if len(params) == 0 {
			return nil
		}

		allowedValues := make(map[string]bool)
		for _, p := range params {
			allowedValues[p] = true
		}

		field := ctx.Value()
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
