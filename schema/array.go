package schema

import (
	"encoding/json"
	"fmt"

	"github.com/weilence/schema-validator/data"
)

// ArraySchema validates arrays/slices
type ArraySchema struct {
	elementSchema Schema
	minItems      *int
	maxItems      *int
	validators    []Validator
}

// NewArraySchema creates a new array schema
func NewArraySchema(elementSchema Schema) *ArraySchema {
	return &ArraySchema{
		elementSchema: elementSchema,
		validators:    make([]Validator, 0),
	}
}

// Type returns SchemaTypeArray
func (a *ArraySchema) Type() SchemaType {
	return SchemaTypeArray
}

// Validate validates an array
func (a *ArraySchema) Validate(ctx *Context) error {
	for _, validator := range a.validators {
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

// AddValidatorByName adds an array validator by name from the global registry
func (a *ArraySchema) AddValidator(name string, params ...string) *ArraySchema {
	v := DefaultRegistry().BuildValidator(name, params)
	if v != nil {
		a.validators = append(a.validators, v)
	}
	return a
}

// SetMinItems sets minimum items constraint
func (a *ArraySchema) SetMinItems(min int) *ArraySchema {
	a.minItems = &min
	return a
}

// SetMaxItems sets maximum items constraint
func (a *ArraySchema) SetMaxItems(max int) *ArraySchema {
	a.maxItems = &max
	return a
}

// GetMinItems returns minimum items constraint
func (a *ArraySchema) GetMinItems() *int {
	return a.minItems
}

// GetMaxItems returns maximum items constraint
func (a *ArraySchema) GetMaxItems() *int {
	return a.maxItems
}

// ToString returns a JSON representation of the array schema
func (a *ArraySchema) ToString() string {
	result := map[string]any{
		"type": "array",
	}

	// Add element schema as nested JSON
	if a.elementSchema != nil {
		// Parse the nested schema's JSON string back to a map for proper nesting
		var elementMap map[string]any
		elementJSON := a.elementSchema.ToString()
		if err := json.Unmarshal([]byte(elementJSON), &elementMap); err == nil {
			result["element"] = elementMap
		} else {
			result["element"] = elementJSON
		}
	}

	if a.minItems != nil {
		result["minItems"] = *a.minItems
	}

	if a.maxItems != nil {
		result["maxItems"] = *a.maxItems
	}

	if len(a.validators) > 0 {
		validators := make([]map[string]any, 0, len(a.validators))
		for _, v := range a.validators {
			validators = append(validators, validatorToMap(v))
		}
		result["validators"] = validators
	}

	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"type":"array","error":"%s"}`, err.Error())
	}
	return string(bytes)
}
