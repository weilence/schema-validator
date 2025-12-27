package builder

import (
	"github.com/weilence/schema-validator/schema"
	"github.com/weilence/schema-validator/validators"
)

// SchemaBuilder provides a unified fluent API for building different schema types
type SchemaBuilder struct {
	schema   schema.Schema
	registry *validators.Registry
}

// Field creates a new field schema builder
func Field() *SchemaBuilder {
	return &SchemaBuilder{schema: schema.NewFieldSchema(), registry: validators.DefaultRegistry()}
}

// Array creates a new array schema builder
func Array(elementSchema schema.Schema) *SchemaBuilder {
	return &SchemaBuilder{schema: schema.NewArraySchema(elementSchema), registry: validators.DefaultRegistry()}
}

// Object creates a new object schema builder
func Object() *SchemaBuilder {
	return &SchemaBuilder{schema: schema.NewObjectSchema(), registry: validators.DefaultRegistry()}
}

// FieldWithRegistry creates a new field schema builder with a custom registry
func (b *SchemaBuilder) Registry(r *validators.Registry) *SchemaBuilder {
	b.registry = r
	return b
}

// Required marks the field as required (no-op for non-field builders)
func (b *SchemaBuilder) Required() *SchemaBuilder {
	b.AddValidator("required")
	return b
}

// Optional marks the field as optional (no-op for non-field builders)
func (b *SchemaBuilder) Optional() *SchemaBuilder {
	b.schema.RemoveValidator("required")
	return b
}

// AddValidator adds a custom validator to the underlying schema
func (b *SchemaBuilder) AddValidator(name string, params ...any) *SchemaBuilder {
	v := b.registry.NewValidator(name, params...)
	b.schema.AddValidator(v)
	return b
}

// Field adds a field to the underlying object schema (no-op for non-object builders)
func (b *SchemaBuilder) Field(name string, fieldSchema schema.Schema) *SchemaBuilder {
	if os, ok := b.schema.(*schema.ObjectSchema); ok {
		os.AddField(name, fieldSchema)
	}
	return b
}

// FieldName sets the mapping for an object field (no-op for non-object builders)
func (b *SchemaBuilder) FieldName(name string, fieldName string) *SchemaBuilder {
	if os, ok := b.schema.(*schema.ObjectSchema); ok {
		os.AddFieldName(name, fieldName)
	}
	return b
}

// Build returns the built schema as the Schema interface
func (b *SchemaBuilder) Build() schema.Schema {
	return b.schema
}
