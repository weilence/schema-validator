package tags

import (
	"reflect"
	"strings"

	"github.com/weilence/schema-validator/schema"
)

// ParseStructTags parses struct tags and builds an ObjectSchema
func ParseStructTags(typ reflect.Type) (*schema.ObjectSchema, error) {
	return ParseStructTagsWithRegistry(typ, DefaultRegistry())
}

// ParseStructTagsWithRegistry parses struct tags with a custom registry
func ParseStructTagsWithRegistry(typ reflect.Type, registry *Registry) (*schema.ObjectSchema, error) {
	objSchema := schema.NewObjectSchema()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Skip unexported fields (unless embedded)
		if !field.Anonymous && field.PkgPath != "" {
			continue
		}

		// Handle embedded structs
		if field.Anonymous {
			fieldType := field.Type
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}

			if fieldType.Kind() == reflect.Struct {
				// Parse embedded struct and merge fields
				embeddedSchema, err := ParseStructTagsWithRegistry(fieldType, registry)
				if err != nil {
					return nil, err
				}

				// Merge embedded fields into parent
				for _, fieldName := range embeddedSchema.Fields() {
					if fieldSchema, ok := embeddedSchema.GetField(fieldName); ok {
						objSchema.AddField(fieldName, fieldSchema)
					}
				}
				continue
			}
		}

		// Get field name
		fieldName := getFieldName(field)
		if fieldName == "" || fieldName == "-" {
			continue
		}

		// Parse validation tag
		validateTag := field.Tag.Get("validate")
		if validateTag == "" || validateTag == "-" {
			// No validation, add optional field schema
			objSchema.AddField(fieldName, schema.NewFieldSchema().SetOptional(true))
			continue
		}

		// Parse field schema from tag
		fieldSchema, err := parseFieldTag(field.Type, validateTag, registry)
		if err != nil {
			return nil, err
		}

		objSchema.AddField(fieldName, fieldSchema)
	}

	return objSchema, nil
}

func parseFieldTag(fieldType reflect.Type, tag string, registry *Registry) (schema.Schema, error) {
	// Handle slice/array types
	if fieldType.Kind() == reflect.Slice || fieldType.Kind() == reflect.Array {
		elemType := fieldType.Elem()
		elemSchema, err := parseFieldTag(elemType, tag, registry)
		if err != nil {
			return nil, err
		}

		arraySchema := schema.NewArraySchema(elemSchema)

		// Parse array-specific validators
		rules := parseTag(tag)
		for _, rule := range rules {
			validator, err := registry.GetArrayValidator(rule.Name, rule.Param)
			if err != nil {
				return nil, err
			}
			if validator != nil {
				arraySchema.AddValidator(validator)
			}
		}

		return arraySchema, nil
	}

	// Handle struct types
	if fieldType.Kind() == reflect.Struct {
		// Recursively parse nested struct
		return ParseStructTagsWithRegistry(fieldType, registry)
	}

	// Handle map types
	if fieldType.Kind() == reflect.Map {
		// For maps, create an object schema
		return schema.NewObjectSchema(), nil
	}

	// Primitive field
	fieldSchema := schema.NewFieldSchema()

	// Parse validation rules
	rules := parseTag(tag)

	hasRequired := false
	for _, rule := range rules {
		if rule.Name == "required" {
			hasRequired = true
			fieldSchema.SetOptional(false)
		}

		validator, err := registry.GetFieldValidator(rule.Name, rule.Param)
		if err != nil {
			return nil, err
		}

		if validator != nil {
			fieldSchema.AddValidator(validator)
		}
	}

	if !hasRequired {
		fieldSchema.SetOptional(true)
	}

	return fieldSchema, nil
}

// TagRule represents a parsed validation rule
type TagRule struct {
	Name  string
	Param string
}

// parseTag parses a validation tag into rules
// Example: "required,min=5,max=100" -> [{required, ""}, {min, "5"}, {max, "100"}]
// For multi-param: "between=10:100" -> [{between, "10:100"}]
// The params will be split later in the registry
func parseTag(tag string) []TagRule {
	if tag == "" {
		return nil
	}

	// Split by comma, but be smart about it
	// We need to handle cases like: "between=10,100" where the comma is part of the parameter
	// Strategy: Split only at commas that are NOT between an '=' and the next validator name

	rules := make([]TagRule, 0)
	currentRule := ""
	inParam := false

	for i := 0; i < len(tag); i++ {
		ch := tag[i]

		if ch == '=' {
			inParam = true
			currentRule += string(ch)
		} else if ch == ',' {
			// Check if we're in a parameter value
			// Look ahead to see if there's another '=' before the next ','
			if inParam {
				// Check if the next part looks like a new rule (contains '=') or is just a param
				nextPart := ""
				for j := i + 1; j < len(tag); j++ {
					if tag[j] == ',' {
						break
					}
					nextPart += string(tag[j])
				}

				// If nextPart doesn't contain '=' and doesn't look like a rule name, it's a parameter
				if !strings.Contains(nextPart, "=") && !isValidatorName(nextPart) {
					// This comma is part of the parameter
					currentRule += string(ch)
				} else {
					// This comma ends the current rule
					inParam = false
					if currentRule != "" {
						rules = append(rules, parseRule(currentRule))
						currentRule = ""
					}
				}
			} else {
				// Not in a parameter, this comma separates rules
				if currentRule != "" {
					rules = append(rules, parseRule(currentRule))
					currentRule = ""
				}
			}
		} else {
			currentRule += string(ch)
		}
	}

	// Don't forget the last rule
	if currentRule != "" {
		rules = append(rules, parseRule(currentRule))
	}

	return rules
}

// parseRule parses a single rule string like "min=5" or "required"
func parseRule(ruleStr string) TagRule {
	ruleStr = strings.TrimSpace(ruleStr)

	if idx := strings.Index(ruleStr, "="); idx != -1 {
		return TagRule{
			Name:  strings.TrimSpace(ruleStr[:idx]),
			Param: strings.TrimSpace(ruleStr[idx+1:]),
		}
	}

	return TagRule{
		Name:  ruleStr,
		Param: "",
	}
}

// isValidatorName checks if a string looks like a validator name
// Validator names typically don't contain digits or special characters
func isValidatorName(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}

	// Check if it's a known validator name (simple heuristic)
	// If it starts with a letter and contains only letters and underscores, it's likely a validator name
	for i, ch := range s {
		if i == 0 {
			if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')) {
				return false
			}
		} else {
			if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_') {
				return false
			}
		}
	}

	return true
}

func getFieldName(field reflect.StructField) string {
	// Check json tag first
	if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
		// Extract name before comma
		if idx := strings.Index(jsonTag, ","); idx != -1 {
			return jsonTag[:idx]
		}
		return jsonTag
	}

	// Default to field name
	return field.Name
}
