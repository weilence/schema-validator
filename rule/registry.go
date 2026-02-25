package rule

import (
	"errors"
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

// Registry maps validator names to factory functions.
//
// Thread-safety: This type is NOT safe for concurrent use.
// The Registry is expected to be configured during initialization
// (via Register/Alias) and then used read-only (via NewValidator).
// If you need concurrent registration, use external synchronization.
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
func (r *Registry) Register(code string, fn any) error {
	rv := reflect.ValueOf(fn)
	rvType := rv.Type()
	if rvType.Kind() != reflect.Func {
		return fmt.Errorf("validator factory must be a function")
	}
	if rvType.NumIn() < 1 || rvType.In(0) != reflect.TypeFor[*schema.Context]() {
		return fmt.Errorf("first parameter of validator factory must be *schema.Context")
	}
	if rvType.NumOut() != 1 || rvType.Out(0) != reflect.TypeFor[error]() {
		return fmt.Errorf("validator factory must return a single error value")
	}

	rvParamTypes := make([]reflect.Type, 0)
	for i := 1; i < rvType.NumIn(); i++ {
		rvParamTypes = append(rvParamTypes, rvType.In(i))
	}

	var newFn func(ctx *schema.Context, params []any) error
	if typedFn, ok := fn.(func(*schema.Context, []any) error); ok {
		newFn = typedFn
	} else {
		newFn = func(ctx *schema.Context, params []any) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("validator factory panic: name=%s, params=%v, err=%v", code, params, r)
				}
			}()

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
				return out.(error)
			}

			return nil
		}
	}

	newFn2 := func(ctx *schema.Context, params []any) error {
		err := newFn(ctx, params)
		if err != nil {
			newErr := schema.ValidationError{
				Path:   ctx.Path(),
				Code:   code,
				Params: params,
				Err:    err,
			}

			if errors.Is(err, schema.ErrCheckFailed) {
				ctx.AddError(newErr)
			} else {
				return newErr
			}
		}

		return nil
	}

	r.validators[code] = validatorFactory{
		name:       code,
		paramTypes: rvParamTypes,
		fn:         newFn2,
	}

	return nil
}

func (r *Registry) Alias(oldName, newName string) error {
	factory, ok := r.validators[oldName]
	if !ok {
		return fmt.Errorf("validator '%s' not found in registry", oldName)
	}

	r.validators[newName] = factory

	return nil
}

// NewValidator gets a field validator by name, returns error if not found
// params is a slice of parameter strings
func (r *Registry) NewValidator(name string, params ...any) (schema.Validator, error) {
	factory, ok := r.validators[name]

	if !ok {
		return nil, fmt.Errorf("validator '%s' not found in registry", name)
	}

	return factory.Build(params), nil
}

// GetValidatorParamTypes gets validator param types, returns error if not found
func (r *Registry) GetValidatorParamTypes(name string) ([]reflect.Type, error) {
	factory, ok := r.validators[name]

	if !ok {
		return nil, fmt.Errorf("validator '%s' not found in registry", name)
	}

	return factory.paramTypes, nil
}

// DefaultRegistry returns the default registry
func DefaultRegistry() *Registry {
	return defaultRegistry
}

var defaultRegistry = NewRegistry()

// Register registers a validator to the default registry, returns error
func Register(name string, fn any) error {
	return defaultRegistry.Register(name, fn)
}

// NewValidator creates a validator from the default registry, returns error
func NewValidator(name string, params ...any) (schema.Validator, error) {
	return defaultRegistry.NewValidator(name, params...)
}
