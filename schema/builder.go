package schema

import (
	"fmt"
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
	b.schema.AddValidator("required")
	return b
}

// Optional marks the field as optional
func (b *FieldSchemaBuilder) Optional() *FieldSchemaBuilder {
	b.schema.RemoveValidator("required")
	return b
}

// AddValidator adds a custom validator
func (b *FieldSchemaBuilder) AddValidator(name string, params ...string) *FieldSchemaBuilder {
	b.schema.AddValidator(name, params...)
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
	b.schema.AddValidator("min_items", fmt.Sprint(min))
	return b
}

// MaxItems sets maximum items constraint
func (b *ArraySchemaBuilder) MaxItems(max int) *ArraySchemaBuilder {
	b.schema.SetMaxItems(max)
	b.schema.AddValidator("max_items", fmt.Sprint(max))
	return b
}

// AddValidator adds a custom array validator
func (b *ArraySchemaBuilder) AddValidator(name string, params ...string) *ArraySchemaBuilder {
	b.schema.AddValidator(name, params...)
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

func (b *ObjectSchemaBuilder) FieldName(name string, fieldName string) *ObjectSchemaBuilder {
	b.schema.AddFieldName(name, fieldName)
	return b
}

// AddValidator adds a custom object validator
func (b *ObjectSchemaBuilder) AddValidator(name string, params ...string) *ObjectSchemaBuilder {
	b.schema.AddValidator(name, params...)
	return b
}

// AddValidatorByName adds a validator by name from the global registry
func (b *ObjectSchemaBuilder) AddValidatorByName(name string, params ...string) *ObjectSchemaBuilder {
	b.schema.AddValidator(name, params...)
	return b
}

// Build returns the built object schema
func (b *ObjectSchemaBuilder) Build() *ObjectSchema {
	return b.schema
}
