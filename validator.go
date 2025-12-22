package validator

import (
	"fmt"
	"reflect"

	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/errors"
	"github.com/weilence/schema-validator/schema"
	"github.com/weilence/schema-validator/tags"
	"github.com/weilence/schema-validator/validation"
)

// Engine is the validation execution engine
type Engine struct {
	schema schema.Schema
}

// newEngine creates a new validation engine
func newEngine(s schema.Schema) *Engine {
	return &Engine{schema: s}
}

// Validate validates data and returns validation result
func (e *Engine) Validate(dataValue interface{}) (*errors.ValidationResult, error) {
	// Create data accessor
	accessor := data.NewValue(dataValue)

	// Create validation context
	ctx := validation.NewContext(accessor)

	// Execute validation
	result := errors.NewValidationResult()

	if err := e.validateRecursive(ctx, accessor); err != nil {
		if verr, ok := err.(*errors.ValidationError); ok {
			result.AddError(verr)
		} else {
			// System error, not validation error
			return nil, err
		}
	}

	return result, nil
}

// validateRecursive recursively validates data and collects all errors
func (e *Engine) validateRecursive(ctx *validation.Context, accessor data.Accessor) error {
	return e.schema.Validate(ctx, accessor)
}

// Validator is the main entry point for validation
type Validator struct {
	schema schema.Schema
	engine *Engine
}

// New creates a validator from a code-based schema
func New(s schema.Schema) *Validator {
	return &Validator{
		schema: s,
		engine: newEngine(s),
	}
}

// NewFromStruct creates a validator from struct tags
func NewFromStruct(prototype interface{}) (*Validator, error) {
	typ := reflect.TypeOf(prototype)
	if typ == nil {
		return nil, fmt.Errorf("prototype is nil")
	}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %v", typ.Kind())
	}

	objSchema, err := tags.ParseStructTags(typ)
	if err != nil {
		return nil, err
	}

	return New(objSchema), nil
}

// Validate validates data and returns validation result
func (v *Validator) Validate(data interface{}) (*errors.ValidationResult, error) {
	return v.engine.Validate(data)
}

// ValidateStruct is a convenience method for validating structs
func (v *Validator) ValidateStruct(s interface{}) (*errors.ValidationResult, error) {
	return v.Validate(s)
}

// MustValidate validates and panics on system errors (not validation errors)
func (v *Validator) MustValidate(data interface{}) *errors.ValidationResult {
	result, err := v.Validate(data)
	if err != nil {
		panic(err)
	}
	return result
}

// IsValid validates and returns true if data is valid
func (v *Validator) IsValid(data interface{}) bool {
	result, err := v.Validate(data)
	if err != nil {
		return false
	}
	return result.IsValid()
}
