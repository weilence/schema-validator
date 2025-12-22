package validator_test

import (
	"testing"

	validator "github.com/weilence/schema-validator"
	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/schema"
	"github.com/weilence/schema-validator/validation"
)

// Test 1: Tag-based validation
func TestTagBasedValidation(t *testing.T) {
	type User struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min_length=8"`
		Confirm  string `json:"confirm" validate:"required,eqfield=password"`
		Age      int    `json:"age" validate:"min=18,max=120"`
	}

	v, err := validator.NewFromStruct(User{})
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	// Valid user
	validUser := User{
		Email:    "test@example.com",
		Password: "password123",
		Confirm:  "password123",
		Age:      25,
	}

	result, err := v.Validate(validUser)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if !result.IsValid() {
		t.Errorf("Expected valid user, got errors: %v", result.Errors())
	}

	// Invalid user - password mismatch
	invalidUser := User{
		Email:    "test@example.com",
		Password: "password123",
		Confirm:  "different",
		Age:      25,
	}

	result, err = v.Validate(invalidUser)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if result.IsValid() {
		t.Error("Expected validation errors for password mismatch")
	}

	if !result.HasFieldError("confirm") {
		t.Error("Expected error on confirm field")
	}

	// Invalid user - missing required field
	invalidUser2 := User{
		Email:    "",
		Password: "password123",
		Confirm:  "password123",
		Age:      25,
	}

	result, err = v.Validate(invalidUser2)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if result.IsValid() {
		t.Error("Expected validation errors for missing email")
	}

	if !result.HasFieldError("email") {
		t.Error("Expected error on email field")
	}
}

// Test 2: Code-based validation
func TestCodeBasedValidation(t *testing.T) {
	userSchema := schema.Object().
		Field("email", schema.Field().AddValidator(validation.Required()).AddValidator(validation.Email()).Build()).
		Field("password", schema.Field().AddValidator(validation.Required()).AddValidator(validation.MinLength(8)).Build()).
		Field("age", schema.Field().AddValidator(validation.Min(18)).AddValidator(validation.Max(120)).Build()).
		Build()

	v := validator.New(userSchema)

	// Valid data (map)
	validData := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
		"age":      25,
	}

	result, err := v.Validate(validData)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if !result.IsValid() {
		t.Errorf("Expected valid data, got errors: %v", result.Errors())
	}

	// Invalid data - invalid email
	invalidData := map[string]interface{}{
		"email":    "not-an-email",
		"password": "password123",
		"age":      25,
	}

	result, err = v.Validate(invalidData)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if result.IsValid() {
		t.Error("Expected validation errors for invalid email")
	}

	if !result.HasFieldError("email") {
		t.Error("Expected error on email field")
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

	v, err := validator.NewFromStruct(Person{})
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

	result, err := v.Validate(validPerson)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if !result.IsValid() {
		t.Errorf("Expected valid person, got errors: %v", result.Errors())
	}

	// Invalid person - missing embedded field
	invalidPerson := Person{
		Name: "John Doe",
		Address: Address{
			Street: "", // Required field is empty
		},
	}

	result, err = v.Validate(invalidPerson)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if result.IsValid() {
		t.Error("Expected validation errors for missing street")
	}

	if !result.HasFieldError("street") {
		t.Error("Expected error on street field")
	}
}

// Test 4: Array validation
func TestArrayValidation(t *testing.T) {
	type TodoList struct {
		Items []string `json:"items" validate:"min_items=1,max_items=10"`
	}

	v, err := validator.NewFromStruct(TodoList{})
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	// Valid todo list
	validList := TodoList{
		Items: []string{"item1", "item2", "item3"},
	}

	result, err := v.Validate(validList)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if !result.IsValid() {
		t.Errorf("Expected valid list, got errors: %v", result.Errors())
	}

	// Invalid - empty array
	invalidList := TodoList{
		Items: []string{},
	}

	result, err = v.Validate(invalidList)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if result.IsValid() {
		t.Error("Expected validation errors for empty array")
	}

	if !result.HasFieldError("items") {
		t.Error("Expected error on items field")
	}
}

// Test 5: Cross-field validation with code
func TestCrossFieldValidationWithCode(t *testing.T) {
	// Custom password match validator
	passwordMatchValidator := validation.ObjectValidatorFunc(func(ctx *validation.Context, obj data.ObjectAccessor) error {
		// This would be implemented properly in a real scenario
		return nil
	})

	userSchema := schema.Object().
		Field("password", schema.Field().AddValidator(validation.Required()).AddValidator(validation.MinLength(8)).Build()).
		Field("confirmPassword", schema.Field().AddValidator(validation.Required()).AddValidator(validation.EqField("password")).Build()).
		CrossField(passwordMatchValidator).
		Build()

	v := validator.New(userSchema)

	// Valid data
	validData := map[string]interface{}{
		"password":        "password123",
		"confirmPassword": "password123",
	}

	result, err := v.Validate(validData)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if !result.IsValid() {
		t.Errorf("Expected valid data, got errors: %v", result.Errors())
	}

	// Invalid - password mismatch
	invalidData := map[string]interface{}{
		"password":        "password123",
		"confirmPassword": "different",
	}

	result, err = v.Validate(invalidData)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if result.IsValid() {
		t.Error("Expected validation errors for password mismatch")
	}
}

// Test 6: Map validation
func TestMapValidation(t *testing.T) {
	userSchema := schema.Object().
		Field("name", schema.Field().AddValidator(validation.Required()).Build()).
		Field("age", schema.Field().AddValidator(validation.Min(0)).Build()).
		Build()

	v := validator.New(userSchema)

	// Valid map
	validMap := map[string]interface{}{
		"name": "John Doe",
		"age":  30,
	}

	result, err := v.Validate(validMap)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if !result.IsValid() {
		t.Errorf("Expected valid map, got errors: %v", result.Errors())
	}

	// Invalid map
	invalidMap := map[string]interface{}{
		"name": "",
		"age":  30,
	}

	result, err = v.Validate(invalidMap)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if result.IsValid() {
		t.Error("Expected validation errors for empty name")
	}
}
