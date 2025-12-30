package builder

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/weilence/schema-validator/schema"
	"github.com/weilence/schema-validator/validators"
)

type ParseConfig struct {
	registry *validators.Registry

	ruleSplitter       rune
	nameParamSeparator rune
	paramsSeparator    rune
	diveTag            string
}

func defaultParseConfig() *ParseConfig {
	return &ParseConfig{
		registry:           validators.DefaultRegistry(),
		ruleSplitter:       '|',
		nameParamSeparator: '=',
		paramsSeparator:    ',',
		diveTag:            "dive",
	}
}

type ParseOption func(*ParseConfig)

func WithRegistry(registry *validators.Registry) ParseOption {
	return func(cfg *ParseConfig) {
		cfg.registry = registry
	}
}

func WithRuleSplitter(r rune) ParseOption {
	return func(cfg *ParseConfig) {
		cfg.ruleSplitter = r
	}
}

func WithNameParamSeparator(r rune) ParseOption {
	return func(cfg *ParseConfig) {
		cfg.nameParamSeparator = r
	}
}

func WithParamsSeparator(r rune) ParseOption {
	return func(cfg *ParseConfig) {
		cfg.paramsSeparator = r
	}
}

func WithDiveTag(tag string) ParseOption {
	return func(cfg *ParseConfig) {
		cfg.diveTag = tag
	}
}

// Parse parses struct tags and builds an ObjectSchema
func Parse(rt reflect.Type, opts ...ParseOption) (*schema.ObjectSchema, error) {
	cfg := defaultParseConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	return parse(rt, cfg)
}

func parse(rt reflect.Type, cfg *ParseConfig) (*schema.ObjectSchema, error) {
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}

	objSchema := schema.NewObject()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if err := parseStructField(objSchema, field, cfg); err != nil {
			return nil, err
		}
	}

	return objSchema, nil
}

func parseStructField(schema *schema.ObjectSchema, field reflect.StructField, cfg *ParseConfig) error {
	fieldType := field.Type
	if fieldType.Kind() == reflect.Pointer {
		fieldType = fieldType.Elem()
	}

	if field.Anonymous {
		if fieldType.Kind() != reflect.Struct {
			return nil
		}

		for i := 0; i < fieldType.NumField(); i++ {
			embeddedField := fieldType.Field(i)
			if err := parseStructField(schema, embeddedField, cfg); err != nil {
				return err
			}
		}

		return nil
	}

	if !field.IsExported() {
		return nil
	}

	fieldName := getFieldName(field)
	validateTag := field.Tag.Get("validate")
	if validateTag == "-" {
		return nil
	}

	tags := parseTag(validateTag, cfg)
	fieldSchema, err := parseField(field.Type, tags, cfg)
	if err != nil {
		return err
	}

	schema.AddField(fieldName, fieldSchema).AddFieldName(fieldName, field.Name)
	return nil
}

func parseField(fieldType reflect.Type, tags []TagRule, cfg *ParseConfig) (schema.Schema, error) {
	if fieldType.Kind() == reflect.Pointer {
		fieldType = fieldType.Elem()
	}

	// Handle slice/array types
	if fieldType.Kind() == reflect.Slice || fieldType.Kind() == reflect.Array {
		// Parse array-specific validators
		diveRuleInex := slices.IndexFunc(tags, func(item TagRule) bool { return item.Name == cfg.diveTag })
		var arrayTags, itemTags []TagRule
		if diveRuleInex >= 0 {
			arrayTags = tags[:diveRuleInex]
			itemTags = tags[diveRuleInex+1:]
		} else {
			arrayTags = tags
		}

		elemType := fieldType.Elem()

		elemSchema, err := parseField(elemType, itemTags, cfg)
		if err != nil {
			return nil, err
		}

		arraySchema := schema.NewArray(elemSchema)

		for _, rule := range arrayTags {
			params := convertValidatorParams(rule.Name, rule.Params, cfg)
			v := cfg.registry.NewValidator(rule.Name, params...)
			if v != nil {
				arraySchema.AddValidator(v)
			}
		}

		return arraySchema, nil
	}

	// Handle struct types
	if fieldType.Kind() == reflect.Struct {
		// Recursively parse nested struct
		return parse(fieldType, cfg)
	}

	// Handle map types
	if fieldType.Kind() == reflect.Map {
		// For maps, create an object schema
		return schema.NewObject(), nil
	}

	// Primitive field
	fieldSchema := schema.NewField()
	for _, tag := range tags {
		params := convertValidatorParams(tag.Name, tag.Params, cfg)
		v := cfg.registry.NewValidator(tag.Name, params...)
		if v != nil {
			fieldSchema.AddValidator(v)
		}
	}

	return fieldSchema, nil
}

