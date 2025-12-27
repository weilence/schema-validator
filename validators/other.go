package validators

import (
	"fmt"
	"slices"
	"strings"

	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/schema"
)
func compareValidator(ct compareType) func(*schema.Context, any) error {
	return func(ctx *schema.Context, value any) error {
		field := ctx.Value()
		otherValue := data.NewValue(value)
		ok, err := compareValue(ct, field, otherValue)
		if err != nil {
			return err
		}

		if !ok {
			return schema.NewValidationError(ctx.Path(), ct.String(), map[string]any{
				"expected": value,
				"actual":   field.Raw(),
			})
		}

		return nil
	}
}

func registerOther(r *Registry) {
	r.Register("min", compareValidator(GreaterThanOrEqual))

	r.Register("max", compareValidator(LessThanOrEqual))

	r.Register("oneof", func(ctx *schema.Context, params []string) error {
		field := ctx.Value()
		val := field.String()
		if val == "" {
			return nil
		}

		if slices.Contains(params, val) {
			return nil
		}

		return schema.NewValidationError(ctx.Path(), "oneof", map[string]any{
			"allowed": params,
			"actual":  val,
		})
	})

	r.Register("required", func(ctx *schema.Context) error {
		field := ctx.Value()
		str := field.String()
		if strings.TrimSpace(str) == "" {
			return schema.NewValidationError(ctx.Path(), "required", nil)
		}
		return nil
	})

	r.Register("required_if", func(ctx *schema.Context, params []string) error {
		if len(params) < 2 {
			return fmt.Errorf("required_if validator needs 2 parameters")
		}

		otherFieldName := params[0]
		expectedValue := params[1]

		otherValue, err := ctx.Parent().GetValue(otherFieldName)
		if err != nil {
			return fmt.Errorf("failed to get field '%s': %v", otherFieldName, err)
		}

		currentField := ctx.Value()
		if otherValue.String() == expectedValue {
			// Check if current field is non-empty
			if currentField.String() == "" {
				return schema.NewValidationError(ctx.Path(), "required_if", map[string]any{
					"field":         otherFieldName,
					"expectedValue": expectedValue,
				})
			}
		}

		return nil
	})
}
