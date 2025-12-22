package main

import (
	"fmt"

	validator "github.com/weilence/schema-validator"
	"github.com/weilence/schema-validator/schema"
	"github.com/weilence/schema-validator/validation"
)

// Example 1: Tag-based validation
func example1() {
	fmt.Println("=== Example 1: Tag-based Validation ===")

	type User struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min_length=8"`
		Confirm  string `json:"confirm" validate:"required,eqfield=Password"`
		Age      int    `json:"age" validate:"min=18,max=120"`
	}

	v, err := validator.NewFromStruct(User{})
	if err != nil {
		panic(err)
	}

	// Valid user
	validUser := User{
		Email:    "john@example.com",
		Password: "securepass123",
		Confirm:  "securepass123",
		Age:      30,
	}

	result, _ := v.Validate(validUser)
	fmt.Printf("Valid user result: %v\n", result.IsValid())

	// Invalid user
	invalidUser := User{
		Email:    "not-an-email",
		Password: "short",
		Confirm:  "different",
		Age:      15,
	}

	result, _ = v.Validate(invalidUser)
	fmt.Printf("Invalid user result: %v\n", result.IsValid())
	for _, err := range result.Errors() {
		fmt.Printf("  - %s: %s\n", err.FieldPath, err.ErrorCode)
	}

	fmt.Println()
}

// Example 2: Code-based validation with maps
func example2() {
	fmt.Println("=== Example 2: Code-based Validation ===")

	userSchema := schema.Object().
		Field("email", schema.Field().
			AddValidator(validation.Required()).
			AddValidator(validation.Email()).
			Build()).
		Field("password", schema.Field().
			AddValidator(validation.Required()).
			AddValidator(validation.MinLength(8)).
			Build()).
		Field("age", schema.Field().
			AddValidator(validation.Min(18)).
			AddValidator(validation.Max(120)).
			Build()).
		Build()

	v := validator.New(userSchema)

	// Valid map data
	validData := map[string]interface{}{
		"email":    "jane@example.com",
		"password": "password123",
		"age":      25,
	}

	result, _ := v.Validate(validData)
	fmt.Printf("Valid map result: %v\n", result.IsValid())

	// Invalid map data
	invalidData := map[string]interface{}{
		"email":    "invalid-email",
		"password": "short",
		"age":      150,
	}

	result, _ = v.Validate(invalidData)
	fmt.Printf("Invalid map result: %v\n", result.IsValid())
	for _, err := range result.Errors() {
		fmt.Printf("  - %s: %s (params: %v)\n", err.FieldPath, err.ErrorCode, err.Params)
	}

	fmt.Println()
}

// Example 3: Embedded struct with private fields
func example3() {
	fmt.Println("=== Example 3: Embedded Struct with Private Fields ===")

	type Address struct {
		Street  string `json:"street" validate:"required"`
		City    string `json:"city" validate:"required"`
		private string // Private field in embedded struct
	}

	type Person struct {
		Name string `json:"name" validate:"required,min_length=2"`
		Age  int    `json:"age" validate:"min=0,max=150"`
		Address
	}

	v, _ := validator.NewFromStruct(Person{})

	// Valid person
	validPerson := Person{
		Name: "Alice",
		Age:  28,
		Address: Address{
			Street:  "123 Main St",
			City:    "New York",
			private: "this is accessible when embedded",
		},
	}

	result, _ := v.Validate(validPerson)
	fmt.Printf("Valid person result: %v\n", result.IsValid())

	// Invalid person - missing embedded field
	invalidPerson := Person{
		Name: "B", // Too short
		Age:  28,
		Address: Address{
			Street: "", // Required but empty
			City:   "New York",
		},
	}

	result, _ = v.Validate(invalidPerson)
	fmt.Printf("Invalid person result: %v\n", result.IsValid())
	for _, err := range result.Errors() {
		fmt.Printf("  - %s: %s\n", err.FieldPath, err.ErrorCode)
	}

	fmt.Println()
}

