package schema

import (
	"errors"
	"fmt"
)

var ErrCheckFailed = fmt.Errorf("validation check failed")

// ValidationError represents a single validation failure with field path and error code
type ValidationError struct {
	Path string

	// Code is the validation error code (e.g., "required", "min", "eqfield")
	Code string

	// Params contains error parameters (positional)
	Params []any

	Err error
}

func (e ValidationError) Unwrap() error {
	return e.Err
}

// Error implements the error interface
func (e ValidationError) Error() string {
	if len(e.Params) > 0 {
		return fmt.Sprintf("%s: %s %v", e.Path, e.Code, e.Params)
	}
	return fmt.Sprintf("%s: %s", e.Path, e.Code)
}

// ValidationErrors holds all validation errors
type ValidationErrors []ValidationError

func (r ValidationErrors) Unwrap() []error {
	if len(r) == 0 {
		return nil
	}

	errs := make([]error, len(r))
	for i, err := range r {
		errs[i] = err
	}

	return errs
}

func (r ValidationErrors) Error() string {
	errs := r.Unwrap()
	if len(errs) == 0 {
		return ""
	}

	return errors.Join(errs...).Error()
}

// AddError adds a validation error to the result
func (r *ValidationErrors) AddError(err ValidationError) {
	*r = append(*r, err)
}

func (r ValidationErrors) HasFieldError(field string) bool {
	for _, err := range r {
		if err.Path == field {
			return true
		}
	}

	return false
}
