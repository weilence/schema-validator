package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/errors"
)

// RequiredValidator validates that a field is not nil/empty
type RequiredValidator struct{}

func (v *RequiredValidator) Validate(ctx *Context, value data.FieldAccessor) error {
	if value.IsNil() {
		return errors.NewValidationError(ctx.Path, "required", nil)
	}

	// Check for empty string
	str := value.String()
	if strings.TrimSpace(str) == "" {
		return errors.NewValidationError(ctx.Path, "required", nil)
	}

	return nil
}

// MinValidator validates minimum value/length
type MinValidator struct {
	Min interface{}
}

func (v *MinValidator) Validate(ctx *Context, value data.FieldAccessor) error {
	if value.IsNil() {
		return nil // Skip validation for nil values
	}

	switch min := v.Min.(type) {
	case int, int64:
		minVal := toInt64(v.Min)
		val, err := value.Int()
		if err == nil && val < minVal {
			return errors.NewValidationError(ctx.Path, "min", map[string]interface{}{"min": minVal, "actual": val})
		}
	case float64:
		val, err := value.Float()
		if err == nil && val < min {
			return errors.NewValidationError(ctx.Path, "min", map[string]interface{}{"min": min, "actual": val})
		}
	case string:
		// For strings, min is length
		str := value.String()
		minLen := len(min)
		if len(str) < minLen {
			return errors.NewValidationError(ctx.Path, "min_length", map[string]interface{}{"min": minLen, "actual": len(str)})
		}
	}

	return nil
}

// MaxValidator validates maximum value/length
type MaxValidator struct {
	Max interface{}
}

func (v *MaxValidator) Validate(ctx *Context, value data.FieldAccessor) error {
	if value.IsNil() {
		return nil // Skip validation for nil values
	}

	switch max := v.Max.(type) {
	case int, int64:
		maxVal := toInt64(v.Max)
		val, err := value.Int()
		if err == nil && val > maxVal {
			return errors.NewValidationError(ctx.Path, "max", map[string]interface{}{"max": maxVal, "actual": val})
		}
	case float64:
		val, err := value.Float()
		if err == nil && val > max {
			return errors.NewValidationError(ctx.Path, "max", map[string]interface{}{"max": max, "actual": val})
		}
	case string:
		// For strings, max is length
		str := value.String()
		maxLen := len(max)
		if len(str) > maxLen {
			return errors.NewValidationError(ctx.Path, "max_length", map[string]interface{}{"max": maxLen, "actual": len(str)})
		}
	}

	return nil
}

// MinLengthValidator validates minimum string length
type MinLengthValidator struct {
	MinLength int
}

func (v *MinLengthValidator) Validate(ctx *Context, value data.FieldAccessor) error {
	if value.IsNil() {
		return nil
	}

	str := value.String()
	if len(str) < v.MinLength {
		return errors.NewValidationError(ctx.Path, "min_length", map[string]interface{}{"min": v.MinLength, "actual": len(str)})
	}

	return nil
}

// MaxLengthValidator validates maximum string length
type MaxLengthValidator struct {
	MaxLength int
}

func (v *MaxLengthValidator) Validate(ctx *Context, value data.FieldAccessor) error {
	if value.IsNil() {
		return nil
	}

	str := value.String()
	if len(str) > v.MaxLength {
		return errors.NewValidationError(ctx.Path, "max_length", map[string]interface{}{"max": v.MaxLength, "actual": len(str)})
	}

	return nil
}

// EmailValidator validates email format
type EmailValidator struct{}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func (v *EmailValidator) Validate(ctx *Context, value data.FieldAccessor) error {
	if value.IsNil() {
		return nil
	}

	str := value.String()
	if !emailRegex.MatchString(str) {
		return errors.NewValidationError(ctx.Path, "email", nil)
	}

	return nil
}

// URLValidator validates URL format
type URLValidator struct{}

var urlRegex = regexp.MustCompile(`^https?://[^\s]+$`)

func (v *URLValidator) Validate(ctx *Context, value data.FieldAccessor) error {
	if value.IsNil() {
		return nil
	}

	str := value.String()
	if !urlRegex.MatchString(str) {
		return errors.NewValidationError(ctx.Path, "url", nil)
	}

	return nil
}

// PatternValidator validates against a regex pattern
type PatternValidator struct {
	Pattern *regexp.Regexp
}

func (v *PatternValidator) Validate(ctx *Context, value data.FieldAccessor) error {
	if value.IsNil() {
		return nil
	}

	str := value.String()
	if !v.Pattern.MatchString(str) {
		return errors.NewValidationError(ctx.Path, "pattern", map[string]interface{}{"pattern": v.Pattern.String()})
	}

	return nil
}

// MinItemsValidator validates minimum array items
type MinItemsValidator struct {
	MinItems int
}

func (v *MinItemsValidator) Validate(ctx *Context, arr data.ArrayAccessor) error {
	if arr.IsNil() {
		return nil
	}

	if arr.Len() < v.MinItems {
		return errors.NewValidationError(ctx.Path, "min_items", map[string]interface{}{"min": v.MinItems, "actual": arr.Len()})
	}

	return nil
}

// MaxItemsValidator validates maximum array items
type MaxItemsValidator struct {
	MaxItems int
}

func (v *MaxItemsValidator) Validate(ctx *Context, arr data.ArrayAccessor) error {
	if arr.IsNil() {
		return nil
	}

	if arr.Len() > v.MaxItems {
		return errors.NewValidationError(ctx.Path, "max_items", map[string]interface{}{"max": v.MaxItems, "actual": arr.Len()})
	}

	return nil
}

// Helper function to convert to int64
func toInt64(v interface{}) int64 {
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

// Convenience functions to create validators

// Required creates a required validator
func Required() FieldValidator {
	return &RequiredValidator{}
}

// Min creates a min validator
func Min(min interface{}) FieldValidator {
	return &MinValidator{Min: min}
}

// Max creates a max validator
func Max(max interface{}) FieldValidator {
	return &MaxValidator{Max: max}
}

// MinLength creates a min length validator
func MinLength(min int) FieldValidator {
	return &MinLengthValidator{MinLength: min}
}

// MaxLength creates a max length validator
func MaxLength(max int) FieldValidator {
	return &MaxLengthValidator{MaxLength: max}
}

// Email creates an email validator
func Email() FieldValidator {
	return &EmailValidator{}
}

// URL creates a URL validator
func URL() FieldValidator {
	return &URLValidator{}
}

// Pattern creates a pattern validator
func Pattern(pattern string) FieldValidator {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("invalid regex pattern: %s", pattern))
	}
	return &PatternValidator{Pattern: regex}
}

// MinItems creates a min items validator
func MinItems(min int) ArrayValidator {
	return &MinItemsValidator{MinItems: min}
}

// MaxItems creates a max items validator
func MaxItems(max int) ArrayValidator {
	return &MaxItemsValidator{MaxItems: max}
}
