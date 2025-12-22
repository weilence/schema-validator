package schema

import (
	"github.com/weilence/schema-validator/validation"
)

// FieldSchemaBuilder provides fluent API for building field schemas
type FieldSchemaBuilder struct {
	schema *FieldSchema
}

// Field creates a new field schema builder
func Field() *FieldSchemaBuilder {
	return &FieldSchemaBuilder{
		schema: NewFieldSchema(),
	}
}

// Required marks the field as required
func (b *FieldSchemaBuilder) Required() *FieldSchemaBuilder {
	b.schema.SetOptional(false)
	return b
}

// Optional marks the field as optional
func (b *FieldSchemaBuilder) Optional() *FieldSchemaBuilder {
	b.schema.SetOptional(true)
	return b
}

// SetOptional sets the optional flag
func (b *FieldSchemaBuilder) SetOptional(optional bool) *FieldSchemaBuilder {
	b.schema.SetOptional(optional)
	return b
}

// AddValidator adds a custom validator
func (b *FieldSchemaBuilder) AddValidator(validator validation.FieldValidator) *FieldSchemaBuilder {
	b.schema.AddValidator(validator)
	return b
}

// Build returns the built field schema
func (b *FieldSchemaBuilder) Build() *FieldSchema {
	return b.schema
}

// ArraySchemaBuilder provides fluent API for building array schemas
type ArraySchemaBuilder struct {
	schema *ArraySchema
}

// Array creates a new array schema builder
func Array(elementSchema Schema) *ArraySchemaBuilder {
	return &ArraySchemaBuilder{
		schema: NewArraySchema(elementSchema),
	}
}

// MinItems sets minimum items constraint
func (b *ArraySchemaBuilder) MinItems(min int) *ArraySchemaBuilder {
	b.schema.SetMinItems(min)
	b.schema.AddValidator(&validation.MinItemsValidator{MinItems: min})
	return b
}

// MaxItems sets maximum items constraint
func (b *ArraySchemaBuilder) MaxItems(max int) *ArraySchemaBuilder {
	b.schema.SetMaxItems(max)
	b.schema.AddValidator(&validation.MaxItemsValidator{MaxItems: max})
	return b
}

// AddValidator adds a custom array validator
func (b *ArraySchemaBuilder) AddValidator(validator validation.ArrayValidator) *ArraySchemaBuilder {
	b.schema.AddValidator(validator)
	return b
}

// Build returns the built array schema
func (b *ArraySchemaBuilder) Build() *ArraySchema {
	return b.schema
}

// ObjectSchemaBuilder provides fluent API for building object schemas
type ObjectSchemaBuilder struct {
	schema *ObjectSchema
}

// Object creates a new object schema builder
func Object() *ObjectSchemaBuilder {
	return &ObjectSchemaBuilder{
		schema: NewObjectSchema(),
	}
}

// Field adds a field to the object schema
func (b *ObjectSchemaBuilder) Field(name string, fieldSchema Schema) *ObjectSchemaBuilder {
	b.schema.AddField(name, fieldSchema)
	return b
}

// Strict enables strict mode (disallow unknown fields)
func (b *ObjectSchemaBuilder) Strict() *ObjectSchemaBuilder {
	b.schema.SetStrict(true)
	return b
}

// CrossField adds a cross-field validator
func (b *ObjectSchemaBuilder) CrossField(validator validation.ObjectValidator) *ObjectSchemaBuilder {
	b.schema.AddValidator(validator)
	return b
}

// AddValidator adds a custom object validator
func (b *ObjectSchemaBuilder) AddValidator(validator validation.ObjectValidator) *ObjectSchemaBuilder {
	b.schema.AddValidator(validator)
	return b
}

// Build returns the built object schema
func (b *ObjectSchemaBuilder) Build() *ObjectSchema {
	return b.schema
}
