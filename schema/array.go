package schema

import (
	"fmt"

	"github.com/weilence/schema-validator/data"
)

// ArraySchema validates arrays/slices
type ArraySchema struct {
	elementSchema Schema

	validators     []Validator
	itemValidators []Validator
}

// NewArraySchema creates a new array schema
func NewArraySchema(elementSchema Schema) *ArraySchema {
	return &ArraySchema{
		elementSchema: elementSchema,
		validators:    make([]Validator, 0),
	}
}

// Validate validates an array
func (a *ArraySchema) Validate(ctx *Context) error {
	for _, validator := range a.validators {
		if ctx.skipRest {
			break
		}

		if err := validator.Validate(ctx); err != nil {
			return err
		}
	}

	// Validate each element
	accessor, ok := ctx.Accessor().(*data.ArrayAccessor)
	if !ok {
		return fmt.Errorf("expected ArrayAccessor, got %T", ctx.Accessor())
	}

	return accessor.Iterate(func(idx int, childAccessor data.Accessor) error {
		elemCtx := ctx.WithChild(fmt.Sprintf("[%d]", idx), a.elementSchema, childAccessor)
		return a.elementSchema.Validate(elemCtx)
	})
}

func (a *ArraySchema) AddValidator(v Validator) Schema {
	a.validators = append(a.validators, v)
	return a
}

func (a *ArraySchema) RemoveValidator(name string) Schema {
	newValidators := make([]Validator, 0)
	for _, v := range a.validators {
		if v.Name() != name {
			newValidators = append(newValidators, v)
		}
	}
	a.validators = newValidators
	return a
}
