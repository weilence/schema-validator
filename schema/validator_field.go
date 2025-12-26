package schema

func init() {
	Register("eqfield", func(ctx *Context, params []string) error {
		if len(params) == 0 {
			return nil
		}
		fieldName := params[0]
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
			return NewValidationError(ctx.Path(), "eqfield", map[string]any{"field": fieldName})
		}
		return nil
	})
	Register("nefield", func(ctx *Context, params []string) error {
		if len(params) == 0 {
			return nil
		}
		fieldName := params[0]
		otherValue, err := ctx.Parent().GetValue(fieldName)
		if err != nil {
			return err
		}
		currentValue, err := ctx.Value()
		if err != nil {
			return err
		}
		if currentValue.String() == otherValue.String() {
			return NewValidationError(ctx.Path(), "nefield", map[string]any{"field": fieldName})
		}
		return nil
	})
	Register("gtfield", func(ctx *Context, params []string) error {
		if len(params) == 0 {
			return nil
		}
		fieldName := params[0]
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
				return NewValidationError(ctx.Path(), "gtfield", map[string]any{"field": fieldName})
			}
			return nil
		}
		fval, err1 := currentValue.Float()
		fotherVal, err2 := otherValue.Float()
		if err1 == nil && err2 == nil {
			if fval <= fotherVal {
				return NewValidationError(ctx.Path(), "gtfield", map[string]any{"field": fieldName})
			}
			return nil
		}
		if currentValue.String() <= otherValue.String() {
			return NewValidationError(ctx.Path(), "gtfield", map[string]any{"field": fieldName})
		}
		return nil
	})
	Register("ltfield", func(ctx *Context, params []string) error {
		if len(params) == 0 {
			return nil
		}
		fieldName := params[0]
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
				return NewValidationError(ctx.Path(), "ltfield", map[string]any{"field": fieldName})
			}
			return nil
		}
		fval, err1 := currentValue.Float()
		fotherVal, err2 := otherValue.Float()
		if err1 == nil && err2 == nil {
			if fval >= fotherVal {
				return NewValidationError(ctx.Path(), "ltfield", map[string]any{"field": fieldName})
			}
			return nil
		}
		if currentValue.String() >= otherValue.String() {
			return NewValidationError(ctx.Path(), "ltfield", map[string]any{"field": fieldName})
		}
		return nil
	})
}
