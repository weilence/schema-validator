package validator

import (
	"testing"

	"github.com/weilence/schema-validator/builder"
	"github.com/weilence/schema-validator/schema"
	"github.com/weilence/schema-validator/validators"
)

// Test 1: Tag-based validation
func TestTagBasedValidation(t *testing.T) {
	type User struct {
		Gender   string `json:"gender" validate:"oneof=male,female,other"`
		Email    string `json:"email" validate:"required|email"`
		Password string `json:"password" validate:"required|min=8"`
		Confirm  string `json:"confirm" validate:"required|eqfield=Password"`
		Age      int    `json:"age" validate:"min=18|max=120"`
	}

	v, err := New(User{})
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	// Valid user
	validUser := User{
		Email:    "test@example.com",
		Password: "password123",
		Confirm:  "password123",
		Age:      25,
		Gender:   "male",
	}

	err = v.Validate(validUser)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	// Invalid user - password mismatch
	invalidUser := User{
		Email:    "test@example.com",
		Password: "password123",
		Confirm:  "different",
		Age:      25,
		Gender:   "male",
	}

	err = v.Validate(invalidUser)
	if err == nil {
		t.Error("Expected validation errors for password mismatch")
	} else {
		switch e := err.(type) {
		case *schema.ValidationResult:
			if !e.HasFieldError("confirm") {
				t.Error("Expected error on confirm field")
			}
		case *schema.ValidationError:
			if e.Path != "confirm" {
				t.Errorf("Expected error on confirm field, got %s", e.Path)
			}
		default:
			t.Fatalf("expected ValidationResult or ValidationError, got %T", err)
		}
	}

	// Invalid user - missing required field
	invalidUser2 := User{
		Email:    "",
		Password: "password123",
		Confirm:  "password123",
		Age:      25,
	}

	err = v.Validate(invalidUser2)
	if err == nil {
		t.Error("Expected validation errors for missing email")
	} else {
		switch e := err.(type) {
		case *schema.ValidationResult:
			if !e.HasFieldError("email") {
				t.Error("Expected error on email field")
			}
		case *schema.ValidationError:
			if e.Path != "email" {
				t.Errorf("Expected error on email field, got %s", e.Path)
			}
		default:
			t.Fatalf("expected ValidationResult or ValidationError, got %T", err)
		}
	}
}

// Test 2: Code-based validation
func TestCodeBasedValidation(t *testing.T) {
	userSchema := builder.Object().
		Field("email", builder.Field().AddValidator("required").AddValidator("email").Build()).
		Field("password", builder.Field().AddValidator("required").AddValidator("min", 8).Build()).
		Field("age", builder.Field().AddValidator("min", 18).AddValidator("max", 120).Build()).
		Build()

	v := NewFromSchema(userSchema)

	// Valid data (map)
	validData := map[string]any{
		"email":    "test@example.com",
		"password": "password123",
		"age":      25,
	}

	err := v.Validate(validData)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	// Invalid data - invalid email
	invalidData := map[string]any{
		"email":    "not-an-email",
		"password": "password123",
		"age":      25,
	}

	err = v.Validate(invalidData)
	if err == nil {
		t.Error("Expected validation errors for invalid email")
	} else {
		switch e := err.(type) {
		case *schema.ValidationResult:
			if !e.HasFieldError("email") {
				t.Error("Expected error on email field")
			}
		case *schema.ValidationError:
			if e.Path != "email" {
				t.Errorf("Expected error on email field, got %s", e.Path)
			}
		default:
			t.Fatalf("expected ValidationResult or ValidationError, got %T", err)
		}
	}
}

