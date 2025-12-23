package schema

import (
	"encoding/json"
	"fmt"
)

// FieldSchema validates primitive/scalar values
type FieldSchema struct {
	validators []Validator
}

// NewFieldSchema creates a new field schema
func NewFieldSchema() *FieldSchema {
	return &FieldSchema{
		validators: make([]Validator, 0),
	}
}

// Type returns SchemaTypeField
func (f *FieldSchema) Type() SchemaType {
	return SchemaTypeField
}

// Validate validates a field value
func (f *FieldSchema) Validate(ctx *Context) error {
	// Run all validators
	for _, validator := range f.validators {
		if err := validator.Validate(ctx); err != nil {
			return err
		}
	}

	return nil
}

// AddValidator adds a field validator
func (f *FieldSchema) AddValidator(name string, params ...string) *FieldSchema {
	v := DefaultRegistry().BuildValidator(name, params)
	if v != nil {
		f.validators = append(f.validators, v)
	}
	return f
}

func (f *FieldSchema) RemoveValidator(name string) *FieldSchema {
	newValidators := make([]Validator, 0)
	for _, v := range f.validators {
		switch validator := v.(type) {
		case validator:
			if validator.Name() != name {
				newValidators = append(newValidators, v)
			}
		case *validator:
			if validator.Name() != name {
				newValidators = append(newValidators, v)
			}
		default:
			newValidators = append(newValidators, v)
		}
	}

	f.validators = newValidators
	return f
}

// ToString returns a JSON representation of the field schema
func (f *FieldSchema) ToString() string {
	result := map[string]any{
		"type": "field",
	}

	if len(f.validators) > 0 {
		validators := make([]map[string]any, 0, len(f.validators))
		for _, v := range f.validators {
			validators = append(validators, validatorToMap(v))
		}
		result["validators"] = validators
	}

	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"type":"field","error":"%s"}`, err.Error())
	}
	return string(bytes)
}
