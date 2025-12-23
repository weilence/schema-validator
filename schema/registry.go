package schema

import "fmt"

// Registry maps validator names to factory functions
type Registry struct {
	fieldValidators map[string]validatorFactory
}

var defaultRegistry = NewRegistry()

// NewRegistry creates a new validator registry
func NewRegistry() *Registry {
	r := &Registry{
		fieldValidators: make(map[string]validatorFactory),
	}

	return r
}

// Register registers a field validator factory
func (r *Registry) Register(name string, fn func(ctx *Context, params []string) error) {
	r.fieldValidators[name] = validatorFactory{
		name: name,
		fn:   fn,
	}
}

// BuildValidator gets a field validator by name
// params is a slice of parameter strings
func (r *Registry) BuildValidator(name string, params []string) Validator {
	factory, ok := r.fieldValidators[name]
	if !ok {
		panic(fmt.Sprintf("validator '%s' not found in registry", name))
	}

	return factory.Build(params)
}

// DefaultRegistry returns the default registry
func DefaultRegistry() *Registry {
	return defaultRegistry
}

func Register(name string, fn func(ctx *Context, params []string) error) {
	defaultRegistry.fieldValidators[name] = validatorFactory{
		name: name,
		fn:   fn,
	}
}
