package tags_test

import (
	"reflect"
	"testing"

	validator "github.com/weilence/schema-validator"
	"github.com/weilence/schema-validator/errors"
	"github.com/weilence/schema-validator/tags"
	"github.com/weilence/schema-validator/validation"
)

// Test multi-parameter validator factory
func TestMultiParameterValidator(t *testing.T) {
	// Register a custom validator that takes multiple parameters
	registry := tags.NewRegistry()

	// Example: between=10,20 validator
	registry.RegisterField("between", func(params []string) (validation.FieldValidator, error) {
		if len(params) < 2 {
			return nil, nil
		}

		min := parseIntHelper(params[0])
		max := parseIntHelper(params[1])

		return validation.FieldValidatorFunc(func(ctx *validation.Context) error {
			field, err := ctx.AsField()
			if err != nil {
				return nil
			}

			val, err := field.Int()
			if err != nil {
				return nil
			}

			if val < int64(min) || val > int64(max) {
				return errors.NewValidationError(ctx.Path(), "between", map[string]interface{}{
					"min":    min,
					"max":    max,
					"actual": val,
				})
			}
			return nil
		}), nil
	})

	type Product struct {
		Price int `json:"price" validate:"between=10:100"`
	}

	// Parse with custom registry
	typ := reflect.TypeOf(Product{})
	objSchema, err := tags.ParseStructTagsWithRegistry(typ, registry)
	if err != nil {
		t.Fatalf("Failed to parse tags: %v", err)
	}

	v := validator.New(objSchema)

	// Test valid price
	product1 := Product{Price: 50}
	result, _ := v.Validate(product1)
	if !result.IsValid() {
		t.Errorf("Expected valid price, got errors: %v", result.Errors())
	}

	// Test price too low
	product2 := Product{Price: 5}
	result, _ = v.Validate(product2)
	if result.IsValid() {
		t.Error("Expected validation to fail for price < 10")
	}
	if !result.HasFieldError("price") {
		t.Error("Expected error on price field")
	}

	// Test price too high
	product3 := Product{Price: 150}
	result, _ = v.Validate(product3)
	if result.IsValid() {
		t.Error("Expected validation to fail for price > 100")
	}
}

// Test multi-parameter with three params
func TestThreeParameterValidator(t *testing.T) {
	registry := tags.NewRegistry()

	// Example: enum=option1:option2:option3 validator
	registry.RegisterField("enum", func(params []string) (validation.FieldValidator, error) {
		if len(params) == 0 {
			return nil, nil
		}

		allowedValues := make(map[string]bool)
		for _, p := range params {
			allowedValues[p] = true
		}

		return validation.FieldValidatorFunc(func(ctx *validation.Context) error {
			field, _ := ctx.AsField()
			val := field.String()
			if !allowedValues[val] {
				return errors.NewValidationError(ctx.Path(), "enum", map[string]interface{}{
					"allowed": params,
					"actual":  val,
				})
			}
			return nil
		}), nil
	})

	type Settings struct {
		Theme string `json:"theme" validate:"enum=light:dark:auto"`
	}

	typ := reflect.TypeOf(Settings{})
	objSchema, err := tags.ParseStructTagsWithRegistry(typ, registry)
	if err != nil {
		t.Fatalf("Failed to parse tags: %v", err)
	}

	v := validator.New(objSchema)

	// Test valid theme
	settings1 := Settings{Theme: "dark"}
	result, _ := v.Validate(settings1)
	if !result.IsValid() {
		t.Errorf("Expected valid theme, got errors: %v", result.Errors())
	}

	// Test invalid theme
	settings2 := Settings{Theme: "blue"}
	result, _ = v.Validate(settings2)
	if result.IsValid() {
		t.Error("Expected validation to fail for invalid theme")
	}
	if !result.HasFieldError("theme") {
		t.Error("Expected error on theme field")
	}
}

// Helper function
func parseIntHelper(s string) int {
	var result int
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			result = result*10 + int(s[i]-'0')
		}
	}
	return result
}
