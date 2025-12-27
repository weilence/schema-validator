package validators

import (
	"cmp"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/spf13/cast"
	"github.com/weilence/schema-validator/schema"
)

type compareType int

const (
	LessThan           compareType = iota // <
	LessThanOrEqual                       // <=
	GreaterThan                           // >
	GreaterThanOrEqual                    // >=
	Equal                                 // ==
	NotEqual                              // !=
)

func compareFn[T cmp.Ordered](t compareType, a, b T) bool {
	switch t {
	case LessThan:
		return a < b
	case LessThanOrEqual:
		return a <= b
	case GreaterThan:
		return a > b
	case GreaterThanOrEqual:
		return a >= b
	case Equal:
		return a == b
	case NotEqual:
		return a != b
	default:
		panic("unknown compare type")
	}
}

func compareValidator(ct compareType) func(*schema.Context, any) error {
	return func(ctx *schema.Context, value any) error {
		field, err := ctx.Value()
		if err != nil {
			return fmt.Errorf("failed to get field value: %v", err)
		}

		kind := field.Kind()
		if kind == reflect.String || kind == reflect.Slice || kind == reflect.Array || kind == reflect.Map {
			min, err := cast.ToIntE(value)
			if err != nil {
				return fmt.Errorf("invalid min parameter: %v", err)
			}

			if !compareFn(ct, field.Len(), min) {
				return schema.NewValidationError(ctx.Path(), "min", map[string]any{"min": min, "actual": field.Len()})
			}

			return nil
		} else if field.IsInt() {
			min, err := cast.ToInt64E(value)
			if err != nil {
				return fmt.Errorf("invalid min parameter: %v", err)
			}

			intValue, err := cast.ToInt64E(field.Raw())
			if err != nil {
				return nil
			}

			if !compareFn(ct, intValue, int64(min)) {
				return schema.NewValidationError(ctx.Path(), "min", map[string]any{"min": min, "actual": intValue})
			}
			return nil
		} else if field.IsUint() {
			min, err := cast.ToUint64E(value)
			if err != nil {
				return fmt.Errorf("invalid min parameter: %v", err)
			}

			uintValue, err := cast.ToUint64E(field.Raw())
			if err != nil {
				return nil
			}

			if !compareFn(ct, uintValue, uint64(min)) {
				return schema.NewValidationError(ctx.Path(), "min", map[string]any{"min": min, "actual": uintValue})
			}
			return nil
		} else if field.IsFloat() {
			min, err := cast.ToFloat64E(value)
			if err != nil {
				return fmt.Errorf("invalid min parameter: %v", err)
			}

			floatValue, err := cast.ToFloat64E(field.Raw())
			if err != nil {
				return nil
			}

			if !compareFn(ct, floatValue, float64(min)) {
				return schema.NewValidationError(ctx.Path(), "min", map[string]any{"min": min, "actual": floatValue})
			}
			return nil
		} else {
			return fmt.Errorf("min validator not supported for kind %s", kind.String())
		}
	}
}

func registerOther(r *Registry) {
	r.Register("min", compareValidator(GreaterThanOrEqual))

	r.Register("max", compareValidator(LessThanOrEqual))

	r.Register("oneof", func(ctx *schema.Context, params []string) error {
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

		return schema.NewValidationError(ctx.Path(), "oneof", map[string]any{
			"allowed": params,
			"actual":  val,
		})
	})

	r.Register("required", func(ctx *schema.Context) error {
		field, err := ctx.Value()
		if err != nil {
			return nil
		}
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

		currentField, err := ctx.Value()
		if err != nil {
			return fmt.Errorf("failed to get current field value: %v", err)
		}

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
