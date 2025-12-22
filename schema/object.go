package schema

import (
	"encoding/json"
	"fmt"

	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/validation"
)

// ObjectSchema validates objects/structs/maps
type ObjectSchema struct {
	fields     map[string]Schema
	validators []validation.ObjectValidator
	strict     bool // disallow unknown fields
}

// NewObjectSchema creates a new object schema
func NewObjectSchema() *ObjectSchema {
	return &ObjectSchema{
		fields:     make(map[string]Schema),
		validators: make([]validation.ObjectValidator, 0),
		strict:     false,
	}
}

// Type returns SchemaTypeObject
func (o *ObjectSchema) Type() SchemaType {
	return SchemaTypeObject
}

// Validate validates an object
func (o *ObjectSchema) Validate(ctx *validation.Context) error {
	// 设置当前 schema 到 context
	ctx.SetSchema(o)

	// 使用缓存的 AsObject
	objAcc, err := ctx.AsObject()
	if err != nil {
		return err
	}

	// Check if the underlying value implements SchemaModifier
	// This allows the struct to dynamically modify its validation schema
	if modifier, ok := ctx.Accessor().(*data.Value); ok {
		if !modifier.IsNil() {
			rawVal := modifier.Raw()
			if schemaModifier, ok := rawVal.(SchemaModifier); ok {
				// Call ModifySchema to allow dynamic schema modification
				// ctx now contains schema, accessor, and context information
				schemaModifier.ModifySchema(ctx)
			}
		}
	}

	// Validate each defined field
	for fieldName, fieldSchema := range o.fields {
		fieldData, exists := objAcc.GetField(fieldName)
		if !exists {
			fieldData = data.NewValue(nil) // nil value for missing field
		}

		// 创建子 context - 自动父级追踪和 schema 传递
		fieldCtx := ctx.WithChild(fieldName, fieldData, fieldSchema)

		if err := fieldSchema.Validate(fieldCtx); err != nil {
			return err
		}
	}

	// Run object-level validators (cross-field)
	for _, validator := range o.validators {
		if err := validator.Validate(ctx); err != nil {
			return err
		}
	}

	return nil
}

// AddField adds a field schema
func (o *ObjectSchema) AddField(name string, schema Schema) *ObjectSchema {
	o.fields[name] = schema
	return o
}

// AddValidator adds an object validator
func (o *ObjectSchema) AddValidator(validator validation.ObjectValidator) *ObjectSchema {
	o.validators = append(o.validators, validator)
	return o
}

// SetStrict enables/disables strict mode
func (o *ObjectSchema) SetStrict(strict bool) *ObjectSchema {
	o.strict = strict
	return o
}

// GetField returns the schema for a field
func (o *ObjectSchema) GetField(name string) (Schema, bool) {
	schema, ok := o.fields[name]
	return schema, ok
}

// Fields returns all field names
func (o *ObjectSchema) Fields() []string {
	names := make([]string, 0, len(o.fields))
	for name := range o.fields {
		names = append(names, name)
	}
	return names
}

// ToString returns a JSON representation of the object schema
func (o *ObjectSchema) ToString() string {
	result := map[string]interface{}{
		"type":   "object",
		"strict": o.strict,
	}

	// Add fields as nested schemas
	if len(o.fields) > 0 {
		fields := make(map[string]interface{})
		for fieldName, fieldSchema := range o.fields {
			// Parse the nested schema's JSON string back to a map for proper nesting
			var fieldMap map[string]interface{}
			fieldJSON := fieldSchema.ToString()
			if err := json.Unmarshal([]byte(fieldJSON), &fieldMap); err == nil {
				fields[fieldName] = fieldMap
			} else {
				fields[fieldName] = fieldJSON
			}
		}
		result["fields"] = fields
	}

	if len(o.validators) > 0 {
		validators := make([]map[string]interface{}, 0, len(o.validators))
		for _, v := range o.validators {
			validators = append(validators, objectValidatorToMap(v))
		}
		result["validators"] = validators
	}

	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"type":"object","error":"%s"}`, err.Error())
	}
	return string(bytes)
}

// objectValidatorToMap converts an object validator to a map representation
func objectValidatorToMap(v interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"name": "custom",
		"type": fmt.Sprintf("%T", v),
	}
	return result
}
