package schema_test

import (
	"encoding/json"
	"testing"

	"github.com/weilence/schema-validator/schema"
	"github.com/weilence/schema-validator/validation"
)

// TestFieldSchemaToString tests ToString for FieldSchema
func TestFieldSchemaToString(t *testing.T) {
	// Test simple field schema with required validator
	fieldSchema := schema.Field().
		AddValidator(validation.Required()).
		Build()

	jsonStr := fieldSchema.ToString()

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("ToString output is not valid JSON: %v", err)
	}

	// Check type
	if result["type"] != "field" {
		t.Errorf("Expected type 'field', got %v", result["type"])
	}

	// Check optional flag
	if result["optional"] != false {
		t.Errorf("Expected optional false, got %v", result["optional"])
	}

	// Check validators
	validators, ok := result["validators"].([]interface{})
	if !ok || len(validators) == 0 {
		t.Error("Expected validators array")
	}
}

// TestFieldSchemaWithMultipleValidators tests field schema with multiple validators
func TestFieldSchemaWithMultipleValidators(t *testing.T) {
	fieldSchema := schema.Field().
		AddValidator(validation.Required()).
		AddValidator(validation.MinLength(5)).
		AddValidator(validation.MaxLength(100)).
		AddValidator(validation.Email()).
		SetOptional(false).
		Build()

	jsonStr := fieldSchema.ToString()

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("ToString output is not valid JSON: %v", err)
	}

	validators, ok := result["validators"].([]interface{})
	if !ok {
		t.Fatal("Expected validators array")
	}

	if len(validators) != 4 {
		t.Errorf("Expected 4 validators, got %d", len(validators))
	}

	// Check each validator
	validatorMap := validators[0].(map[string]interface{})
	if validatorMap["name"] != "required" {
		t.Errorf("Expected first validator to be 'required', got %v", validatorMap["name"])
	}

	validatorMap = validators[1].(map[string]interface{})
	if validatorMap["name"] != "min_length" {
		t.Errorf("Expected second validator to be 'min_length', got %v", validatorMap["name"])
	}
	if validatorMap["value"] != float64(5) {
		t.Errorf("Expected min_length value 5, got %v", validatorMap["value"])
	}
}

// TestArraySchemaToString tests ToString for ArraySchema
func TestArraySchemaToString(t *testing.T) {
	// Create an array schema with field element
	elementSchema := schema.Field().
		AddValidator(validation.Required()).
		Build()

	arraySchema := schema.Array(elementSchema).
		MinItems(1).
		MaxItems(10).
		Build()

	jsonStr := arraySchema.ToString()

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("ToString output is not valid JSON: %v", err)
	}

	// Check type
	if result["type"] != "array" {
		t.Errorf("Expected type 'array', got %v", result["type"])
	}

	// Check minItems and maxItems
	if result["minItems"] != float64(1) {
		t.Errorf("Expected minItems 1, got %v", result["minItems"])
	}
	if result["maxItems"] != float64(10) {
		t.Errorf("Expected maxItems 10, got %v", result["maxItems"])
	}

	// Check element schema is nested properly
	element, ok := result["element"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected element to be an object")
	}

	if element["type"] != "field" {
		t.Errorf("Expected element type 'field', got %v", element["type"])
	}
}

// TestObjectSchemaToString tests ToString for ObjectSchema
func TestObjectSchemaToString(t *testing.T) {
	// Create a simple object schema
	objSchema := schema.Object().
		Field("name", schema.Field().AddValidator(validation.Required()).Build()).
		Field("email", schema.Field().AddValidator(validation.Email()).Build()).
		Field("age", schema.Field().AddValidator(validation.Min(18)).Build()).
		Build()

	jsonStr := objSchema.ToString()

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("ToString output is not valid JSON: %v", err)
	}

	// Check type
	if result["type"] != "object" {
		t.Errorf("Expected type 'object', got %v", result["type"])
	}

	// Check fields
	fields, ok := result["fields"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected fields to be an object")
	}

	if len(fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(fields))
	}

	// Check name field
	nameField, ok := fields["name"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected name field to be an object")
	}
	if nameField["type"] != "field" {
		t.Errorf("Expected name field type 'field', got %v", nameField["type"])
	}
}

