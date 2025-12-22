package validation

// FieldValidator validates a single field value
type FieldValidator interface {
	Validate(ctx *Context) error
}

// ArrayValidator validates an array as a whole
type ArrayValidator interface {
	Validate(ctx *Context) error
}

// ObjectValidator validates an object as a whole (for cross-field validation)
type ObjectValidator interface {
	Validate(ctx *Context) error
}

// FieldValidatorFunc is a function that validates a field
type FieldValidatorFunc func(ctx *Context) error

// Validate implements FieldValidator
func (f FieldValidatorFunc) Validate(ctx *Context) error {
	return f(ctx)
}

// ArrayValidatorFunc is a function that validates an array
type ArrayValidatorFunc func(ctx *Context) error

// Validate implements ArrayValidator
func (f ArrayValidatorFunc) Validate(ctx *Context) error {
	return f(ctx)
}

// ObjectValidatorFunc is a function that validates an object
type ObjectValidatorFunc func(ctx *Context) error

// Validate implements ObjectValidator
func (f ObjectValidatorFunc) Validate(ctx *Context) error {
	return f(ctx)
}
