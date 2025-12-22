package main

import (
	"fmt"

	"github.com/weilence/schema-validator/schema"
	"github.com/weilence/schema-validator/validation"
)

func main() {
	fmt.Println("=== Schema ToString Examples ===\n")

	// Example 1: Simple Field Schema
	fmt.Println("1. Simple Field Schema:")
	fieldSchema := schema.Field().
		AddValidator(validation.Required()).
		AddValidator(validation.MinLength(5)).
		AddValidator(validation.Email()).
		Build()
	fmt.Println(fieldSchema.ToString())
	fmt.Println()

	// Example 2: Array Schema
	fmt.Println("2. Array Schema:")
	arraySchema := schema.Array(
		schema.Field().AddValidator(validation.Required()).Build(),
	).MinItems(1).MaxItems(10).Build()
	fmt.Println(arraySchema.ToString())
	fmt.Println()

	// Example 3: Simple Object Schema
	fmt.Println("3. Simple Object Schema:")
	userSchema := schema.Object().
		Field("name", schema.Field().AddValidator(validation.Required()).Build()).
		Field("email", schema.Field().AddValidator(validation.Email()).Build()).
		Field("age", schema.Field().AddValidator(validation.Min(18)).Build()).
		Build()
	fmt.Println(userSchema.ToString())
	fmt.Println()

	// Example 4: Nested Object Schema
	fmt.Println("4. Nested Object Schema:")
	addressSchema := schema.Object().
		Field("street", schema.Field().AddValidator(validation.Required()).Build()).
		Field("city", schema.Field().AddValidator(validation.Required()).Build()).
		Field("zipCode", schema.Field().AddValidator(validation.MinLength(5)).Build()).
		Build()

	personSchema := schema.Object().
		Field("name", schema.Field().AddValidator(validation.Required()).Build()).
		Field("email", schema.Field().AddValidator(validation.Email()).Build()).
		Field("address", addressSchema).
		Build()
	fmt.Println(personSchema.ToString())
	fmt.Println()

	// Example 5: Complex Schema with Arrays and Objects
	fmt.Println("5. Complex Schema (User with Phones):")
	phoneSchema := schema.Object().
		Field("type", schema.Field().AddValidator(validation.Required()).Build()).
		Field("number", schema.Field().AddValidator(validation.Required()).Build()).
		Build()

	complexUserSchema := schema.Object().
		Field("name", schema.Field().AddValidator(validation.Required()).Build()).
		Field("email", schema.Field().AddValidator(validation.Email()).Build()).
		Field("phones", schema.Array(phoneSchema).MinItems(1).Build()).
		Field("address", addressSchema).
		Build()
	fmt.Println(complexUserSchema.ToString())
}
