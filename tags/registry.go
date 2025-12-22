package tags

import (
	"strings"

	"github.com/weilence/schema-validator/validation"
)

// ValidatorFactory creates a validator from tag parameters
// params is a slice of parameter strings extracted from the tag
// For example: validate:"between=10,20" -> params = ["10", "20"]
type ValidatorFactory func(params []string) (validation.FieldValidator, error)

// ArrayValidatorFactory creates an array validator from tag parameters
type ArrayValidatorFactory func(params []string) (validation.ArrayValidator, error)

// Registry maps validator names to factory functions
type Registry struct {
	fieldValidators map[string]ValidatorFactory
	arrayValidators map[string]ArrayValidatorFactory
}

var defaultRegistry = NewRegistry()

// NewRegistry creates a new validator registry
func NewRegistry() *Registry {
	r := &Registry{
		fieldValidators: make(map[string]ValidatorFactory),
		arrayValidators: make(map[string]ArrayValidatorFactory),
	}

	// Register built-in validators
	r.registerBuiltins()

	return r
}

func (r *Registry) registerBuiltins() {
	// Field validators
	r.RegisterField("required", func(params []string) (validation.FieldValidator, error) {
		return validation.Required(), nil
	})

	r.RegisterField("min", func(params []string) (validation.FieldValidator, error) {
		if len(params) == 0 {
			return nil, nil
		}
		return validation.Min(parseIntOrString(params[0])), nil
	})

	r.RegisterField("max", func(params []string) (validation.FieldValidator, error) {
		if len(params) == 0 {
			return nil, nil
		}
		return validation.Max(parseIntOrString(params[0])), nil
	})

	r.RegisterField("min_length", func(params []string) (validation.FieldValidator, error) {
		if len(params) == 0 {
			return nil, nil
		}
		return validation.MinLength(parseInt(params[0])), nil
	})

	r.RegisterField("max_length", func(params []string) (validation.FieldValidator, error) {
		if len(params) == 0 {
			return nil, nil
		}
		return validation.MaxLength(parseInt(params[0])), nil
	})

	r.RegisterField("email", func(params []string) (validation.FieldValidator, error) {
		return validation.Email(), nil
	})

	r.RegisterField("url", func(params []string) (validation.FieldValidator, error) {
		return validation.URL(), nil
	})

	// Cross-field validators
	r.RegisterField("eqfield", func(params []string) (validation.FieldValidator, error) {
		if len(params) == 0 {
			return nil, nil
		}
		return validation.EqField(params[0]), nil
	})

	r.RegisterField("nefield", func(params []string) (validation.FieldValidator, error) {
		if len(params) == 0 {
			return nil, nil
		}
		return validation.NeField(params[0]), nil
	})

	r.RegisterField("gtfield", func(params []string) (validation.FieldValidator, error) {
		if len(params) == 0 {
			return nil, nil
		}
		return validation.GtField(params[0]), nil
	})

	r.RegisterField("ltfield", func(params []string) (validation.FieldValidator, error) {
		if len(params) == 0 {
			return nil, nil
		}
		return validation.LtField(params[0]), nil
	})

	// Array validators
	r.RegisterArray("min_items", func(params []string) (validation.ArrayValidator, error) {
		if len(params) == 0 {
			return nil, nil
		}
		return validation.MinItems(parseInt(params[0])), nil
	})

	r.RegisterArray("max_items", func(params []string) (validation.ArrayValidator, error) {
		if len(params) == 0 {
			return nil, nil
		}
		return validation.MaxItems(parseInt(params[0])), nil
	})
}

// RegisterField registers a field validator factory
func (r *Registry) RegisterField(name string, factory ValidatorFactory) {
	r.fieldValidators[name] = factory
}

// RegisterArray registers an array validator factory
func (r *Registry) RegisterArray(name string, factory ArrayValidatorFactory) {
	r.arrayValidators[name] = factory
}

// GetFieldValidator gets a field validator by name
// params is a comma-separated string that will be split into a slice
func (r *Registry) GetFieldValidator(name, params string) (validation.FieldValidator, error) {
	factory, ok := r.fieldValidators[name]
	if !ok {
		return nil, nil // Unknown validator, skip
	}

	// Split params by comma to support multiple parameters
	paramSlice := splitParams(params)
	return factory(paramSlice)
}

// GetArrayValidator gets an array validator by name
func (r *Registry) GetArrayValidator(name, params string) (validation.ArrayValidator, error) {
	factory, ok := r.arrayValidators[name]
	if !ok {
		return nil, nil // Unknown validator, skip
	}

	// Split params by comma to support multiple parameters
	paramSlice := splitParams(params)
	return factory(paramSlice)
}

// DefaultRegistry returns the default registry
func DefaultRegistry() *Registry {
	return defaultRegistry
}

// Helper functions

// splitParams splits a parameter string by colon (for multi-param validators)
// For example: "10:20:30" -> ["10", "20", "30"]
// Single params like "field" -> ["field"]
func splitParams(params string) []string {
	if params == "" {
		return []string{}
	}

	// Check if params contains colon (multi-param)
	if strings.Contains(params, ":") {
		parts := strings.Split(params, ":")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}

	// Single parameter
	return []string{params}
}

func parseInt(s string) int {
	var result int
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			result = result*10 + int(s[i]-'0')
		}
	}
	return result
}

func parseIntOrString(s string) interface{} {
	// Try to parse as int
	isNum := true
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			isNum = false
			break
		}
	}

	if isNum && len(s) > 0 {
		return parseInt(s)
	}

	return s
}
