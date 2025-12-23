package schema

import (
	"fmt"
	"net"
	"regexp"
	"slices"
	"strings"

	"github.com/miekg/dns"
	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/errors"
)

func init() {
	r := DefaultRegistry()

	// Field validators
	r.Register("required", Required)
	r.Register("min", Min)
	r.Register("max", Max)
	r.Register("min_length", MinLength)
	r.Register("max_length", MaxLength)
	r.Register("email", Email)
	r.Register("url", URL)
	r.Register("pattern", Pattern)

	// Cross-field validators
	r.Register("eqfield", EqField)
	r.Register("nefield", NeField)
	r.Register("gtfield", GtField)
	r.Register("ltfield", LtField)

	// Array validators
	r.Register("min_items", MinItems)
	r.Register("max_items", MaxItems)

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

		return errors.NewValidationError(ctx.Path(), "oneof", map[string]any{
			"allowed": params,
			"actual":  val,
		})
	})

	Register("required_if", func(ctx *Context, params []string) error {
		if len(params) < 2 {
			return fmt.Errorf("required_if validator needs 2 parameters")
		}

		otherFieldName := params[0]
		expectedValue := params[1]

		otherValue, err := ctx.GetValue(otherFieldName)
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
				return errors.NewValidationError(ctx.Path(), "required_if", map[string]any{
					"field":         otherFieldName,
					"expectedValue": expectedValue,
				})
			}
		}

		return nil
	})

	Register("ip", func(ctx *Context, params []string) error {
		field, err := ctx.Value()
		if err != nil {
			return nil
		}

		val := field.String()
		if net.ParseIP(val) != nil {
			return nil
		}

		return errors.NewValidationError(ctx.Path(), "ip", map[string]any{
			"actual": val,
		})
	})

	Register("port", func(ctx *Context, params []string) error {
		field, err := ctx.Value()
		if err != nil {
			return nil
		}

		val, err := field.Int()
		if err != nil {
			return nil
		}

		if val >= 1 && val <= 65535 {
			return nil
		}

		return errors.NewValidationError(ctx.Path(), "port", map[string]any{
			"actual": val,
		})
	})

	Register("domain", func(ctx *Context, params []string) error {
		field, err := ctx.Value()
		if err != nil {
			return nil
		}

		val := field.String()
		_, ok := dns.IsDomainName(val)
		if ok {
			return nil
		}

		return errors.NewValidationError(ctx.Path(), "domain", map[string]any{
			"actual": val,
		})
	})

}

// Factory functions (used by registry)
var Required = func(ctx *Context, params []string) error {
	field, err := ctx.Value()
	if err != nil {
		return nil
	}

	// Check for empty string
	str := field.String()
	if strings.TrimSpace(str) == "" {
		return errors.NewValidationError(ctx.Path(), "required", nil)
	}

	return nil
}

var Min = func(ctx *Context, params []string) error {
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
			return errors.NewValidationError(ctx.Path(), "min", map[string]any{"min": minVal, "actual": val})
		}
	case float64:
		val, err := field.Float()
		if err == nil && val < m {
			return errors.NewValidationError(ctx.Path(), "min", map[string]any{"min": m, "actual": val})
		}
	case string:
		// For strings, min is length
		str := field.String()
		minLen := len(m)
		if len(str) < minLen {
			return errors.NewValidationError(ctx.Path(), "min_length", map[string]any{"min": minLen, "actual": len(str)})
		}
	}

	return nil
}

var Max = func(ctx *Context, params []string) error {
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
			return errors.NewValidationError(ctx.Path(), "max", map[string]any{"max": maxVal, "actual": val})
		}
	case float64:
		val, err := field.Float()
		if err == nil && val > m {
			return errors.NewValidationError(ctx.Path(), "max", map[string]any{"max": m, "actual": val})
		}
	case string:
		// For strings, max is length
		str := field.String()
		maxLen := len(m)
		if len(str) > maxLen {
			return errors.NewValidationError(ctx.Path(), "max_length", map[string]any{"max": maxLen, "actual": len(str)})
		}
	}

	return nil
}

var MinLength = func(ctx *Context, params []string) error {
	if len(params) == 0 {
		return nil
	}
	min := parseInt(params[0])

	field, err := ctx.Value()
	if err != nil {
		return nil
	}

	str := field.String()
	if len(str) < min {
		return errors.NewValidationError(ctx.Path(), "min_length", map[string]any{"min": min, "actual": len(str)})
	}
	return nil
}

