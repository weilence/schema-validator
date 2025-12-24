package main

import (
	"fmt"

	validator "github.com/weilence/schema-validator"
	sverrors "github.com/weilence/schema-validator/errors"
	"github.com/weilence/schema-validator/schema"
)

// Example 1: Tag-based validation
func example1() {
	fmt.Println("=== Example 1: Tag-based Validation ===")

	type User struct {
		Email    string `json:"email" validate:"required|email"`
		Password string `json:"password" validate:"required|min_length=8"`
		Confirm  string `json:"confirm" validate:"required|eqfield=Password"`
		Age      int    `json:"age" validate:"min=18,max=120"`
	}

	v, err := validator.New(User{})
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

	err = v.Validate(validUser)
	fmt.Printf("Valid user result: %v\n", err == nil)

	// Invalid user
	invalidUser := User{
		Email:    "not-an-email",
		Password: "short",
		Confirm:  "different",
		Age:      15,
	}

	err = v.Validate(invalidUser)
	fmt.Printf("Invalid user result: %v\n", err == nil)
	if err != nil {
		res, ok := err.(*sverrors.ValidationResult)
		if ok {
			for _, e := range res.Errors() {
				fmt.Printf("  - %s: %s\n", e.FieldPath, e.ErrorCode)
			}
		} else {
			fmt.Printf("  - %v\n", err)
		}
	}

	fmt.Println()
}

// Example 2: Code-based validation with maps
func example2() {
	fmt.Println("=== Example 2: Code-based Validation ===")

	userSchema := schema.Object().
		Field("email", schema.Field().
			AddValidator("required").
			AddValidator("email").
			Build()).
		Field("password", schema.Field().
			AddValidator("required").
			AddValidator("min_length", "8").
			Build()).
		Field("age", schema.Field().
			AddValidator("min", "18").
			AddValidator("max", "120").
			Build()).
		Build()

	v := validator.NewFromSchema(userSchema)

	// Valid map data
	validData := map[string]any{
		"email":    "jane@example.com",
		"password": "password123",
		"age":      25,
	}

	err := v.Validate(validData)
	fmt.Printf("Valid map result: %v\n", err == nil)

	// Invalid map data
	invalidData := map[string]any{
		"email":    "invalid-email",
		"password": "short",
		"age":      150,
	}

	err = v.Validate(invalidData)
	fmt.Printf("Invalid map result: %v\n", err == nil)
	if err != nil {
		res, ok := err.(*sverrors.ValidationResult)
		if ok {
			for _, e := range res.Errors() {
				fmt.Printf("  - %s: %s (params: %v)\n", e.FieldPath, e.ErrorCode, e.Params)
			}
		} else {
			fmt.Printf("  - %v\n", err)
		}
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
		Name string `json:"name" validate:"required|min_length=2"`
		Age  int    `json:"age" validate:"min=0,max=150"`
		Address
	}

	v, _ := validator.New(Person{})

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

	err := v.Validate(validPerson)
	fmt.Printf("Valid person result: %v\n", err == nil)

	// Invalid person - missing embedded field
	invalidPerson := Person{
		Name: "B", // Too short
		Age:  28,
		Address: Address{
			Street: "", // Required but empty
			City:   "New York",
		},
	}

	err = v.Validate(invalidPerson)
	fmt.Printf("Invalid person result: %v\n", err == nil)
	if err != nil {
		res, ok := err.(*sverrors.ValidationResult)
		if ok {
			for _, e := range res.Errors() {
				fmt.Printf("  - %s: %s\n", e.FieldPath, e.ErrorCode)
			}
		} else {
			fmt.Printf("  - %v\n", err)
		}
	}

	fmt.Println()
}

