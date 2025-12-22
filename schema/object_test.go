package schema_test

import (
	"testing"

	validator "github.com/weilence/schema-validator"
	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/schema"
	"github.com/weilence/schema-validator/validation"
)

// DynamicForm implements SchemaModifier to dynamically modify validation rules
type DynamicForm struct {
	Type     string `json:"type" validate:"required"`
	Value    string `json:"value"`
	Required bool   `json:"required"`
}

func (f DynamicForm) ModifySchema(ctx *validation.Context, accessor data.ObjectAccessor, s *schema.ObjectSchema) {
	// Access current object's fields
	requiredField, _ := accessor.GetField("required")
	if requiredField != nil {
		fieldAcc, _ := requiredField.AsField()
		isRequired, _ := fieldAcc.Bool()

		// Dynamically modify "value" field validation based on required flag
		if isRequired {
			// Add required validation
			valueSchema := schema.Field().
				AddValidator(validation.Required()).
				Build()
			s.AddField("value", valueSchema)
		} else {
			// Make value optional (override any existing schema)
			valueSchema := schema.Field().
				SetOptional(true).
				Build()
			s.AddField("value", valueSchema)
		}
	}

	// Can also access parent context
	if ctx.Parent != nil {
		// Access parent object fields if needed
	}

	// Can access root
	if ctx.Root != nil {
		// Access root object if needed
	}
}

// Test SchemaModifier interface
func TestSchemaModifier(t *testing.T) {
	v, err := validator.NewFromStruct(DynamicForm{})
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	// Test case 1: required=true, value should be required
	form1 := DynamicForm{
		Type:     "text",
		Value:    "", // Empty value
		Required: true,
	}

	result, _ := v.Validate(form1)
	if result.IsValid() {
		t.Error("Expected validation to fail when required=true and value is empty")
	}
	if !result.HasFieldError("value") {
		t.Error("Expected error on value field when required=true")
	}

	// Test case 2: required=false, value can be empty
	form2 := DynamicForm{
		Type:     "text",
		Value:    "", // Empty value
		Required: false,
	}

	result, _ = v.Validate(form2)
	if !result.IsValid() {
		t.Errorf("Expected validation to pass when required=false, got errors: %v", result.Errors())
	}

	// Test case 3: required=true, value is provided
	form3 := DynamicForm{
		Type:     "text",
		Value:    "some value",
		Required: true,
	}

	result, _ = v.Validate(form3)
	if !result.IsValid() {
		t.Errorf("Expected validation to pass when required=true and value is provided, got errors: %v", result.Errors())
	}
}

// NestedUser implements SchemaModifier to add zip code validation based on country
type NestedUser struct {
	Name string `json:"name" validate:"required"`
	NestedAddress
}

type NestedAddress struct {
	Country string `json:"country"`
	ZipCode string `json:"zipCode"`
}

func (u NestedUser) ModifySchema(ctx *validation.Context, accessor data.ObjectAccessor, s *schema.ObjectSchema) {
	// For embedded structs, access fields directly
	countryField, exists := accessor.GetField("country")
	if !exists {
		return
	}

	countryAcc, _ := countryField.AsField()
	country := countryAcc.String()

	// Add different validation rules based on country
	if country == "US" {
		// US zip codes should be 5 digits
		zipCodeSchema := schema.Field().
			AddValidator(validation.MinLength(5)).
			AddValidator(validation.MaxLength(5)).
			Build()
		s.AddField("zipCode", zipCodeSchema)
	}
}

// Test accessing nested values in SchemaModifier
func TestSchemaModifierNestedAccess(t *testing.T) {
	v, err := validator.NewFromStruct(NestedUser{})
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	// Test US address with invalid zip code
	user := NestedUser{
		Name: "John",
		NestedAddress: NestedAddress{
			Country: "US",
			ZipCode: "123", // Too short
		},
	}

	result, _ := v.Validate(user)
	if result.IsValid() {
		t.Error("Expected validation to fail for US zip code with length < 5")
	}
}
