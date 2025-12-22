package schema

import (
	"encoding/json"
	"fmt"

	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/validation"
)

// ArraySchema validates arrays/slices
type ArraySchema struct {
	elementSchema Schema
	minItems      *int
	maxItems      *int
	validators    []validation.ArrayValidator
}

// NewArraySchema creates a new array schema
func NewArraySchema(elementSchema Schema) *ArraySchema {
	return &ArraySchema{
		elementSchema: elementSchema,
		validators:    make([]validation.ArrayValidator, 0),
	}
}

// Type returns SchemaTypeArray
func (a *ArraySchema) Type() SchemaType {
	return SchemaTypeArray
}

// Validate validates an array
func (a *ArraySchema) Validate(ctx *validation.Context, accessor data.Accessor) error {
	arrAcc, err := accessor.AsArray()
	if err != nil {
		return err
	}

	// Validate array-level constraints (min/max items)
	for _, validator := range a.validators {
		if err := validator.Validate(ctx, arrAcc); err != nil {
			return err
		}
	}

	// Validate each element
	return arrAcc.Iterate(func(idx int, elem data.Accessor) error {
		elemCtx := ctx.WithPath(fmt.Sprintf("[%d]", idx))
		return a.elementSchema.Validate(elemCtx, elem)
	})
}

// AddValidator adds an array validator
func (a *ArraySchema) AddValidator(validator validation.ArrayValidator) *ArraySchema {
	a.validators = append(a.validators, validator)
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
	result := map[string]interface{}{
		"type": "array",
	}

	// Add element schema as nested JSON
	if a.elementSchema != nil {
		// Parse the nested schema's JSON string back to a map for proper nesting
		var elementMap map[string]interface{}
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
		validators := make([]map[string]interface{}, 0, len(a.validators))
		for _, v := range a.validators {
			validators = append(validators, arrayValidatorToMap(v))
		}
		result["validators"] = validators
	}

	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"type":"array","error":"%s"}`, err.Error())
	}
	return string(bytes)
}

// arrayValidatorToMap converts an array validator to a map representation
func arrayValidatorToMap(v interface{}) map[string]interface{} {
	result := map[string]interface{}{}

	switch validator := v.(type) {
	case *validation.MinItemsValidator:
		result["name"] = "min_items"
		result["value"] = validator.MinItems
	case *validation.MaxItemsValidator:
		result["name"] = "max_items"
		result["value"] = validator.MaxItems
	default:
		result["name"] = "custom"
		result["type"] = fmt.Sprintf("%T", v)
	}

	return result
}
