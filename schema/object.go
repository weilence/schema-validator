package schema

import (
	"fmt"

	"github.com/weilence/schema-validator/data"
)

// ObjectSchema validates objects/structs/maps
type ObjectSchema struct {
	fields       map[string]Schema
	fieldNameMap map[string]string // mapping of lower-case field names to actual names
	validators   []Validator
}

// NewObjectSchema creates a new object schema
func NewObjectSchema() *ObjectSchema {
	return &ObjectSchema{
		fields:       make(map[string]Schema),
		fieldNameMap: make(map[string]string),
		validators:   make([]Validator, 0),
	}
}

// Validate validates an object
func (o *ObjectSchema) Validate(ctx *Context) error {
	oa := ctx.Accessor().(data.ObjectAccessor)
	for _, accessor := range oa.Accessors() {
		if v, ok := accessor.Raw().(SchemaModifier); ok {
			v.ModifySchema(ctx)
		}
	}

	for _, validator := range o.validators {
		if ctx.skipRest {
			break
		}

		if err := validator.Validate(ctx); err != nil {
			return err
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

	return nil
}

// AddField adds a field schema
func (o *ObjectSchema) AddField(name string, schema Schema) *ObjectSchema {
	if oldSchema, ok := o.fields[name]; ok {
		o.fields[name] = mergeSchema(oldSchema, schema)
	} else {
		o.fields[name] = schema
	}

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

func (o *ObjectSchema) AddValidator(v Validator) Schema {
	o.validators = append(o.validators, v)
	return o
}

func (o *ObjectSchema) RemoveValidator(name string) Schema {
	newValidators := make([]Validator, 0)
	for _, v := range o.validators {
		if v.Name() != name {
			newValidators = append(newValidators, v)
		}
	}
	o.validators = newValidators
	return o
}

func mergeSchema(s1, s2 Schema) Schema {
	switch s := s1.(type) {
	case *FieldSchema:
		fs2 := s2.(*FieldSchema)
		for _, v := range fs2.validators {
			s.AddValidator(v)
		}
		return s
	case *ArraySchema:
		as2 := s2.(*ArraySchema)
		s.elementSchema = mergeSchema(s.elementSchema, as2.elementSchema)
		for _, v := range as2.validators {
			s.AddValidator(v)
		}
		return s
	case *ObjectSchema:
		os2 := s2.(*ObjectSchema)
		for name, fieldSchema2 := range os2.fields {
			if fieldSchema1, ok := s.fields[name]; ok {
				s.fields[name] = mergeSchema(fieldSchema1, fieldSchema2)
			} else {
				s.fields[name] = fieldSchema2
			}
		}
		for _, v := range os2.validators {
			s.AddValidator(v)
		}
		return s
	default:
		panic("unknown schema type")
	}
}