// Test 3: Embedded struct with private fields
func TestEmbeddedStructWithPrivateFields(t *testing.T) {
	type Address struct {
		Street  string `json:"street" validate:"required"`
		private string // Private field in embedded struct
	}

	type Person struct {
		Name string `json:"name" validate:"required"`
		Address
	}

	v, err := New(Person{})
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	// Valid person
	validPerson := Person{
		Name: "John Doe",
		Address: Address{
			Street:  "123 Main St",
			private: "should be accessible",
		},
	}

	err = v.Validate(validPerson)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	// Invalid person - missing embedded field
	invalidPerson := Person{
		Name: "John Doe",
		Address: Address{
			Street: "", // Required field is empty
		},
	}

	err = v.Validate(invalidPerson)
	if err == nil {
		t.Error("Expected validation errors for missing street")
	} else {
		switch e := err.(type) {
		case *schema.ValidationResult:
			if !e.HasFieldError("street") {
				t.Error("Expected error on street field")
			}
		case *schema.ValidationError:
			if e.Path != "street" {
				t.Errorf("Expected error on street field, got %s", e.Path)
			}
		default:
			t.Fatalf("expected ValidationResult or ValidationError, got %T", err)
		}
	}
}

// Test 4: Array validation
func TestArrayValidation(t *testing.T) {
	type TodoList struct {
		Items []string `json:"items" validate:"min=1|max=10"`
	}

	v, err := New(TodoList{})
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	// Valid todo list
	validList := TodoList{
		Items: []string{"item1", "item2", "item3"},
	}

	err = v.Validate(validList)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	// Invalid - empty array
	invalidList := TodoList{
		Items: []string{},
	}

	err = v.Validate(invalidList)
	if err == nil {
		t.Error("Expected validation errors for empty array")
	} else {
		switch e := err.(type) {
		case *schema.ValidationResult:
			if !e.HasFieldError("items") {
				t.Error("Expected error on items field")
			}
		case *schema.ValidationError:
			if e.Path != "items" {
				t.Errorf("Expected error on items field, got %s", e.Path)
			}
		default:
			t.Fatalf("expected ValidationResult or ValidationError, got %T", err)
		}
	}
}

// Test 5: Cross-field validation with code
func TestCrossFieldValidationWithCode(t *testing.T) {
	// register passwordMatchValidator into registry and add by name
	validators.Register("password", func(ctx *schema.Context, params []any) error {
		return nil
	})

	userSchema := builder.Object().
		Field("password", builder.Field().AddValidator("required").AddValidator("min", 8).Build()).
		Field("confirmPassword", builder.Field().AddValidator("required").AddValidator("eqfield", "password").Build()).
		AddValidator("password").
		Build()

	v := NewFromSchema(userSchema)

	// Valid data
	validData := map[string]any{
		"password":        "password123",
		"confirmPassword": "password123",
	}

	err := v.Validate(validData)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	// Invalid - password mismatch
	invalidData := map[string]any{
		"password":        "password123",
		"confirmPassword": "different",
	}

	err = v.Validate(invalidData)
	if err == nil {
		t.Error("Expected validation errors for password mismatch")
	}
}

// Test 6: Map validation
func TestMapValidation(t *testing.T) {
	userSchema := builder.Object().
		Field("name", builder.Field().AddValidator("required").Build()).
		Field("age", builder.Field().AddValidator("min", 0).Build()).
		Build()

	v := NewFromSchema(userSchema)

	// Valid map
	validMap := map[string]any{
		"name": "John Doe",
		"age":  30,
	}

	err := v.Validate(validMap)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	// Invalid map
	invalidMap := map[string]any{
		"name": "",
		"age":  30,
	}

	err = v.Validate(invalidMap)
	if err == nil {
		t.Error("Expected validation errors for empty name")
	}
}

// DynamicForm implements SchemaModifier to dynamically modify validation rules
type DynamicForm struct {
	Type     string `json:"type" validate:"required"`
	Value    string `json:"value"`
	Required bool   `json:"required"`
}