var MaxLength = func(ctx *Context, params []string) error {
	if len(params) == 0 {
		return nil
	}
	max := parseInt(params[0])

	field, err := ctx.Value()
	if err != nil {
		return nil
	}

	str := field.String()
	if len(str) > max {
		return errors.NewValidationError(ctx.Path(), "max_length", map[string]any{"max": max, "actual": len(str)})
	}
	return nil
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
var Email = func(ctx *Context, params []string) error {
	field, err := ctx.Value()
	if err != nil {
		return nil
	}

	str := field.String()
	if !emailRegex.MatchString(str) {
		return errors.NewValidationError(ctx.Path(), "email", nil)
	}
	return nil
}

var urlRegex = regexp.MustCompile(`^https?://[^\s]+$`)
var URL = func(ctx *Context, params []string) error {
	field, err := ctx.Value()
	if err != nil {
		return nil
	}

	str := field.String()
	if !urlRegex.MatchString(str) {
		return errors.NewValidationError(ctx.Path(), "url", nil)
	}
	return nil
}

var Pattern = func(ctx *Context, params []string) error {
	if len(params) == 0 {
		return nil
	}
	pattern := params[0]
	regex, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("invalid regex pattern: %s", pattern))
	}

	field, err := ctx.Value()
	if err != nil {
		return nil
	}

	str := field.String()
	if !regex.MatchString(str) {
		return errors.NewValidationError(ctx.Path(), "pattern", map[string]any{"pattern": regex.String()})
	}
	return nil
}

var MinItems = func(ctx *Context, params []string) error {
	if len(params) == 0 {
		return nil
	}
	min := parseInt(params[0])

	field, ok := ctx.Accessor().(*data.ArrayAccessor)
	if !ok {
		return nil
	}

	if field.Len() < min {
		return errors.NewValidationError(ctx.Path(), "min_items", map[string]any{"min": min, "actual": field.Len()})
	}
	return nil
}

var MaxItems = func(ctx *Context, params []string) error {
	if len(params) == 0 {
		return nil
	}
	max := parseInt(params[0])

	field, ok := ctx.Accessor().(*data.ArrayAccessor)
	if !ok {
		return nil
	}

	if field.Len() > max {
		return errors.NewValidationError(ctx.Path(), "max_items", map[string]any{"max": max, "actual": field.Len()})
	}
	return nil
}

// Cross-field factory functions
var EqField = func(ctx *Context, params []string) error {
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
		return errors.NewValidationError(ctx.Path(), "eqfield",
			map[string]any{"field": fieldName})
	}

	return nil
}

var NeField = func(ctx *Context, params []string) error {
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
		return errors.NewValidationError(ctx.Path(), "nefield",
			map[string]any{"field": fieldName})
	}

	return nil
}

var GtField = func(ctx *Context, params []string) error {
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

	// Try numeric comparison
	val, err1 := currentValue.Int()
	otherVal, err2 := otherValue.Int()

	if err1 == nil && err2 == nil {
		if val <= otherVal {
			return errors.NewValidationError(ctx.Path(), "gtfield",
				map[string]any{"field": fieldName})
		}
		return nil
	}

	// Try float comparison
	fval, err1 := currentValue.Float()
	fotherVal, err2 := otherValue.Float()

	if err1 == nil && err2 == nil {
		if fval <= fotherVal {
			return errors.NewValidationError(ctx.Path(), "gtfield",
				map[string]any{"field": fieldName})
		}
		return nil
	}

	// String comparison
	if currentValue.String() <= otherValue.String() {
		return errors.NewValidationError(ctx.Path(), "gtfield",
			map[string]any{"field": fieldName})
	}

	return nil
}

var LtField = func(ctx *Context, params []string) error {
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

	// Try numeric comparison
	val, err1 := currentValue.Int()
	otherVal, err2 := otherValue.Int()

	if err1 == nil && err2 == nil {
		if val >= otherVal {
			return errors.NewValidationError(ctx.Path(), "ltfield",
				map[string]any{"field": fieldName})
		}
		return nil
	}

	// Try float comparison
	fval, err1 := currentValue.Float()
	fotherVal, err2 := otherValue.Float()

	if err1 == nil && err2 == nil {
		if fval >= fotherVal {
			return errors.NewValidationError(ctx.Path(), "ltfield",
				map[string]any{"field": fieldName})
		}
		return nil
	}

	// String comparison
	if currentValue.String() >= otherValue.String() {
		return errors.NewValidationError(ctx.Path(), "ltfield",
			map[string]any{"field": fieldName})
	}

	return nil
}

// Helper function to convert to int64
func toInt64(v any) int64 {
	switch val := v.(type) {
	case int:
		return int64(val)
	case int64:
		return val
	case int32:
		return int64(val)
	case int16:
		return int64(val)
	case int8:
		return int64(val)
	default:
		return 0
	}
}

func parseInt(s string) int {
	var result int
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			result = result*10 + int(s[i]-'0')
		}
	}
	return result
}

func parseIntOrString(s string) any {
	// Try to parse as int
	isNum := true
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			isNum = false
			break
		}
	}

	if isNum && len(s) > 0 {
		return parseInt(s)
	}

	return s
}
