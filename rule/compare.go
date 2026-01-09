package rule

import (
	"strings"

	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/schema"
)

func registerCompare(r *Registry) {
	r.Register("eq", func(ctx *schema.Context, other string) error {
		currentValue := ctx.Value()
		otherValue := data.NewValue(other)
		ok, err := compareValue(Equal, currentValue, otherValue)
		if err != nil {
			return err
		}
		if !ok {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("eq_ignore_case", func(ctx *schema.Context, other string) error {
		currentStr := strings.ToLower(ctx.Value().String())
		otherStr := strings.ToLower(other)
		if currentStr != otherStr {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("gt", func(ctx *schema.Context, other string) error {
		currentValue := ctx.Value()
		otherValue := data.NewValue(other)
		ok, err := compareValue(GreaterThan, currentValue, otherValue)
		if err != nil {
			return err
		}
		if !ok {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("gte", func(ctx *schema.Context, other string) error {
		currentValue := ctx.Value()
		otherValue := data.NewValue(other)
		ok, err := compareValue(GreaterThanOrEqual, currentValue, otherValue)
		if err != nil {
			return err
		}
		if !ok {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("lt", func(ctx *schema.Context, other string) error {
		currentValue := ctx.Value()
		otherValue := data.NewValue(other)
		ok, err := compareValue(LessThan, currentValue, otherValue)
		if err != nil {
			return err
		}
		if !ok {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("lte", func(ctx *schema.Context, other string) error {
		currentValue := ctx.Value()
		otherValue := data.NewValue(other)
		ok, err := compareValue(LessThanOrEqual, currentValue, otherValue)
		if err != nil {
			return err
		}
		if !ok {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("ne", func(ctx *schema.Context, other string) error {
		currentValue := ctx.Value()
		otherValue := data.NewValue(other)
		ok, err := compareValue(NotEqual, currentValue, otherValue)
		if err != nil {
			return err
		}
		if !ok {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("ne_ignore_case", func(ctx *schema.Context, other string) error {
		currentStr := strings.ToLower(ctx.Value().String())
		otherStr := strings.ToLower(other)
		if currentStr == otherStr {
			return schema.ErrCheckFailed
		}
		return nil
	})
}