// Example 4: Array validation
func example4() {
	fmt.Println("=== Example 4: Array Validation ===")

	itemSchema := schema.Field().AddValidator("min_length", "1").Build()

	arraySchema := schema.Array(itemSchema).
		MinItems(1).
		MaxItems(5).
		Build()

	listSchema := schema.Object().
		Field("items", arraySchema).
		Build()

	v := validator.NewFromSchema(listSchema)

	// Valid list
	validList := map[string]any{
		"items": []string{"item1", "item2", "item3"},
	}

	err := v.Validate(validList)
	fmt.Printf("Valid list result: %v\n", err == nil)

	// Invalid - empty array
	invalidList := map[string]any{
		"items": []string{},
	}

	err = v.Validate(invalidList)
	fmt.Printf("Invalid list result (empty): %v\n", err == nil)
	if err != nil {
		res, ok := err.(*sverrors.ValidationResult)
		if ok {
			for _, e := range res.Errors() {
				fmt.Printf("  - %s: %s\n", e.FieldPath, e.ErrorCode)
			}
		} else {
			fmt.Printf("  - %v\n", err)
		}
	}

	// Invalid - too many items
	tooManyItems := map[string]any{
		"items": []string{"1", "2", "3", "4", "5", "6"},
	}

	err = v.Validate(tooManyItems)
	fmt.Printf("Invalid list result (too many): %v\n", err == nil)
	if err != nil {
		res, ok := err.(*sverrors.ValidationResult)
		if ok {
			for _, e := range res.Errors() {
				fmt.Printf("  - %s: %s (params: %v)\n", e.FieldPath, e.ErrorCode, e.Params)
			}
		} else {
			fmt.Printf("  - %v\n", err)
		}
	}

	fmt.Println()
}

// Example 5: Cross-field validation
func example5() {
	fmt.Println("=== Example 5: Cross-field Validation ===")

	type PasswordForm struct {
		Password        string `json:"password" validate:"required|min_length=8"`
		ConfirmPassword string `json:"confirmPassword" validate:"required|eqfield=Password"`
	}

	v, _ := validator.New(PasswordForm{})

	// Valid - passwords match
	validForm := PasswordForm{
		Password:        "securepass123",
		ConfirmPassword: "securepass123",
	}

	err := v.Validate(validForm)
	fmt.Printf("Matching passwords result: %v\n", err == nil)

	// Invalid - passwords don't match
	invalidForm := PasswordForm{
		Password:        "securepass123",
		ConfirmPassword: "different",
	}

	err = v.Validate(invalidForm)
	fmt.Printf("Non-matching passwords result: %v\n", err == nil)
	if err != nil {
		res, ok := err.(*sverrors.ValidationResult)
		if ok {
			for _, e := range res.Errors() {
				fmt.Printf("  - %s: %s (params: %v)\n", e.FieldPath, e.ErrorCode, e.Params)
			}
		} else {
			fmt.Printf("  - %v\n", err)
		}
	}

	fmt.Println()
}

// Example 6: Error handling patterns
func example6() {
	fmt.Println("=== Example 6: Error Handling Patterns ===")

	type Product struct {
		Name  string `json:"name" validate:"required|min_length=3"`
		Price int    `json:"price" validate:"required|min=0"`
		SKU   string `json:"sku" validate:"required"`
	}

	v, _ := validator.New(Product{})

	invalidProduct := Product{
		Name:  "AB",
		Price: -10,
		SKU:   "",
	}

	err := v.Validate(invalidProduct)

	// Pattern 1: Check if valid
	if err != nil {
		fmt.Println("Product validation failed!")

		// Pattern 2: Iterate all errors
		fmt.Println("\nAll errors:")
		if res, ok := err.(*sverrors.ValidationResult); ok {
			for _, e := range res.Errors() {
				fmt.Printf("  %s: %s %v\n", e.FieldPath, e.ErrorCode, e.Params)
			}
		} else {
			fmt.Printf("  %v\n", err)
		}

		// Pattern 3: Group by field
		fmt.Println("\nErrors by field:")
		if res, ok := err.(*sverrors.ValidationResult); ok {
			errorsByField := res.ErrorsByField()
			for field, errs := range errorsByField {
				fmt.Printf("  %s: ", field)
				for _, e := range errs {
					fmt.Printf("%s ", e.ErrorCode)
				}
				fmt.Println()
			}

			// Pattern 4: Get first error
			firstErr := res.FirstError()
			fmt.Printf("\nFirst error: %s - %s\n", firstErr.FieldPath, firstErr.ErrorCode)
		}
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
