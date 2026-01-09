package validator

import (
	"reflect"

	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/schema"
)

// Validator is the main entry point for validation
type Validator struct {
	schema schema.Schema
}

func New(prototype any, opts ...ParseOption) (*Validator, error) {
	objSchema, err := Parse(reflect.TypeOf(prototype), opts...)
	if err != nil {
		return nil, err
	}

	return NewFromSchema(objSchema), nil
}

// NewFromSchema creates a validator from a code-based schema
func NewFromSchema(s schema.Schema) *Validator {
	return &Validator{
		schema: s,
	}
}

// Validate validates data and returns validation result
func (v *Validator) Validate(value any) error {
	// Create data accessor
	accessor := data.New(value)

	// Create validation context
	ctx := schema.NewContext(v.schema, accessor)
	err := v.schema.Validate(ctx)
	if err != nil {
		return err
	}

	errs := ctx.Errors()
	if len(errs) == 0 {
		return nil
	}

	return errs
}
