package schema

// SchemaModifier is an interface that allows structs to modify their validation schema dynamically
// The struct can implement this interface to add/remove validation rules based on runtime values
type SchemaModifier interface {
	// ModifySchema is called before validation with access to the current object's data
	// ctx provides access to:
	//   - validation context (path, parent, root)
	//   - current object's accessor (via ctx.AsObject())
	//   - current ObjectSchema (via ctx.ObjectSchema())
	ModifySchema(ctx *Context)
}

// Schema represents a validation schema for any data type
type Schema interface {
	// Validate validates data against this schema
	// ctx contains both the validation context and the data accessor
	Validate(ctx *Context) error

	AddValidator(v Validator) Schema

	RemoveValidator(name string) Schema
}
