package validators

import (
	"github.com/weilence/schema-validator/schema"
)

func registerField(r *Registry) {
	r.Register("eqfield", func(ctx *schema.Context, fieldName string) error {
		currentValue, err := ctx.Value()
		if err != nil {
			return err
		}
		parentObj := ctx.Parent()
		otherValue, err := parentObj.GetValue(fieldName)
		if err != nil {
			return err
		}
		if currentValue.String() != otherValue.String() {
			return schema.NewValidationError(ctx.Path(), "eqfield", map[string]any{"field": fieldName})
		}
		return nil
	})
	r.Register("nefield", func(ctx *schema.Context, fieldName string) error {
		otherValue, err := ctx.Parent().GetValue(fieldName)
		if err != nil {
			return err
		}
		currentValue, err := ctx.Value()
		if err != nil {
			return err
		}
		if currentValue.String() == otherValue.String() {
			return schema.NewValidationError(ctx.Path(), "nefield", map[string]any{"field": fieldName})
		}
		return nil
	})
	r.Register("gtfield", func(ctx *schema.Context, fieldName string) error {
		otherValue, err := ctx.Parent().GetValue(fieldName)
		if err != nil {
			return err
		}
		currentValue, err := ctx.Value()
		if err != nil {
			return err
		}
		val, err1 := currentValue.Int()
		otherVal, err2 := otherValue.Int()
		if err1 == nil && err2 == nil {
			if val <= otherVal {
				return schema.NewValidationError(ctx.Path(), "gtfield", map[string]any{"field": fieldName})
			}
			return nil
		}
		fval, err1 := currentValue.Float()
		fotherVal, err2 := otherValue.Float()
		if err1 == nil && err2 == nil {
			if fval <= fotherVal {
				return schema.NewValidationError(ctx.Path(), "gtfield", map[string]any{"field": fieldName})
			}
			return nil
		}
		if currentValue.String() <= otherValue.String() {
			return schema.NewValidationError(ctx.Path(), "gtfield", map[string]any{"field": fieldName})
		}
		return nil
	})
	r.Register("ltfield", func(ctx *schema.Context, fieldName string) error {
		otherValue, err := ctx.Parent().GetValue(fieldName)
		if err != nil {
			return err
		}
		currentValue, err := ctx.Value()
		if err != nil {
			return err
		}
		val, err1 := currentValue.Int()
		otherVal, err2 := otherValue.Int()
		if err1 == nil && err2 == nil {
			if val >= otherVal {
				return schema.NewValidationError(ctx.Path(), "ltfield", map[string]any{"field": fieldName})
			}
			return nil
		}
		fval, err1 := currentValue.Float()
		fotherVal, err2 := otherValue.Float()
		if err1 == nil && err2 == nil {
			if fval >= fotherVal {
				return schema.NewValidationError(ctx.Path(), "ltfield", map[string]any{"field": fieldName})
			}
			return nil
		}
		if currentValue.String() >= otherValue.String() {
			return schema.NewValidationError(ctx.Path(), "ltfield", map[string]any{"field": fieldName})
		}
		return nil
	})
}