// Example 4: Array validation
func example4() {
	fmt.Println("=== Example 4: Array Validation ===")

	itemSchema := schema.Field().AddValidator(validation.MinLength(1)).Build()

	arraySchema := schema.Array(itemSchema).
		MinItems(1).
		MaxItems(5).
		Build()

	listSchema := schema.Object().
		Field("items", arraySchema).
		Build()

	v := validator.New(listSchema)

	// Valid list
	validList := map[string]interface{}{
		"items": []string{"item1", "item2", "item3"},
	}

	result, _ := v.Validate(validList)
	fmt.Printf("Valid list result: %v\n", result.IsValid())

	// Invalid - empty array
	invalidList := map[string]interface{}{
		"items": []string{},
	}

	result, _ = v.Validate(invalidList)
	fmt.Printf("Invalid list result (empty): %v\n", result.IsValid())
	for _, err := range result.Errors() {
		fmt.Printf("  - %s: %s\n", err.FieldPath, err.ErrorCode)
	}

	// Invalid - too many items
	tooManyItems := map[string]interface{}{
		"items": []string{"1", "2", "3", "4", "5", "6"},
	}

	result, _ = v.Validate(tooManyItems)
	fmt.Printf("Invalid list result (too many): %v\n", result.IsValid())
	for _, err := range result.Errors() {
		fmt.Printf("  - %s: %s (params: %v)\n", err.FieldPath, err.ErrorCode, err.Params)
	}

	fmt.Println()
}

// Example 5: Cross-field validation
func example5() {
	fmt.Println("=== Example 5: Cross-field Validation ===")

	type PasswordForm struct {
		Password        string `json:"password" validate:"required,min_length=8"`
		ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
	}

	v, _ := validator.NewFromStruct(PasswordForm{})

	// Valid - passwords match
	validForm := PasswordForm{
		Password:        "securepass123",
		ConfirmPassword: "securepass123",
	}

	result, _ := v.Validate(validForm)
	fmt.Printf("Matching passwords result: %v\n", result.IsValid())

	// Invalid - passwords don't match
	invalidForm := PasswordForm{
		Password:        "securepass123",
		ConfirmPassword: "different",
	}

	result, _ = v.Validate(invalidForm)
	fmt.Printf("Non-matching passwords result: %v\n", result.IsValid())
	for _, err := range result.Errors() {
		fmt.Printf("  - %s: %s (params: %v)\n", err.FieldPath, err.ErrorCode, err.Params)
	}

	fmt.Println()
}

// Example 6: Error handling patterns
func example6() {
	fmt.Println("=== Example 6: Error Handling Patterns ===")

	type Product struct {
		Name  string `json:"name" validate:"required,min_length=3"`
		Price int    `json:"price" validate:"required,min=0"`
		SKU   string `json:"sku" validate:"required"`
	}

	v, _ := validator.NewFromStruct(Product{})

	invalidProduct := Product{
		Name:  "AB",
		Price: -10,
		SKU:   "",
	}

	result, _ := v.Validate(invalidProduct)

	// Pattern 1: Check if valid
	if !result.IsValid() {
		fmt.Println("Product validation failed!")

		// Pattern 2: Iterate all errors
		fmt.Println("\nAll errors:")
		for _, err := range result.Errors() {
			fmt.Printf("  %s: %s %v\n", err.FieldPath, err.ErrorCode, err.Params)
		}

		// Pattern 3: Group by field
		fmt.Println("\nErrors by field:")
		errorsByField := result.ErrorsByField()
		for field, errs := range errorsByField {
			fmt.Printf("  %s: ", field)
			for _, err := range errs {
				fmt.Printf("%s ", err.ErrorCode)
			}
			fmt.Println()
		}

		// Pattern 4: Get first error
		firstErr := result.FirstError()
		fmt.Printf("\nFirst error: %s - %s\n", firstErr.FieldPath, firstErr.ErrorCode)
	}

	fmt.Println()
}

func main() {
	example1()
	example2()
	example3()
	example4()
	example5()
	example6()

	fmt.Println("All examples completed!")
}
