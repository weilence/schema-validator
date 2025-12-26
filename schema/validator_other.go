package schema

import (
	"fmt"
	"slices"
	"strings"
)

func init() {
	Register("min", func(ctx *Context, params []string) error {
		if len(params) == 0 {
			return nil
		}
		minAny := parseIntOrString(params[0])
		field, err := ctx.Value()
		if err != nil {
			return nil
		}
		switch m := minAny.(type) {
		case int, int64:
			minVal := toInt64(m)
			val, err := field.Int()
			if err == nil && val < minVal {
				return NewValidationError(ctx.Path(), "min", map[string]any{"min": minVal, "actual": val})
			}
		case float64:
			val, err := field.Float()
			if err == nil && val < m {
				return NewValidationError(ctx.Path(), "min", map[string]any{"min": m, "actual": val})
			}
		case string:
			str := field.String()
			minLen := len(m)
			if len(str) < minLen {
				return NewValidationError(ctx.Path(), "min_length", map[string]any{"min": minLen, "actual": len(str)})
			}
		}
		return nil
	})

	Register("max", func(ctx *Context, params []string) error {
		if len(params) == 0 {
			return nil
		}
		maxAny := parseIntOrString(params[0])
		field, err := ctx.Value()
		if err != nil {
			return nil
		}
		switch m := maxAny.(type) {
		case int, int64:
			maxVal := toInt64(m)
			val, err := field.Int()
			if err == nil && val > maxVal {
				return NewValidationError(ctx.Path(), "max", map[string]any{"max": maxVal, "actual": val})
			}
		case float64:
			val, err := field.Float()
			if err == nil && val > m {
				return NewValidationError(ctx.Path(), "max", map[string]any{"max": m, "actual": val})
			}
		case string:
			str := field.String()
			maxLen := len(m)
			if len(str) > maxLen {
				return NewValidationError(ctx.Path(), "max_length", map[string]any{"max": maxLen, "actual": len(str)})
			}
		}
		return nil
	})

	Register("oneof", func(ctx *Context, params []string) error {
		field, err := ctx.GetValue("")
		if err != nil {
			return nil
		}

		val := field.String()
		if val == "" {
			return nil
		}

		if slices.Contains(params, val) {
			return nil
		}

		return NewValidationError(ctx.Path(), "oneof", map[string]any{
			"allowed": params,
			"actual":  val,
		})
	})

	Register("required", func(ctx *Context, params []string) error {
		field, err := ctx.Value()
		if err != nil {
			return nil
		}
		str := field.String()
		if strings.TrimSpace(str) == "" {
			return NewValidationError(ctx.Path(), "required", nil)
		}
		return nil
	})

	Register("required_if", func(ctx *Context, params []string) error {
		if len(params) < 2 {
			return fmt.Errorf("required_if validator needs 2 parameters")
		}

		otherFieldName := params[0]
		expectedValue := params[1]

		otherValue, err := ctx.Parent().GetValue(otherFieldName)
		if err != nil {
			return fmt.Errorf("failed to get field '%s': %v", otherFieldName, err)
		}

		currentField, err := ctx.Value()
		if err != nil {
			return fmt.Errorf("failed to get current field value: %v", err)
		}

		if otherValue.String() == expectedValue {
			// Check if current field is non-empty
			if currentField.String() == "" {
				return NewValidationError(ctx.Path(), "required_if", map[string]any{
					"field":         otherFieldName,
					"expectedValue": expectedValue,
				})
			}
		}

		return nil
	})
}
