package errors

import "fmt"

// ValidationError represents a single validation failure with field path and error code
type ValidationError struct {
	// FieldPath is the path to the field (e.g., "user.email", "items[0].name")
	FieldPath string

	// ErrorCode is the validation error code (e.g., "required", "min", "eqfield")
	ErrorCode string

	// Params contains additional error parameters (e.g., min value, field name)
	Params map[string]interface{}
}

// NewValidationError creates a new validation error
func NewValidationError(path, code string, params map[string]interface{}) *ValidationError {
	if params == nil {
		params = make(map[string]interface{})
	}
	return &ValidationError{
		FieldPath: path,
		ErrorCode: code,
		Params:    params,
	}
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	if len(e.Params) > 0 {
		return fmt.Sprintf("%s: %s %v", e.FieldPath, e.ErrorCode, e.Params)
	}
	return fmt.Sprintf("%s: %s", e.FieldPath, e.ErrorCode)
}