// TestNestedObjectSchemaToString tests ToString with nested objects
func TestNestedObjectSchemaToString(t *testing.T) {
	// Create nested object schema
	addressSchema := schema.Object().
		Field("street", schema.Field().AddValidator(validation.Required()).Build()).
		Field("city", schema.Field().AddValidator(validation.Required()).Build()).
		Field("zipCode", schema.Field().AddValidator(validation.MinLength(5)).Build()).
		Build()

	userSchema := schema.Object().
		Field("name", schema.Field().AddValidator(validation.Required()).Build()).
		Field("address", addressSchema).
		Build()

	jsonStr := userSchema.ToString()

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("ToString output is not valid JSON: %v", err)
	}

	// Check nested structure
	fields := result["fields"].(map[string]interface{})
	addressField := fields["address"].(map[string]interface{})

	if addressField["type"] != "object" {
		t.Errorf("Expected address type 'object', got %v", addressField["type"])
	}

	// Check nested fields
	addressFields := addressField["fields"].(map[string]interface{})
	if len(addressFields) != 3 {
		t.Errorf("Expected 3 address fields, got %d", len(addressFields))
	}

	streetField := addressFields["street"].(map[string]interface{})
	if streetField["type"] != "field" {
		t.Errorf("Expected street type 'field', got %v", streetField["type"])
	}
}

// TestArrayOfObjectsToString tests ToString with array of objects
func TestArrayOfObjectsToString(t *testing.T) {
	// Create an array of objects schema
	itemSchema := schema.Object().
		Field("id", schema.Field().AddValidator(validation.Required()).Build()).
		Field("name", schema.Field().AddValidator(validation.Required()).Build()).
		Build()

	arraySchema := schema.Array(itemSchema).
		MinItems(1).
		Build()

	jsonStr := arraySchema.ToString()

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("ToString output is not valid JSON: %v", err)
	}

	// Check array element is object
	element := result["element"].(map[string]interface{})
	if element["type"] != "object" {
		t.Errorf("Expected element type 'object', got %v", element["type"])
	}

	// Check object fields
	fields := element["fields"].(map[string]interface{})
	if len(fields) != 2 {
		t.Errorf("Expected 2 fields in element object, got %d", len(fields))
	}
}

// TestComplexNestedSchemaToString tests a complex nested structure
func TestComplexNestedSchemaToString(t *testing.T) {
	// Create a complex schema with multiple levels of nesting
	phoneSchema := schema.Object().
		Field("type", schema.Field().AddValidator(validation.Required()).Build()).
		Field("number", schema.Field().AddValidator(validation.Required()).Build()).
		Build()

	addressSchema := schema.Object().
		Field("street", schema.Field().AddValidator(validation.Required()).Build()).
		Field("city", schema.Field().AddValidator(validation.Required()).Build()).
		Build()

	userSchema := schema.Object().
		Field("name", schema.Field().AddValidator(validation.Required()).Build()).
		Field("email", schema.Field().AddValidator(validation.Email()).Build()).
		Field("phones", schema.Array(phoneSchema).MinItems(1).Build()).
		Field("address", addressSchema).
		Build()

	jsonStr := userSchema.ToString()

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("ToString output is not valid JSON: %v", err)
	}

	// Verify structure
	if result["type"] != "object" {
		t.Errorf("Expected type 'object', got %v", result["type"])
	}

	fields := result["fields"].(map[string]interface{})

	// Check phones field is array
	phonesField := fields["phones"].(map[string]interface{})
	if phonesField["type"] != "array" {
		t.Errorf("Expected phones type 'array', got %v", phonesField["type"])
	}

	// Check phones element is object
	phonesElement := phonesField["element"].(map[string]interface{})
	if phonesElement["type"] != "object" {
		t.Errorf("Expected phones element type 'object', got %v", phonesElement["type"])
	}

	// Check address field is object
	addressField := fields["address"].(map[string]interface{})
	if addressField["type"] != "object" {
		t.Errorf("Expected address type 'object', got %v", addressField["type"])
	}
}

// TestCrossFieldValidatorInToString tests cross-field validators in ToString
func TestCrossFieldValidatorInToString(t *testing.T) {
	fieldSchema := schema.Field().
		AddValidator(validation.Required()).
		AddValidator(validation.EqField("password")).
		Build()

	jsonStr := fieldSchema.ToString()

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("ToString output is not valid JSON: %v", err)
	}

	validators := result["validators"].([]interface{})
	if len(validators) != 2 {
		t.Errorf("Expected 2 validators, got %d", len(validators))
	}

	// Check second validator is eqfield
	eqfieldValidator := validators[1].(map[string]interface{})
	if eqfieldValidator["name"] != "eqfield" {
		t.Errorf("Expected validator name 'eqfield', got %v", eqfieldValidator["name"])
	}
	if eqfieldValidator["value"] != "password" {
		t.Errorf("Expected eqfield value 'password', got %v", eqfieldValidator["value"])
	}
}
