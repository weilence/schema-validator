package validation

import (
	"github.com/weilence/schema-validator/errors"
)

// EqFieldValidator validates that a field equals another field
type EqFieldValidator struct {
	FieldName string
}

func (v *EqFieldValidator) Validate(ctx *Context) error {
	// 直接通过 ctx 获取父级对象
	parentObj, err := ctx.GetParentObject()
	if err != nil {
		return err
	}

	otherField, exists := parentObj.GetField(v.FieldName)
	if !exists {
		return errors.NewValidationError(ctx.Path(), "eqfield",
			map[string]interface{}{"field": v.FieldName, "error": "field not found"})
	}

	otherValue, err := otherField.AsField()
	if err != nil {
		return err
	}

	currentField, _ := ctx.AsField() // 使用缓存

	if currentField.String() != otherValue.String() {
		return errors.NewValidationError(ctx.Path(), "eqfield",
			map[string]interface{}{"field": v.FieldName})
	}

	return nil
}

// NeFieldValidator validates that a field does not equal another field
type NeFieldValidator struct {
	FieldName string
}

func (v *NeFieldValidator) Validate(ctx *Context) error {
	// 直接通过 ctx 获取父级对象
	parentObj, err := ctx.GetParentObject()
	if err != nil {
		return err
	}

	otherField, exists := parentObj.GetField(v.FieldName)
	if !exists {
		// If field doesn't exist, it's not equal
		return nil
	}

	otherValue, err := otherField.AsField()
	if err != nil {
		return err
	}

	currentField, _ := ctx.AsField() // 使用缓存

	if currentField.String() == otherValue.String() {
		return errors.NewValidationError(ctx.Path(), "nefield",
			map[string]interface{}{"field": v.FieldName})
	}

	return nil
}

// GtFieldValidator validates that a field is greater than another field
type GtFieldValidator struct {
	FieldName string
}

func (v *GtFieldValidator) Validate(ctx *Context) error {
	// 直接通过 ctx 获取父级对象
	parentObj, err := ctx.GetParentObject()
	if err != nil {
		return err
	}

	otherField, exists := parentObj.GetField(v.FieldName)
	if !exists {
		return errors.NewValidationError(ctx.Path(), "gtfield",
			map[string]interface{}{"field": v.FieldName, "error": "field not found"})
	}

	otherValue, err := otherField.AsField()
	if err != nil {
		return err
	}

	currentField, _ := ctx.AsField() // 使用缓存

	// Try numeric comparison
	val, err1 := currentField.Int()
	otherVal, err2 := otherValue.Int()

	if err1 == nil && err2 == nil {
		if val <= otherVal {
			return errors.NewValidationError(ctx.Path(), "gtfield",
				map[string]interface{}{"field": v.FieldName})
		}
		return nil
	}

	// Try float comparison
	fval, err1 := currentField.Float()
	fotherVal, err2 := otherValue.Float()

	if err1 == nil && err2 == nil {
		if fval <= fotherVal {
			return errors.NewValidationError(ctx.Path(), "gtfield",
				map[string]interface{}{"field": v.FieldName})
		}
		return nil
	}

	// String comparison
	if currentField.String() <= otherValue.String() {
		return errors.NewValidationError(ctx.Path(), "gtfield",
			map[string]interface{}{"field": v.FieldName})
	}

	return nil
}

// LtFieldValidator validates that a field is less than another field
type LtFieldValidator struct {
	FieldName string
}

func (v *LtFieldValidator) Validate(ctx *Context) error {
	// 直接通过 ctx 获取父级对象
	parentObj, err := ctx.GetParentObject()
	if err != nil {
		return err
	}

	otherField, exists := parentObj.GetField(v.FieldName)
	if !exists {
		return errors.NewValidationError(ctx.Path(), "ltfield",
			map[string]interface{}{"field": v.FieldName, "error": "field not found"})
	}

	otherValue, err := otherField.AsField()
	if err != nil {
		return err
	}

	currentField, _ := ctx.AsField() // 使用缓存

	// Try numeric comparison
	val, err1 := currentField.Int()
	otherVal, err2 := otherValue.Int()

	if err1 == nil && err2 == nil {
		if val >= otherVal {
			return errors.NewValidationError(ctx.Path(), "ltfield",
				map[string]interface{}{"field": v.FieldName})
		}
		return nil
	}

	// Try float comparison
	fval, err1 := currentField.Float()
	fotherVal, err2 := otherValue.Float()

	if err1 == nil && err2 == nil {
		if fval >= fotherVal {
			return errors.NewValidationError(ctx.Path(), "ltfield",
				map[string]interface{}{"field": v.FieldName})
		}
		return nil
	}

	// String comparison
	if currentField.String() >= otherValue.String() {
		return errors.NewValidationError(ctx.Path(), "ltfield",
			map[string]interface{}{"field": v.FieldName})
	}

	return nil
}

// Convenience functions

// EqField creates an equal field validator
func EqField(fieldName string) FieldValidator {
	return &EqFieldValidator{FieldName: fieldName}
}

// NeField creates a not equal field validator
func NeField(fieldName string) FieldValidator {
	return &NeFieldValidator{FieldName: fieldName}
}

// GtField creates a greater than field validator
func GtField(fieldName string) FieldValidator {
	return &GtFieldValidator{FieldName: fieldName}
}

// LtField creates a less than field validator
func LtField(fieldName string) FieldValidator {
	return &LtFieldValidator{FieldName: fieldName}
}
