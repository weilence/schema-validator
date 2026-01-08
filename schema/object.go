package schema

import (
	"fmt"
	"reflect"

	"github.com/weilence/schema-validator/data"
)

// ObjectSchema validates objects/structs/maps
type ObjectSchema struct {
	fields       map[string]Schema
	fieldNameMap map[string]string // mapping of lower-case field names to actual names
	validators   []Validator
}

// NewObject creates a new object schema
func NewObject() *ObjectSchema {
	return &ObjectSchema{
		fields:       make(map[string]Schema),
		fieldNameMap: make(map[string]string),
		validators:   make([]Validator, 0),
	}
}

// Validate validates an object
func (o *ObjectSchema) Validate(ctx *Context) error {
	accessor := ctx.Accessor()
	switch oa := accessor.(type) {
	case data.ObjectAccessor:
		for _, accessor := range oa.Accessors() {
			if v, ok := accessor.Raw().(SchemaModifier); ok {
				v.ModifySchema(ctx)
			}
		}
	case *data.Value:
		if oa.Kind() == reflect.Invalid {
			return nil
		}
		if oa.Kind() == reflect.Ptr && oa.IsNilOrZero() {
			return nil
		}

		return fmt.Errorf("expected object accessor, got primitive value")
	default:
		return fmt.Errorf("expected object accessor, got %T", oa)
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

		fieldData, err := accessor.GetField(fieldName)
		if err != nil {
			return fmt.Errorf("error accessing field %s: %w", fieldName, err)
		}

		fieldCtx := ctx.WithChild(name, fieldSchema, fieldData)

		if err := fieldSchema.Validate(fieldCtx); err != nil {
			return err
		}
	}

	return nil
}

func (o *ObjectSchema) Field(name string) Schema {
	return o.fields[name]
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
		s.element = mergeSchema(s.element, as2.element)
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
