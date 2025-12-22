package validation

import (
	"github.com/weilence/schema-validator/data"
)

// FieldValidator validates a single field value
type FieldValidator interface {
	Validate(ctx *Context, value data.FieldAccessor) error
}

// ArrayValidator validates an array as a whole
type ArrayValidator interface {
	Validate(ctx *Context, arr data.ArrayAccessor) error
}

// ObjectValidator validates an object as a whole (for cross-field validation)
type ObjectValidator interface {
	Validate(ctx *Context, obj data.ObjectAccessor) error
}

// FieldValidatorFunc is a function that validates a field
type FieldValidatorFunc func(ctx *Context, value data.FieldAccessor) error

// Validate implements FieldValidator
func (f FieldValidatorFunc) Validate(ctx *Context, value data.FieldAccessor) error {
	return f(ctx, value)
}

// ArrayValidatorFunc is a function that validates an array
type ArrayValidatorFunc func(ctx *Context, arr data.ArrayAccessor) error

// Validate implements ArrayValidator
func (f ArrayValidatorFunc) Validate(ctx *Context, arr data.ArrayAccessor) error {
	return f(ctx, arr)
}

// ObjectValidatorFunc is a function that validates an object
type ObjectValidatorFunc func(ctx *Context, obj data.ObjectAccessor) error

// Validate implements ObjectValidator
func (f ObjectValidatorFunc) Validate(ctx *Context, obj data.ObjectAccessor) error {
	return f(ctx, obj)
}
