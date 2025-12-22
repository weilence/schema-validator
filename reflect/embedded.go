package reflect

import (
	"reflect"
	"strings"
)

// walkFields recursively walks struct fields including embedded structs
func walkFields(typ reflect.Type, index []int, info *StructInfo) {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldIndex := append([]int(nil), index...)
		fieldIndex = append(fieldIndex, i)

		// Handle embedded structs
		if field.Anonymous {
			fieldType := field.Type
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}

			if fieldType.Kind() == reflect.Struct {
				// Recursively process embedded struct
				walkFields(fieldType, fieldIndex, info)
				continue
			}
		}

		// Skip non-embedded private fields (lowercase first letter)
		if !field.Anonymous && field.PkgPath != "" {
			continue
		}

		// Add field info
		fieldInfo := &FieldInfo{
			Name:      field.Name,
			Index:     fieldIndex,
			Type:      field.Type,
			Embedded:  len(index) > 0,
			Anonymous: field.Anonymous,
		}

		// Use tag name if available, otherwise field name
		name := getFieldName(field)

		// Only add if not already present (embedded fields can be shadowed)
		if _, exists := info.Fields[name]; !exists {
			info.Fields[name] = fieldInfo
			info.names = append(info.names, name)
		}
	}
}

// getFieldName extracts the field name from tags or uses the struct field name
func getFieldName(field reflect.StructField) string {
	// Check json tag first
	if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
		// Extract name before comma
		if idx := strings.Index(jsonTag, ","); idx != -1 {
			return jsonTag[:idx]
		}
		return jsonTag
	}

	// Check validate tag for field name
	if validateTag := field.Tag.Get("validate"); validateTag != "" && validateTag != "-" {
		// validate tag doesn't define field names, just validation rules
	}

	// Default to field name
	return field.Name
}
