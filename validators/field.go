package validators

import (
	"github.com/weilence/schema-validator/schema"
)

func compareFieldValidator(ct compareType) func(*schema.Context, string) error {
	return func(ctx *schema.Context, fieldName string) error {
		currentValue := ctx.Value()
		otherValue, err := ctx.Parent().GetValue(fieldName)
		if err != nil {
			return err
		}

		ok, err := compareValue(ct, currentValue, otherValue)
		if err != nil {
			return err
		}

		if !ok {
			return schema.NewValidationError(ctx.Path(), ct.String()+"field", map[string]any{"field": fieldName})
		}

		return nil
	}
}

func registerField(r *Registry) {
	r.Register("eqfield", compareFieldValidator(Equal))
	r.Register("nefield", compareFieldValidator(NotEqual))
	r.Register("gtfield", compareFieldValidator(GreaterThan))
	r.Register("ltfield", compareFieldValidator(LessThan))
	r.Register("gtefield", compareFieldValidator(GreaterThanOrEqual))
	r.Register("ltefield", compareFieldValidator(LessThanOrEqual))
}
