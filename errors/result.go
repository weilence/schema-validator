package errors

import "strings"

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
		result[err.FieldPath] = append(result[err.FieldPath], err)
	}
	return result
}

// HasFieldError checks if a specific field has an error
func (r *ValidationResult) HasFieldError(fieldPath string) bool {
	for _, err := range r.errors {
		if err.FieldPath == fieldPath {
			return true
		}
	}
	return false
}
