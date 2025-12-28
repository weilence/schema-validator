package validators

import (
	"fmt"
	"slices"

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
			return schema.ErrCheckFailed
		}

		return nil
	}
}

func registerOther(r *Registry) {
	r.Register("min", compareValidator(GreaterThanOrEqual))

	r.Register("max", compareValidator(LessThanOrEqual))

	r.Register("oneof", func(ctx *schema.Context, params []string) error {
		val := ctx.Value().String()
		if val == "" {
			return nil
		}

		if slices.Contains(params, val) {
			return nil
		}

		return schema.ErrCheckFailed
	})

	requiredFn := func(ctx *schema.Context) error {
		if ctx.Value().IsNilOrZero() {
			return schema.ErrCheckFailed
		}

		return nil
	}

	r.Register("required", requiredFn)

	r.Register("required_if", func(ctx *schema.Context, fieldName string, expectedValue any) error {
		otherValue, err := ctx.Parent().GetValue(fieldName)
		if err != nil {
			return fmt.Errorf("failed to get field '%s': %v", fieldName, err)
		}

		ok, err := compareValue(Equal, otherValue, data.NewValue(expectedValue))
		if err != nil {
			return err
		}

		if ok {
			return requiredFn(ctx)
		}

		return nil
	})

	r.Register("omitempty", func(ctx *schema.Context) error {
		if ctx.Value().IsNilOrZero() {
			ctx.SkipRest()
		}

		return nil
	})
}