func convertValidatorParams(name string, paramStrs []string, cfg *ParseConfig) []any {
	paramTypes := cfg.registry.GetValidatorParamTypes(name)

	paramTypesLen := len(paramTypes)
	if paramTypesLen == 0 && len(paramStrs) != 0 {
		panic(fmt.Sprintf("%s does not take any parameters", name))
	}

	if paramTypesLen == 1 {
		paramType := paramTypes[0]
		kind := paramType.Kind()
		switch kind {
		case reflect.Array:
			res := reflect.ArrayOf(paramType.Len(), paramType.Elem())
			rv := reflect.New(res).Elem()
			for i, paramStr := range paramStrs {
				elem := parseValidatorParam(paramType.Elem(), paramStr)
				rv.Index(i).Set(reflect.ValueOf(elem))
			}

			return []any{rv.Interface()}
		case reflect.Slice:
			res := reflect.MakeSlice(paramType, 0, 0)
			for _, paramStr := range paramStrs {
				elem := parseValidatorParam(paramType.Elem(), paramStr)
				res = reflect.Append(res, reflect.ValueOf(elem))
			}

			return []any{res.Interface()}
		default:
			if len(paramStrs) != 1 {
				panic(fmt.Sprintf("%s expected 1 parameter, got %d", name, len(paramStrs)))
			}

			return []any{parseValidatorParam(paramType, paramStrs[0])}
		}
	}

	if len(paramStrs) != paramTypesLen {
		panic(fmt.Sprintf("%s expected %d parameters, got %d", name, paramTypesLen, len(paramStrs)))
	}

	params := make([]any, paramTypesLen)
	for i, paramType := range paramTypes {
		params[i] = parseValidatorParam(paramType, paramStrs[i])
	}

	return params
}

func parseValidatorParam(paramType reflect.Type, paramValue string) any {
	switch paramType.Kind() {
	case reflect.Bool:
		var boolVal bool
		_, err := fmt.Sscanf(paramValue, "%t", &boolVal)
		if err != nil {
			panic(fmt.Sprintf("invalid bool parameter: %s", paramValue))
		}
		return boolVal
	case reflect.Int:
		var intVal int
		_, err := fmt.Sscanf(paramValue, "%d", &intVal)
		if err != nil {
			panic(fmt.Sprintf("invalid int parameter: %s", paramValue))
		}
		return intVal
	case reflect.Int8:
		var intVal int8
		_, err := fmt.Sscanf(paramValue, "%d", &intVal)
		if err != nil {
			panic(fmt.Sprintf("invalid int8 parameter: %s", paramValue))
		}
		return intVal
	case reflect.Int16:
		var intVal int16
		_, err := fmt.Sscanf(paramValue, "%d", &intVal)
		if err != nil {
			panic(fmt.Sprintf("invalid int16 parameter: %s", paramValue))
		}
		return intVal
	case reflect.Int32:
		var intVal int32
		_, err := fmt.Sscanf(paramValue, "%d", &intVal)
		if err != nil {
			panic(fmt.Sprintf("invalid int32 parameter: %s", paramValue))
		}
		return intVal
	case reflect.Int64:
		var intVal int64
		_, err := fmt.Sscanf(paramValue, "%d", &intVal)
		if err != nil {
			panic(fmt.Sprintf("invalid int64 parameter: %s", paramValue))
		}
		return intVal
	case reflect.Uint:
		var uintVal uint
		_, err := fmt.Sscanf(paramValue, "%d", &uintVal)
		if err != nil {
			panic(fmt.Sprintf("invalid uint parameter: %s", paramValue))
		}
		return uintVal
	case reflect.Uint8:
		var uintVal uint8
		_, err := fmt.Sscanf(paramValue, "%d", &uintVal)
		if err != nil {
			panic(fmt.Sprintf("invalid uint8 parameter: %s", paramValue))
		}
		return uintVal
	case reflect.Uint16:
		var uintVal uint16
		_, err := fmt.Sscanf(paramValue, "%d", &uintVal)
		if err != nil {
			panic(fmt.Sprintf("invalid uint16 parameter: %s", paramValue))
		}
		return uintVal
	case reflect.Uint32:
		var uintVal uint32
		_, err := fmt.Sscanf(paramValue, "%d", &uintVal)
		if err != nil {
			panic(fmt.Sprintf("invalid uint32 parameter: %s", paramValue))
		}
		return uintVal
	case reflect.Uint64:
		var uintVal uint64
		_, err := fmt.Sscanf(paramValue, "%d", &uintVal)
		if err != nil {
			panic(fmt.Sprintf("invalid uint64 parameter: %s", paramValue))
		}
		return uintVal
	case reflect.Float32:
		var floatVal float32
		_, err := fmt.Sscanf(paramValue, "%f", &floatVal)
		if err != nil {
			panic(fmt.Sprintf("invalid float parameter: %s", paramValue))
		}
		return floatVal
	case reflect.Float64:
		var floatVal float64
		_, err := fmt.Sscanf(paramValue, "%f", &floatVal)
		if err != nil {
			panic(fmt.Sprintf("invalid float parameter: %s", paramValue))
		}
		return floatVal
	case reflect.String:
		return paramValue
	case reflect.Interface:
		return paramValue
	default:
		panic(fmt.Sprintf("unsupported parameter type: %s", paramType.Kind()))
	}
}

