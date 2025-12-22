package schema

import (
	"encoding/json"
	"fmt"

	"github.com/weilence/schema-validator/validation"
)

// FieldSchema validates primitive/scalar values
type FieldSchema struct {
	validators []validation.FieldValidator
	optional   bool
}

// NewFieldSchema creates a new field schema
func NewFieldSchema() *FieldSchema {
	return &FieldSchema{
		validators: make([]validation.FieldValidator, 0),
		optional:   false,
	}
}

// Type returns SchemaTypeField
func (f *FieldSchema) Type() SchemaType {
	return SchemaTypeField
}

// Validate validates a field value
func (f *FieldSchema) Validate(ctx *validation.Context) error {
	// 设置当前 schema
	ctx.SetSchema(f)

	// Check if value is nil and optional
	if ctx.IsNil() && f.optional {
		return nil
	}

	// Run all validators
	for _, validator := range f.validators {
		if err := validator.Validate(ctx); err != nil {
			return err
		}
	}

	return nil
}

// AddValidator adds a field validator
func (f *FieldSchema) AddValidator(validator validation.FieldValidator) *FieldSchema {
	f.validators = append(f.validators, validator)
	return f
}

// SetOptional marks the field as optional
func (f *FieldSchema) SetOptional(optional bool) *FieldSchema {
	f.optional = optional
	return f
}

// ToString returns a JSON representation of the field schema
func (f *FieldSchema) ToString() string {
	result := map[string]interface{}{
		"type":     "field",
		"optional": f.optional,
	}

	if len(f.validators) > 0 {
		validators := make([]map[string]interface{}, 0, len(f.validators))
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

// validatorToMap converts a validator to a map representation
func validatorToMap(v interface{}) map[string]interface{} {
	result := map[string]interface{}{}

	switch validator := v.(type) {
	case *validation.RequiredValidator:
		result["name"] = "required"
	case *validation.MinValidator:
		result["name"] = "min"
		result["value"] = validator.Min
	case *validation.MaxValidator:
		result["name"] = "max"
		result["value"] = validator.Max
	case *validation.MinLengthValidator:
		result["name"] = "min_length"
		result["value"] = validator.MinLength
	case *validation.MaxLengthValidator:
		result["name"] = "max_length"
		result["value"] = validator.MaxLength
	case *validation.EmailValidator:
		result["name"] = "email"
	case *validation.URLValidator:
		result["name"] = "url"
	case *validation.PatternValidator:
		result["name"] = "pattern"
		result["value"] = validator.Pattern.String()
	case *validation.EqFieldValidator:
		result["name"] = "eqfield"
		result["value"] = validator.FieldName
	case *validation.NeFieldValidator:
		result["name"] = "nefield"
		result["value"] = validator.FieldName
	case *validation.GtFieldValidator:
		result["name"] = "gtfield"
		result["value"] = validator.FieldName
	case *validation.LtFieldValidator:
		result["name"] = "ltfield"
		result["value"] = validator.FieldName
	default:
		result["name"] = "custom"
		result["type"] = fmt.Sprintf("%T", v)
	}

	return result
}
