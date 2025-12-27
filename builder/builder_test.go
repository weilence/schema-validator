package builder

import (
	"testing"

	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/schema"
)

// TestFieldSchemaFields tests FieldSchema字段内容
func TestFieldSchemaFields(t *testing.T) {
	fieldSchema := Field().
		AddValidator("required").
		Build()
	ctx := schema.NewContext(fieldSchema, data.New(123))
	if err := fieldSchema.Validate(ctx); err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}
}

// TestFieldSchemaWithMultipleValidators tests field schema with multiple validators
func TestFieldSchemaWithMultipleValidators(t *testing.T) {
	fieldSchema := Field().
		AddValidator("required").
		AddValidator("min", 5).
		AddValidator("max", 100).
		Optional().
		Build()
	ctx := schema.NewContext(fieldSchema, data.New(23))
	if err := fieldSchema.Validate(ctx); err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}
}

// TestArraySchemaFields tests ArraySchema字段内容
func TestArraySchemaFields(t *testing.T) {
	elementSchema := Field().
		AddValidator("email").
		Build()

	arraySchema := Array(elementSchema).
		AddValidator("min", 1).
		AddValidator("max", 10).
		Build()

	ctx := schema.NewContext(arraySchema, data.New([]any{"a@example.com", "b@example.com", "c@example.com"}))
	if err := arraySchema.Validate(ctx); err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}
}

// TestObjectSchemaFields tests ObjectSchema字段内容
func TestObjectSchemaFields(t *testing.T) {
	objSchema := Object().
		Field("name", Field().AddValidator("required").Build()).
		Field("email", Field().AddValidator("email").Build()).
		Field("age", Field().AddValidator("min", 18).Build()).
		Build()

	ctx := schema.NewContext(objSchema, data.New(map[string]any{"name": "a", "email": "a@example.com", "age": 18}))
	if err := objSchema.Validate(ctx); err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}
}

// TestNestedObjectSchemaFields tests 嵌套对象schema字段内容
func TestNestedObjectSchemaFields(t *testing.T) {
	addressSchema := Object().
		Field("street", Field().AddValidator("required").Build()).
		Field("city", Field().AddValidator("required").Build()).
		Field("zipCode", Field().AddValidator("min", 5).Build()).
		Build()

	userSchema := Object().
		Field("name", Field().AddValidator("required").Build()).
		Field("address", addressSchema).
		Build()

	ctx := schema.NewContext(userSchema, data.New(map[string]any{"name": "a", "address": map[string]any{"street": "s", "city": "c", "zipCode": 12345}}))
	if err := userSchema.Validate(ctx); err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}
}

// TestArrayOfObjectsSchemaFields tests 数组元素为对象的schema字段内容
func TestArrayOfObjectsSchemaFields(t *testing.T) {
	itemSchema := Object().
		Field("id", Field().AddValidator("required").Build()).
		Field("name", Field().AddValidator("required").Build()).
		Build()

	arraySchema := Array(itemSchema).
		AddValidator("min", 1).
		Build()
	ctx := schema.NewContext(arraySchema, data.New([]any{
		map[string]any{"id": 1, "name": "a"},
		map[string]any{"id": 2, "name": "b"},
	}))
	if err := arraySchema.Validate(ctx); err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}
}

// TestComplexNestedSchemaFields tests 多层嵌套schema字段内容
func TestComplexNestedSchemaFields(t *testing.T) {
	phoneSchema := Object().
		Field("type", Field().AddValidator("required").Build()).
		Field("number", Field().AddValidator("required").Build()).
		Build()

	addressSchema := Object().
		Field("street", Field().AddValidator("required").Build()).
		Field("city", Field().AddValidator("required").Build()).
		Build()

	userSchema := Object().
		Field("name", Field().AddValidator("required").Build()).
		Field("email", Field().AddValidator("email").Build()).
		Field("phones", Array(phoneSchema).AddValidator("min", 1).Build()).
		Field("address", addressSchema).
		Build()
	ctx := schema.NewContext(userSchema, data.New(map[string]any{
		"name":  "a",
		"email": "a@example.com",
		"phones": []any{
			map[string]any{"type": "mobile", "number": "123"},
		},
		"address": map[string]any{"street": "s", "city": "c"},
	}))
	if err := userSchema.Validate(ctx); err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}
}

// TestCrossFieldValidatorFields tests 跨字段validator
func TestCrossFieldValidatorFields(t *testing.T) {
	fieldSchema := Object().
		Field("password", Field().AddValidator("required").Build()).
		Field("confirmPassword", Field().AddValidator("eqfield", "password").Build()).
		Build()
	ctx := schema.NewContext(fieldSchema, data.New(map[string]any{
		"password":        "secret",
		"confirmPassword": "secret",
	}))
	if err := fieldSchema.Validate(ctx); err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}
}
