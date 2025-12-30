package schema

// FieldSchema validates primitive/scalar values
type FieldSchema struct {
	validators []Validator
}

// NewField creates a new field schema
func NewField() *FieldSchema {
	return &FieldSchema{
		validators: make([]Validator, 0),
	}
}

// Validate validates a field value
func (f *FieldSchema) Validate(ctx *Context) error {
	// Run all validators
	for _, validator := range f.validators {
		if ctx.skipRest {
			break
		}

		if err := validator.Validate(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (f *FieldSchema) AddValidator(v Validator) Schema {
	f.validators = append(f.validators, v)
	return f
}

func (f *FieldSchema) RemoveValidator(name string) Schema {
	newValidators := make([]Validator, 0)
	for _, v := range f.validators {
		if v.Name() != name {
			newValidators = append(newValidators, v)
		}
	}
	f.validators = newValidators
	return f
}
