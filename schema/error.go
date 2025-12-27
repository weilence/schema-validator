package schema

import (
	"fmt"
	"strings"
)

// ValidationError represents a single validation failure with field path and error code
type ValidationError struct {
	// Path is the path to the field (e.g., "user.email", "items[0].name")
	Path string

	// Name is the validation error code (e.g., "required", "min", "eqfield")
	Name string

	// Params contains additional error parameters (e.g., min value, field name)
	Params []any

	Err error
}

func NewValidationError(path, name string, params map[string]any) *ValidationError {
	paramList := make([]any, 0, len(params))
	for k, v := range params {
		paramList = append(paramList, fmt.Sprintf("%s=%v", k, v))
	}
	return &ValidationError{
		Path:   path,
		Name:   name,
		Params: paramList,
	}
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	if len(e.Params) > 0 {
		return fmt.Sprintf("%s: %s %v", e.Path, e.Name, e.Params)
	}
	return fmt.Sprintf("%s: %s", e.Path, e.Name)
}

// ValidationResult holds all validation errors
type ValidationResult struct {
	errors []*ValidationError
}

// NewValidationResult creates a new validation result
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		errors: make([]*ValidationError, 0),
	}
}

func (r *ValidationResult) Error() string {
	var sb strings.Builder

	for i, err := range r.errors {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(err.Error())
	}

	return sb.String()
}

// AddError adds a validation error to the result
func (r *ValidationResult) AddError(err *ValidationError) {
	r.errors = append(r.errors, err)
}

// Errors returns all validation errors
func (r *ValidationResult) Errors() []*ValidationError {
	return r.errors
}

// IsValid returns true if there are no validation errors
func (r *ValidationResult) IsValid() bool {
	return len(r.errors) == 0
}

// FirstError returns the first validation error, or nil if there are none
func (r *ValidationResult) FirstError() *ValidationError {
	if len(r.errors) > 0 {
		return r.errors[0]
	}
	return nil
}

// ErrorsByField groups errors by field path
func (r *ValidationResult) ErrorsByField() map[string][]*ValidationError {
	result := make(map[string][]*ValidationError)
	for _, err := range r.errors {
		result[err.Path] = append(result[err.Path], err)
	}
	return result
}

// HasFieldError checks if a specific field has an error
func (r *ValidationResult) HasFieldError(fieldPath string) bool {
	for _, err := range r.errors {
		if err.Path == fieldPath {
			return true
		}
	}
	return false
}
