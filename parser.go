package validator

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/weilence/schema-validator/rule"
	"github.com/weilence/schema-validator/schema"
	"github.com/weilence/schema-validator/tag"
)

type ParseConfig struct {
	Registry   *rule.Registry
	TagParser  *tag.Parser
	DiveTag    string
	ValueTypes []reflect.Type
}

func defaultParseConfig() *ParseConfig {
	return &ParseConfig{
		Registry:   rule.DefaultRegistry(),
		TagParser:  tag.NewParser(tag.DefaultConfig()),
		DiveTag:    "dive",
		ValueTypes: []reflect.Type{reflect.TypeFor[time.Time]()},
	}
}

type ParseOption func(*ParseConfig)

func WithRegistry(registry *rule.Registry) ParseOption {
	return func(cfg *ParseConfig) {
		cfg.Registry = registry
	}
}

func WithTagParser(parser *tag.Parser) ParseOption {
	return func(cfg *ParseConfig) {
		cfg.TagParser = parser
	}
}

func WithDiveTag(diveTag string) ParseOption {
	return func(cfg *ParseConfig) {
		cfg.DiveTag = diveTag
	}
}

func WithValueTypes(types ...reflect.Type) ParseOption {
	return func(cfg *ParseConfig) {
		cfg.ValueTypes = append(cfg.ValueTypes, types...)
	}
}

// Parse parses a struct type into an ObjectSchema using struct tags
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

