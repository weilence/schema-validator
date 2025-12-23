package schema

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/weilence/schema-validator/data"
)

// ObjectSchema validates objects/structs/maps
type ObjectSchema struct {
	fields       map[string]Schema
	fieldNameMap map[string]string // mapping of lower-case field names to actual names
	validators   []Validator
	strict       bool // disallow unknown fields
}

// NewObjectSchema creates a new object schema
func NewObjectSchema() *ObjectSchema {
	return &ObjectSchema{
		fields:       make(map[string]Schema),
		fieldNameMap: make(map[string]string),
		validators:   make([]Validator, 0),
		strict:       false,
	}
}

// Type returns SchemaTypeObject
func (o *ObjectSchema) Type() SchemaType {
	return SchemaTypeObject
}

// Validate validates an object
func (o *ObjectSchema) Validate(ctx *Context) error {
	oa := ctx.Accessor().(data.ObjectAccessor)
	for _, accessor := range oa.Accessors() {
		if v, ok := accessor.Raw().(SchemaModifier); ok {
			v.ModifySchema(ctx)
		}
	}

	for name, fieldSchema := range o.fields {
		fieldName := name
		if mappedName, ok := o.fieldNameMap[name]; ok {
			fieldName = mappedName
		}

		fieldData, err := oa.GetField(fieldName)
		if err != nil {
			return fmt.Errorf("error accessing field %s: %w", fieldName, err)
		}

		// 创建子 context - 自动父级追踪和 schema 传递
		fieldCtx := ctx.WithChild(name, fieldSchema, fieldData)

		if err := fieldSchema.Validate(fieldCtx); err != nil {
			return err
		}
	}

	// Run object-level validators (cross-field)
	for _, validator := range o.validators {
		if err := validator.Validate(ctx); err != nil {
			return err
		}
	}

	return nil
}

// AddField adds a field schema
func (o *ObjectSchema) AddField(name string, schema Schema) *ObjectSchema {
	// TODO: process duplicate field names
	if _, ok := o.fields[name]; ok {
		slog.Warn("overwriting existing field schema", "field", name)
	}

	o.fields[name] = schema
	return o
}

func (o *ObjectSchema) RemoveField(name string) *ObjectSchema {
	delete(o.fields, name)
	return o
}

func (o *ObjectSchema) AddFieldName(name string, fieldName string) *ObjectSchema {
	o.fieldNameMap[name] = fieldName
	return o
}

// AddValidator adds an object validator
func (o *ObjectSchema) AddValidator(name string, params ...string) *ObjectSchema {
	v := DefaultRegistry().BuildValidator(name, params)
	if v != nil {
		o.validators = append(o.validators, v)
	}
	return o
}

// GetField returns the schema for a field
func (o *ObjectSchema) GetField(name string) (Schema, bool) {
	schema, ok := o.fields[name]
	return schema, ok
}

// Fields returns all field names
func (o *ObjectSchema) Fields() []string {
	names := make([]string, 0, len(o.fields))
	for name := range o.fields {
		names = append(names, name)
	}
	return names
}

// ToString returns a JSON representation of the object schema
func (o *ObjectSchema) ToString() string {
	result := map[string]any{
		"type":   "object",
		"strict": o.strict,
	}

	// Add fields as nested schemas
	if len(o.fields) > 0 {
		fields := make(map[string]any)
		for fieldName, fieldSchema := range o.fields {
			// Parse the nested schema's JSON string back to a map for proper nesting
			var fieldMap map[string]any
			fieldJSON := fieldSchema.ToString()
			if err := json.Unmarshal([]byte(fieldJSON), &fieldMap); err == nil {
				fields[fieldName] = fieldMap
			} else {
				fields[fieldName] = fieldJSON
			}
		}
		result["fields"] = fields
	}

	if len(o.validators) > 0 {
		validators := make([]map[string]any, 0, len(o.validators))
		for _, v := range o.validators {
			validators = append(validators, validatorToMap(v))
		}
		result["validators"] = validators
	}

	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"type":"object","error":"%s"}`, err.Error())
	}
	return string(bytes)
}

func SchemaModifiers(r any) []SchemaModifier {
	var modifiers []SchemaModifier
	if r == nil {
		return modifiers
	}

	rval := reflect.ValueOf(r)
	for rval.Kind() == reflect.Ptr {
		if rval.IsNil() {
			return modifiers
		}
		rval = rval.Elem()
	}

	if schemaModifier, ok := rval.Interface().(SchemaModifier); ok {
		modifiers = append(modifiers, schemaModifier)
	}

	// Iterate struct fields and collect SchemaModifier from embedded fields
	for i := 0; i < rval.NumField(); i++ {
		sf := rval.Type().Field(i)
		// Only consider anonymous (embedded) fields
		if !sf.Anonymous {
			continue
		}

		fieldVal := rval.Field(i)
		if !fieldVal.IsValid() {
			continue
		}

		// If the field cannot be interfaced, skip it
		if !fieldVal.CanInterface() {
			continue
		}

		emods := SchemaModifiers(fieldVal.Interface())
		if len(emods) > 0 {
			modifiers = append(modifiers, emods...)
		}
	}

	return modifiers
}
