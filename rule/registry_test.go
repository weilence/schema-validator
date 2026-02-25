package rule

import (
	"fmt"
	"sync"
	"testing"

	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/schema"
)

// mustNewValidator is a test helper that creates a validator and panics on error
func mustNewValidator(r *Registry, name string, params ...any) schema.Validator {
	v, err := r.NewValidator(name, params...)
	if err != nil {
		panic(fmt.Sprintf("NewValidator(%q, %v) failed: %v", name, params, err))
	}
	return v
}

// Test multi-parameter validator factory
func TestMultiParameterValidator(t *testing.T) {
	// Register a custom validator that takes multiple parameters
	// Example: between=10,20 validator
	Register("between", func(ctx *schema.Context, min any, max any) error {
		field := ctx.Value()

		minValue := data.NewValue(min)
		maxValue := data.NewValue(max)

		ok1, err := compareValue(GreaterThanOrEqual, field, minValue)
		if err != nil {
			return err
		}

		ok2, err := compareValue(LessThanOrEqual, field, maxValue)
		if err != nil {
			return err
		}

		if !ok1 || !ok2 {
			return schema.ErrCheckFailed
		}
		return nil
	})

	// Create a context with a field value
	s := schema.NewField().AddValidator(mustNewValidator(DefaultRegistry(), "between", 10, 20))
	ctx := schema.NewContext(s, data.NewValue(15))
	// Validate the context
	if err := s.Validate(ctx); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Test multi-parameter with three params
func TestThreeParameterValidator(t *testing.T) {
	// Example: enum=option1,option2,option3 validator
	Register("enum", func(ctx *schema.Context, params []string) error {
		if len(params) == 0 {
			return nil
		}

		allowedValues := make(map[string]bool)
		for _, p := range params {
			allowedValues[p] = true
		}

		field := ctx.Value()
		val := field.String()
		if !allowedValues[val] {
			return schema.ErrCheckFailed
		}
		return nil
	})
}

// TestRegistryConcurrentRead tests concurrent reading of validators
func TestRegistryConcurrentRead(t *testing.T) {
	r := NewRegistry()

	// Register some validators first
	for i := range 10 {
		name := fmt.Sprintf("validator_%d", i)
		err := r.Register(name, func(ctx *schema.Context) error {
			return schema.ErrCheckFailed
		})
		if err != nil {
			t.Fatalf("Register failed: %v", err)
		}
	}

	var wg sync.WaitGroup
	numGoroutines := 100
	readsPerGoroutine := 100

	for range numGoroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range readsPerGoroutine {
				for k := range 10 {
					name := fmt.Sprintf("validator_%d", k)
					v, err := r.NewValidator(name)
					if err != nil {
						t.Errorf("NewValidator failed: %v", err)
					}
					if v == nil {
						t.Errorf("NewValidator returned nil for %s", name)
					}
				}
			}
		}()
	}

	wg.Wait()
}

// TestDefaultRegistryConcurrentAccess tests concurrent access to the default registry
func TestDefaultRegistryConcurrentAccess(t *testing.T) {
	var wg sync.WaitGroup
	numGoroutines := 50

	for range numGoroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Access default registry concurrently
			dr := DefaultRegistry()
			if dr == nil {
				t.Error("DefaultRegistry returned nil")
			}
		}()
	}

	wg.Wait()
}