func parseStructField(s *schema.ObjectSchema, field reflect.StructField, cfg *ParseConfig) error {
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
			if err := parseStructField(s, embeddedField, cfg); err != nil {
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

	rules := cfg.TagParser.Parse(validateTag)
	fieldSchema, err := parseField(field.Type, rules, cfg)
	if err != nil {
		return err
	}

	s.AddField(fieldName, fieldSchema).AddFieldName(fieldName, field.Name)
	return nil
}

func parseField(fieldType reflect.Type, rules []tag.Rule, cfg *ParseConfig) (schema.Schema, error) {
	if fieldType.Kind() == reflect.Pointer {
		fieldType = fieldType.Elem()
	}

	if fieldType.Kind() == reflect.Slice || fieldType.Kind() == reflect.Array {
		diveIdx := slices.IndexFunc(rules, func(r tag.Rule) bool { return r.Name == cfg.DiveTag })
		var arrayRules, itemRules []tag.Rule
		if diveIdx >= 0 {
			arrayRules = rules[:diveIdx]
			itemRules = rules[diveIdx+1:]
		} else {
			arrayRules = rules
		}

		elemSchema, err := parseField(fieldType.Elem(), itemRules, cfg)
		if err != nil {
			return nil, err
		}

		arraySchema := schema.NewArray(elemSchema)
		for _, rule := range arrayRules {
			params, err := convertValidatorParams(rule.Name, rule.Params, cfg)
			if err != nil {
				return nil, err
			}
			v, err := cfg.Registry.NewValidator(rule.Name, params...)
			if err != nil {
				return nil, err
			}

			arraySchema.AddValidator(v)
		}

		return arraySchema, nil
	}

	if fieldType.Kind() == reflect.Struct && !slices.Contains(cfg.ValueTypes, fieldType) {
		return parse(fieldType, cfg)
	}

	if fieldType.Kind() == reflect.Map {
		return schema.NewObject(), nil
	}

	fieldSchema := schema.NewField()
	for _, rule := range rules {
		params, err := convertValidatorParams(rule.Name, rule.Params, cfg)
		if err != nil {
			return nil, err
		}
		v, err := cfg.Registry.NewValidator(rule.Name, params...)
		if err != nil {
			return nil, err
		}

		fieldSchema.AddValidator(v)
	}

	return fieldSchema, nil
}

func convertValidatorParams(name string, paramStrs []string, cfg *ParseConfig) ([]any, error) {
	paramTypes, err := cfg.Registry.GetValidatorParamTypes(name)
	if err != nil {
		return nil, err
	}

	paramTypesLen := len(paramTypes)
	if paramTypesLen == 0 && len(paramStrs) != 0 {
		return nil, fmt.Errorf("%s does not take any parameters", name)
	}

	if paramTypesLen == 1 {
		paramType := paramTypes[0]
		kind := paramType.Kind()
		switch kind {
		case reflect.Array:
			res := reflect.ArrayOf(paramType.Len(), paramType.Elem())
			rv := reflect.New(res).Elem()
			for i, paramStr := range paramStrs {
				elem, err := parseValidatorParam(paramType.Elem(), paramStr)
				if err != nil {
					return nil, err
				}
				rv.Index(i).Set(reflect.ValueOf(elem))
			}
			return []any{rv.Interface()}, nil
		case reflect.Slice:
			res := reflect.MakeSlice(paramType, 0, 0)
			for _, paramStr := range paramStrs {
				elem, err := parseValidatorParam(paramType.Elem(), paramStr)
				if err != nil {
					return nil, err
				}
				res = reflect.Append(res, reflect.ValueOf(elem))
			}
			return []any{res.Interface()}, nil
		default:
			if len(paramStrs) != 1 {
				return nil, fmt.Errorf("%s expected 1 parameter, got %d", name, len(paramStrs))
			}
			param, err := parseValidatorParam(paramType, paramStrs[0])
			if err != nil {
				return nil, err
			}
			return []any{param}, nil
		}
	}

	if len(paramStrs) != paramTypesLen {
		return nil, fmt.Errorf("%s expected %d parameters, got %d", name, paramTypesLen, len(paramStrs))
	}

	params := make([]any, paramTypesLen)
	for i, paramType := range paramTypes {
		param, err := parseValidatorParam(paramType, paramStrs[i])
		if err != nil {
			return nil, err
		}
		params[i] = param
	}

	return params, nil
}

func parseValidatorParam(paramType reflect.Type, paramValue string) (any, error) {
	switch paramType.Kind() {
	case reflect.Bool:
		var v bool
		if _, err := fmt.Sscanf(paramValue, "%t", &v); err != nil {
			return nil, fmt.Errorf("invalid bool parameter: %s", paramValue)
		}
		return v, nil
	case reflect.Int:
		var v int
		if _, err := fmt.Sscanf(paramValue, "%d", &v); err != nil {
			return nil, fmt.Errorf("invalid int parameter: %s", paramValue)
		}
		return v, nil
	case reflect.Int8:
		var v int8
		if _, err := fmt.Sscanf(paramValue, "%d", &v); err != nil {
			return nil, fmt.Errorf("invalid int8 parameter: %s", paramValue)
		}
		return v, nil
	case reflect.Int16:
		var v int16
		if _, err := fmt.Sscanf(paramValue, "%d", &v); err != nil {
			return nil, fmt.Errorf("invalid int16 parameter: %s", paramValue)
		}
		return v, nil
	case reflect.Int32:
		var v int32
		if _, err := fmt.Sscanf(paramValue, "%d", &v); err != nil {
			return nil, fmt.Errorf("invalid int32 parameter: %s", paramValue)
		}
		return v, nil
	case reflect.Int64:
		var v int64
		if _, err := fmt.Sscanf(paramValue, "%d", &v); err != nil {
			return nil, fmt.Errorf("invalid int64 parameter: %s", paramValue)
		}
		return v, nil
	case reflect.Uint:
		var v uint
		if _, err := fmt.Sscanf(paramValue, "%d", &v); err != nil {
			return nil, fmt.Errorf("invalid uint parameter: %s", paramValue)
		}
		return v, nil
	case reflect.Uint8:
		var v uint8
		if _, err := fmt.Sscanf(paramValue, "%d", &v); err != nil {
			return nil, fmt.Errorf("invalid uint8 parameter: %s", paramValue)
		}
		return v, nil
	case reflect.Uint16:
		var v uint16
		if _, err := fmt.Sscanf(paramValue, "%d", &v); err != nil {
			return nil, fmt.Errorf("invalid uint16 parameter: %s", paramValue)
		}
		return v, nil
	case reflect.Uint32:
		var v uint32
		if _, err := fmt.Sscanf(paramValue, "%d", &v); err != nil {
			return nil, fmt.Errorf("invalid uint32 parameter: %s", paramValue)
		}
		return v, nil
	case reflect.Uint64:
		var v uint64
		if _, err := fmt.Sscanf(paramValue, "%d", &v); err != nil {
			return nil, fmt.Errorf("invalid uint64 parameter: %s", paramValue)
		}
		return v, nil
	case reflect.Float32:
		var v float32
		if _, err := fmt.Sscanf(paramValue, "%f", &v); err != nil {
			return nil, fmt.Errorf("invalid float32 parameter: %s", paramValue)
		}
		return v, nil
	case reflect.Float64:
		var v float64
		if _, err := fmt.Sscanf(paramValue, "%f", &v); err != nil {
			return nil, fmt.Errorf("invalid float64 parameter: %s", paramValue)
		}
		return v, nil
	case reflect.String, reflect.Interface:
		return paramValue, nil
	default:
		return nil, fmt.Errorf("unsupported parameter type: %s", paramType.Kind())
	}
}

func getFieldName(field reflect.StructField) string {
	if name := extractNameFromTag(field.Tag.Get("json")); name != "" {
		return name
	}
	if name := extractNameFromTag(field.Tag.Get("param")); name != "" {
		return name
	}
	if name := extractNameFromTag(field.Tag.Get("query")); name != "" {
		return name
	}
	return field.Name
}

func extractNameFromTag(t string) string {
	if t == "" || t == "-" {
		return ""
	}
	if before, _, ok := strings.Cut(t, ","); ok {
		return before
	}
	return t
}
