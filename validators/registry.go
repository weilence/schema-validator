package validators

import (
	"fmt"
	"reflect"

	"github.com/weilence/schema-validator/schema"
)

// FieldValidator validates a single field value
type validator struct {
	name   string
	params []any
	fn     func(ctx *schema.Context, params []any) error
}

func (v validator) Name() string {
	return v.name
}

func (v validator) Params() []any {
	return v.params
}

// Validate implements FieldValidator
func (v validator) Validate(ctx *schema.Context) error {
	return v.fn(ctx, v.params)
}

// validatorFactory creates a validator from tag parameters
// params is a slice of parameter strings extracted from the tag
// For example: validate:"between=10,20" -> params = ["10", "20"]
type validatorFactory struct {
	name       string
	paramTypes []reflect.Type
	fn         func(ctx *schema.Context, params []any) error
}

func (vf validatorFactory) Build(params []any) schema.Validator {
	return &validator{
		name:   vf.name,
		params: params,
		fn:     vf.fn,
	}
}

// Registry maps validator names to factory functions
type Registry struct {
	validators map[string]validatorFactory
}

// NewRegistry creates a new validator registry
func NewRegistry() *Registry {
	return &Registry{
		validators: make(map[string]validatorFactory),
	}
}

// Register registers a field validator factory
func (r *Registry) Register(name string, fn any) {
	rv := reflect.ValueOf(fn)
	rvType := rv.Type()
	if rvType.Kind() != reflect.Func {
		panic("validator factory must be a function")
	}
	if rvType.NumIn() < 1 || rvType.In(0) != reflect.TypeFor[*schema.Context]() {
		panic("first parameter of validator factory must be *schema.Context")
	}
	if rvType.NumOut() != 1 || rvType.Out(0) != reflect.TypeFor[error]() {
		panic("validator factory must return a single error value")
	}

	rvParamTypes := make([]reflect.Type, 0)
	for i := 1; i < rvType.NumIn(); i++ {
		rvParamTypes = append(rvParamTypes, rvType.In(i))
	}

	var newFn func(ctx *schema.Context, params []any) error
	if typedFn, ok := fn.(func(*schema.Context, []any) error); ok {
		newFn = typedFn
	} else {
		newFn = func(ctx *schema.Context, params []any) error {
			defer func() {
				if r := recover(); r != nil {
					panic(fmt.Sprintf("validator factory panic: name=%s, params=%v, err=%v", name, params, r))
				}
			}()

			var err error
			rvParams := make([]reflect.Value, len(rvParamTypes)+1)
			rvParams[0] = reflect.ValueOf(ctx)
			for i, param := range params {
				rvParams[i+1] = reflect.ValueOf(param)
			}

			outs := rv.Call(rvParams)
			if len(outs) != 1 {
				return fmt.Errorf("validator factory must return a single error value")
			}

			out := outs[0].Interface()
			if out != nil {
				err = out.(error)
			}

			if err != nil {
				return &schema.ValidationError{
					Path:   ctx.Path(),
					Name:   name,
					Params: params,
					Err:    err,
				}
			}

			return nil
		}
	}

	r.validators[name] = validatorFactory{
		name:       name,
		paramTypes: rvParamTypes,
		fn:         newFn,
	}
}

// NewValidator gets a field validator by name
// params is a slice of parameter strings
func (r *Registry) NewValidator(name string, params ...any) schema.Validator {
	factory, ok := r.validators[name]
	if !ok {
		panic(fmt.Sprintf("validator '%s' not found in registry", name))
	}

	return factory.Build(params)
}

func (r *Registry) GetValidatorParamTypes(name string) []reflect.Type {
	factory, ok := r.validators[name]
	if !ok {
		panic(fmt.Sprintf("validator '%s' not found in registry", name))
	}

	return factory.paramTypes
}

// DefaultRegistry returns the default registry
func DefaultRegistry() *Registry {
	return defaultRegistry
}

var defaultRegistry = NewRegistry()

func Register(name string, fn any) {
	defaultRegistry.Register(name, fn)
}

func NewValidator(name string, params ...any) schema.Validator {
	return defaultRegistry.NewValidator(name, params...)
}