func (f DynamicForm) ModifySchema(ctx *schema.Context) {
	// Access current object's fields
	requiredField, _ := ctx.GetValue("Required")
	if requiredField != nil {
		isRequired, _ := requiredField.Bool()

		// Dynamically modify "value" field validation based on required flag
		if isRequired {
			// Add required validation
			valueSchema := builder.Field().
				AddValidator("required").
				Build()
			ctx.Schema().(*schema.ObjectSchema).AddField("value", valueSchema)
		} else {
			// Make value optional (override any existing schema)
			ctx.Schema().(*schema.ObjectSchema).RemoveField("value")
		}
	}

	// Can also access parent context
	if ctx.Parent() != nil {
		// Access parent object fields if needed
	}
}

// Test SchemaModifier interface
func TestSchemaModifier(t *testing.T) {
	v, err := New(DynamicForm{})
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	// Test case 1: required=true, value should be required
	form1 := DynamicForm{
		Type:     "text",
		Value:    "", // Empty value
		Required: true,
	}

	err = v.Validate(form1)
	if err == nil {
		t.Error("Expected validation to fail when required=true and value is empty")
	} else {
		switch e := err.(type) {
		case *schema.ValidationResult:
			if !e.HasFieldError("value") {
				t.Error("Expected error on value field when required=true")
			}
		case *schema.ValidationError:
			if e.Path != "value" {
				t.Errorf("Expected error on value field, got %s", e.Path)
			}
		default:
			t.Fatalf("expected ValidationResult or ValidationError, got %T", err)
		}
	}

	// Test case 2: required=false, value can be empty
	form2 := DynamicForm{
		Type:     "text",
		Value:    "", // Empty value
		Required: false,
	}

	err = v.Validate(form2)
	if err != nil {
		t.Errorf("Expected validation to pass when required=false, got errors: %v", err)
	}

	// Test case 3: required=true, value is provided
	form3 := DynamicForm{
		Type:     "text",
		Value:    "some value",
		Required: true,
	}

	err = v.Validate(form3)
	if err != nil {
		t.Errorf("Expected validation to pass when required=true and value is provided, got errors: %v", err)
	}
}

// NestedUser implements SchemaModifier to add zip code validation based on country
type NestedUser struct {
	Name string `json:"name" validate:"required"`
	NestedAddress
}

func (u NestedUser) ModifySchema(ctx *schema.Context) {
	// For embedded structs, access fields directly
	countryField, _ := ctx.GetValue("Country")

	country := countryField.String()

	// Add different validation rules based on country
	if country == "US111" {
		// US zip codes should be 5 digits
		zipCodeSchema := builder.Field().
			AddValidator("min", 5).
			AddValidator("max", 5).
			Build()
		ctx.Schema().(*schema.ObjectSchema).AddField("zipCode", zipCodeSchema)
	}
}

type NestedAddress struct {
	Country string `json:"country"`
	ZipCode string `json:"zipCode"`
}

func (u NestedAddress) ModifySchema(ctx *schema.Context) {
	// For embedded structs, access fields directly
	countryField, _ := ctx.GetValue("Country")

	country := countryField.String()

	// Add different validation rules based on country
	if country == "US" {
		// US zip codes should be 5 digits
		zipCodeSchema := builder.Field().
			AddValidator("min", 5).
			AddValidator("max", 5).
			Build()
		ctx.Schema().(*schema.ObjectSchema).AddField("country", zipCodeSchema)
	}
}

// Test accessing nested values in SchemaModifier
func TestSchemaModifierNestedAccess(t *testing.T) {
	v, err := New(NestedUser{})
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	// Test US address with invalid zip code
	user := NestedUser{
		Name: "John",
		NestedAddress: NestedAddress{
			Country: "US111",
			ZipCode: "123", // Too short
		},
	}

	err = v.Validate(user)
	if err == nil {
		t.Error("Expected validation to fail for US zip code with length < 5")
	}

	user = NestedUser{
		Name: "John",
		NestedAddress: NestedAddress{
			Country: "US", // Too short
			ZipCode: "12345",
		},
	}

	err = v.Validate(user)
	if err == nil {
		t.Error("Expected validation to fail for US zip code with length < 5")
	}
}
