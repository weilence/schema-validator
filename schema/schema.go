package schema

import (
	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/validation"
)

// SchemaModifier is an interface that allows structs to modify their validation schema dynamically
// The struct can implement this interface to add/remove validation rules based on runtime values
type SchemaModifier interface {
	// ModifySchema is called before validation with access to the current object's data
	// ctx provides access to the validation context including parent and root objects
	// accessor provides access to the current object's field values
	// schema is the current ObjectSchema that can be modified
	ModifySchema(ctx *validation.Context, accessor data.ObjectAccessor, schema *ObjectSchema)
}

// SchemaType represents the type of schema
type SchemaType int

const (
	SchemaTypeField SchemaType = iota
	SchemaTypeArray
	SchemaTypeObject
)

// Schema represents a validation schema for any data type
type Schema interface {
	// Validate validates data against this schema
	Validate(ctx *validation.Context, accessor data.Accessor) error

	// Type returns the schema type (field/array/object)
	Type() SchemaType

	// ToString returns a JSON representation of the schema
	ToString() string
}