// TagRule represents a parsed validation rule
type TagRule struct {
	Name   string
	Params []string
}

// parseTag parses a validation tag into rules
// Example: "required|min=5,max=100" -> [{required, []}, {min, ["5"]}, {max, ["100"]}]
// For multi-param: "between=10:100" -> [{between, ["10","100"]}]
// The params are split here into a slice
func parseTag(tag string, cfg *ParseConfig) []TagRule {
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

		if ch == byte(cfg.nameParamSeparator) {
			inParam = true
			currentRule += string(ch)
		} else if ch == byte(cfg.ruleSplitter) {
			// Check if we're in a parameter value
			// Look ahead to see if there's another '=' before the next ','
			if inParam {
				// Check if the next part looks like a new rule (contains '=') or is just a param
				nextPart := ""
				for j := i + 1; j < len(tag); j++ {
					if tag[j] == byte(cfg.ruleSplitter) {
						break
					}
					nextPart += string(tag[j])
				}

				// If nextPart doesn't contain '=' and doesn't look like a rule name, it's a parameter
				if !slices.Contains([]byte(nextPart), byte(cfg.nameParamSeparator)) && !isValidatorName(nextPart) {
					// This comma is part of the parameter
					currentRule += string(ch)
				} else {
					// This comma ends the current rule
					inParam = false
					if currentRule != "" {
						rules = append(rules, parseRule(currentRule, cfg))
						currentRule = ""
					}
				}
			} else {
				// Not in a parameter, this comma separates rules
				if currentRule != "" {
					rules = append(rules, parseRule(currentRule, cfg))
					currentRule = ""
				}
			}
		} else {
			currentRule += string(ch)
		}
	}

	// Don't forget the last rule
	if currentRule != "" {
		rules = append(rules, parseRule(currentRule, cfg))
	}

	return rules
}

// parseRule parses a single rule string like "min=5" or "required"
func parseRule(ruleStr string, cfg *ParseConfig) TagRule {
	ruleStr = strings.TrimSpace(ruleStr)

	if before, after, ok := strings.Cut(ruleStr, string(cfg.nameParamSeparator)); ok {
		name := strings.TrimSpace(before)
		raw := strings.TrimSpace(after)
		// split params by paramsSeparator into slice
		parts := []string{}
		if raw != "" {
			for _, p := range strings.Split(raw, string(cfg.paramsSeparator)) {
				tp := strings.TrimSpace(p)
				if tp != "" {
					parts = append(parts, tp)
				}
			}
		}
		return TagRule{
			Name:   name,
			Params: parts,
		}
	}

	return TagRule{
		Name:   ruleStr,
		Params: []string{},
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
