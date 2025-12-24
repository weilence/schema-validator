package tags

import (
	"reflect"
	"slices"
	"strings"

	"github.com/weilence/schema-validator/schema"
)

type ParseConfig struct {
	registry *schema.Registry

	ruleSplitter       rune
	nameParamSeparator rune
	paramsSeparator    rune
	diveTag            string
}

func defaultParseConfig() *ParseConfig {
	return &ParseConfig{
		registry:           schema.DefaultRegistry(),
		ruleSplitter:       '|',
		nameParamSeparator: '=',
		paramsSeparator:    ',',
		diveTag:            "dive",
	}
}

type ParseOption func(*ParseConfig)

func WithRegistry(registry *schema.Registry) ParseOption {
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

	objSchema := schema.NewObjectSchema()
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

		arraySchema := schema.NewArraySchema(elemSchema)

		for _, rule := range arrayTags {
			arraySchema.AddValidator(rule.Name, rule.Params...)
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
		return schema.NewObjectSchema(), nil
	}

	// Primitive field
	fieldSchema := schema.NewFieldSchema()
	for _, tag := range tags {
		fieldSchema.AddValidator(tag.Name, tag.Params...)
	}

	return fieldSchema, nil
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
