package validator

import (
	"github.com/weilence/schema-validator/schema"
	"github.com/weilence/schema-validator/rule"
)

// SchemaBuilder provides a unified fluent API for building different schema types
type SchemaBuilder struct {
	schema   schema.Schema
	registry *rule.Registry
}

// Field creates a new field schema builder
func Field() *SchemaBuilder {
	return &SchemaBuilder{schema: schema.NewField(), registry: rule.DefaultRegistry()}
}

// Array creates a new array schema builder
func Array(elementSchema schema.Schema) *SchemaBuilder {
	return &SchemaBuilder{schema: schema.NewArray(elementSchema), registry: rule.DefaultRegistry()}
}

// Object creates a new object schema builder
func Object() *SchemaBuilder {
	return &SchemaBuilder{schema: schema.NewObject(), registry: rule.DefaultRegistry()}
}

// Registry sets a custom registry for the builder
func (b *SchemaBuilder) Registry(r *rule.Registry) *SchemaBuilder {
	b.registry = r
	return b
}

// Required marks the field as required
func (b *SchemaBuilder) Required() *SchemaBuilder {
	b.AddValidator("required")
	return b
}

// Optional marks the field as optional
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

func (b *SchemaBuilder) WithField(name string, fieldSchema schema.Schema) *SchemaBuilder {
	if os, ok := b.schema.(*schema.ObjectSchema); ok {
		os.AddField(name, fieldSchema)
	}
	return b
}

// FieldName sets the mapping for an object field
func (b *SchemaBuilder) FieldName(name string, fieldName string) *SchemaBuilder {
	if os, ok := b.schema.(*schema.ObjectSchema); ok {
		os.AddFieldName(name, fieldName)
	}
	return b
}

// Build returns the built schema
func (b *SchemaBuilder) Build() schema.Schema {
	return b.schema
}
