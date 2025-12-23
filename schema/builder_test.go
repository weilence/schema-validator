package schema

import (
	"encoding/json"
	"testing"
)

// TestFieldSchemaToString tests ToString for FieldSchema
func TestFieldSchemaToString(t *testing.T) {
	// Test simple field schema with required validator
	fieldSchema := Field().
		AddValidator("required").
		Build()

	jsonStr := fieldSchema.ToString()

	// Verify it's valid JSON
	var result map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("ToString output is not valid JSON: %v", err)
	}

	// Check type
	if result["type"] != "field" {
		t.Errorf("Expected type 'field', got %v", result["type"])
	}

	// Check validators
	validators, ok := result["validators"].([]any)
	if !ok || len(validators) == 0 {
		t.Error("Expected validators array")
	}
}

// TestFieldSchemaWithMultipleValidators tests field schema with multiple validators
func TestFieldSchemaWithMultipleValidators(t *testing.T) {
	fieldSchema := Field().
		AddValidator("required").
		AddValidator("min_length", "5").
		AddValidator("max_length", "100").
		AddValidator("email").
		Optional().
		Build()

	jsonStr := fieldSchema.ToString()

	var result map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("ToString output is not valid JSON: %v", err)
	}

	validators, ok := result["validators"].([]any)
	if !ok {
		t.Fatal("Expected validators array")
	}

	if len(validators) != 3 {
		t.Errorf("Expected 3 validators, got %d", len(validators))
	}

	// Check each validator
	validatorMap := validators[0].(map[string]any)
	if validatorMap["name"] != "min_length" {
		t.Errorf("Expected first validator to be 'min_length', got %v", validatorMap["name"])
	}

	if validatorMap["value"] != float64(5) {
		t.Errorf("Expected min_length value 5, got %v", validatorMap["value"])
	}

	validatorMap = validators[1].(map[string]any)
	if validatorMap["name"] != "max_length" {
		t.Errorf("Expected second validator to be 'max_length', got %v", validatorMap["name"])
	}

	validatorMap = validators[2].(map[string]any)
	if validatorMap["name"] != "email" {
		t.Errorf("Expected third validator to be 'email', got %v", validatorMap["name"])
	}
}

// TestArraySchemaToString tests ToString for ArraySchema
func TestArraySchemaToString(t *testing.T) {
	// Create an array schema with field element
	elementSchema := Field().
		AddValidator("required").
		Build()

	arraySchema := Array(elementSchema).
		MinItems(1).
		MaxItems(10).
		Build()

	jsonStr := arraySchema.ToString()

	var result map[string]any
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
	element, ok := result["element"].(map[string]any)
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
	objSchema := Object().
		Field("name", Field().AddValidator("required").Build()).
		Field("email", Field().AddValidator("email").Build()).
		Field("age", Field().AddValidator("min", "18").Build()).
		Build()

	jsonStr := objSchema.ToString()

	var result map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("ToString output is not valid JSON: %v", err)
	}

	// Check type
	if result["type"] != "object" {
		t.Errorf("Expected type 'object', got %v", result["type"])
	}

	// Check fields
	fields, ok := result["fields"].(map[string]any)
	if !ok {
		t.Fatal("Expected fields to be an object")
	}

	if len(fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(fields))
	}

	// Check name field
	nameField, ok := fields["name"].(map[string]any)
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
	addressSchema := Object().
		Field("street", Field().AddValidator("required").Build()).
		Field("city", Field().AddValidator("required").Build()).
		Field("zipCode", Field().AddValidator("min_length", "5").Build()).
		Build()

	userSchema := Object().
		Field("name", Field().AddValidator("required").Build()).
		Field("address", addressSchema).
		Build()

	jsonStr := userSchema.ToString()

	var result map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("ToString output is not valid JSON: %v", err)
	}

	// Check nested structure
	fields := result["fields"].(map[string]any)
	addressField := fields["address"].(map[string]any)

	if addressField["type"] != "object" {
		t.Errorf("Expected address type 'object', got %v", addressField["type"])
	}

	// Check nested fields
	addressFields := addressField["fields"].(map[string]any)
	if len(addressFields) != 3 {
		t.Errorf("Expected 3 address fields, got %d", len(addressFields))
	}

	streetField := addressFields["street"].(map[string]any)
	if streetField["type"] != "field" {
		t.Errorf("Expected street type 'field', got %v", streetField["type"])
	}
}

// TestArrayOfObjectsToString tests ToString with array of objects
func TestArrayOfObjectsToString(t *testing.T) {
	// Create an array of objects schema
	itemSchema := Object().
		Field("id", Field().AddValidator("required").Build()).
		Field("name", Field().AddValidator("required").Build()).
		Build()

	arraySchema := Array(itemSchema).
		MinItems(1).
		Build()

	jsonStr := arraySchema.ToString()

	var result map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("ToString output is not valid JSON: %v", err)
	}

	// Check array element is object
	element := result["element"].(map[string]any)
	if element["type"] != "object" {
		t.Errorf("Expected element type 'object', got %v", element["type"])
	}

	// Check object fields
	fields := element["fields"].(map[string]any)
	if len(fields) != 2 {
		t.Errorf("Expected 2 fields in element object, got %d", len(fields))
	}
}

// TestComplexNestedSchemaToString tests a complex nested structure
func TestComplexNestedSchemaToString(t *testing.T) {
	// Create a complex schema with multiple levels of nesting
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
		Field("phones", Array(phoneSchema).MinItems(1).Build()).
		Field("address", addressSchema).
		Build()

	jsonStr := userSchema.ToString()

	// Verify it's valid JSON
	var result map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("ToString output is not valid JSON: %v", err)
	}

	// Verify structure
	if result["type"] != "object" {
		t.Errorf("Expected type 'object', got %v", result["type"])
	}

	fields := result["fields"].(map[string]any)

	// Check phones field is array
	phonesField := fields["phones"].(map[string]any)
	if phonesField["type"] != "array" {
		t.Errorf("Expected phones type 'array', got %v", phonesField["type"])
	}

	// Check phones element is object
	phonesElement := phonesField["element"].(map[string]any)
	if phonesElement["type"] != "object" {
		t.Errorf("Expected phones element type 'object', got %v", phonesElement["type"])
	}

	// Check address field is object
	addressField := fields["address"].(map[string]any)
	if addressField["type"] != "object" {
		t.Errorf("Expected address type 'object', got %v", addressField["type"])
	}
}

// TestCrossFieldValidatorInToString tests cross-field validators in ToString
func TestCrossFieldValidatorInToString(t *testing.T) {
	fieldSchema := Field().
		AddValidator("required").
		AddValidator("eqfield", "password").
		Build()

	jsonStr := fieldSchema.ToString()

	var result map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("ToString output is not valid JSON: %v", err)
	}

	validators := result["validators"].([]any)
	if len(validators) != 2 {
		t.Errorf("Expected 2 validators, got %d", len(validators))
	}

	// Check second validator is eqfield
	eqfieldValidator := validators[1].(map[string]any)
	if eqfieldValidator["name"] != "eqfield" {
		t.Errorf("Expected validator name 'eqfield', got %v", eqfieldValidator["name"])
	}
	if eqfieldValidator["value"] != "password" {
		t.Errorf("Expected eqfield value 'password', got %v", eqfieldValidator["value"])
	}
}
