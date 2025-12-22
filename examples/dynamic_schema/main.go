package main

import (
	"fmt"

	validator "github.com/weilence/schema-validator"
	"github.com/weilence/schema-validator/schema"
	"github.com/weilence/schema-validator/validation"
)

// Example 1: Conditional required fields
type DynamicForm struct {
	Type     string `json:"type" validate:"required"`
	Value    string `json:"value"`
	Required bool   `json:"required"`
}

func (f DynamicForm) ModifySchema(ctx *validation.Context) {
	// 从 ctx 获取 ObjectSchema
	s, ok := ctx.ObjectSchema().(*schema.ObjectSchema)
	if !ok || s == nil {
		return
	}

	// 从 ctx 获取 accessor（使用缓存）
	obj, _ := ctx.AsObject()

	// 读取required字段的值
	requiredField, _ := obj.GetField("required")
	fieldAcc, _ := requiredField.AsField()
	isRequired, _ := fieldAcc.Bool()

	// 根据required标志动态修改value字段的验证
	if isRequired {
		s.AddField("value", schema.Field().
			AddValidator(validation.Required()).
			Build())
	} else {
		s.AddField("value", schema.Field().
			SetOptional(true).
			Build())
	}
}

// Example 2: Type-based validation
type User struct {
	Name    string  `json:"name" validate:"required"`
	Address Address `json:"address"`
}

type Address struct {
	Country string `json:"country"`
	ZipCode string `json:"zipCode"`
}

func (u User) ModifySchema(ctx *validation.Context) {
	s, ok := ctx.ObjectSchema().(*schema.ObjectSchema)
	if !ok || s == nil {
		return
	}

	obj, _ := ctx.AsObject()

	// 访问嵌套的address.country
	addressField, _ := obj.GetField("address")
	addressObj, _ := addressField.AsObject()
	countryField, _ := addressObj.GetField("country")
	countryAcc, _ := countryField.AsField()
	country := countryAcc.String()

	// 根据国家添加不同的邮编验证
	switch country {
	case "US":
		addressSchema := schema.Object().
			Field("country", schema.Field().Build()).
			Field("zipCode", schema.Field().
				AddValidator(validation.MinLength(5)).
				AddValidator(validation.MaxLength(5)).
				Build()).
			Build()
		s.AddField("address", addressSchema)
	case "CN":
		addressSchema := schema.Object().
			Field("country", schema.Field().Build()).
			Field("zipCode", schema.Field().
				AddValidator(validation.MinLength(6)).
				AddValidator(validation.MaxLength(6)).
				Build()).
			Build()
		s.AddField("address", addressSchema)
	}
}

func main() {
	fmt.Println("=== Dynamic Schema (SchemaModifier) Examples ===")

	// Example 1: Conditional required fields
	fmt.Println("Example 1: Conditional Required Fields")
	fmt.Println("---------------------------------------")

	v1, _ := validator.NewFromStruct(DynamicForm{})

	// required=true时，value必填
	form1 := DynamicForm{Type: "text", Value: "", Required: true}
	result, _ := v1.Validate(form1)
	fmt.Printf("Form with required=true, empty value: %v\n", result.IsValid())
	if !result.IsValid() {
		for _, err := range result.Errors() {
			fmt.Printf("  - %s: %s\n", err.FieldPath, err.ErrorCode)
		}
	}

	// required=false时，value可选
	form2 := DynamicForm{Type: "text", Value: "", Required: false}
	result, _ = v1.Validate(form2)
	fmt.Printf("Form with required=false, empty value: %v\n", result.IsValid())

	// required=true时，value有值
	form3 := DynamicForm{Type: "text", Value: "hello", Required: true}
	result, _ = v1.Validate(form3)
	fmt.Printf("Form with required=true, non-empty value: %v\n\n", result.IsValid())

	// Example 2: Type-based validation
	fmt.Println("Example 2: Type-based Validation (Country-specific ZipCode)")
	fmt.Println("-------------------------------------------------------------")

	v2, _ := validator.NewFromStruct(User{})

	// US zipcode (5 digits)
	user1 := User{
		Name: "John",
		Address: Address{
			Country: "US",
			ZipCode: "12345",
		},
	}
	result, _ = v2.Validate(user1)
	fmt.Printf("US user with valid zipcode (12345): %v\n", result.IsValid())

	// US zipcode invalid (6 digits)
	user2 := User{
		Name: "John",
		Address: Address{
			Country: "US",
			ZipCode: "123456",
		},
	}
	result, _ = v2.Validate(user2)
	fmt.Printf("US user with invalid zipcode (123456): %v\n", result.IsValid())
	if !result.IsValid() {
		for _, err := range result.Errors() {
			fmt.Printf("  - %s: %s (params: %v)\n", err.FieldPath, err.ErrorCode, err.Params)
		}
	}

	// China zipcode (6 digits)
	user3 := User{
		Name: "李明",
		Address: Address{
			Country: "CN",
			ZipCode: "100000",
		},
	}
	result, _ = v2.Validate(user3)
	fmt.Printf("CN user with valid zipcode (100000): %v\n", result.IsValid())

	// China zipcode invalid (5 digits)
	user4 := User{
		Name: "李明",
		Address: Address{
			Country: "CN",
			ZipCode: "10000",
		},
	}
	result, _ = v2.Validate(user4)
	fmt.Printf("CN user with invalid zipcode (10000): %v\n", result.IsValid())
	if !result.IsValid() {
		for _, err := range result.Errors() {
			fmt.Printf("  - %s: %s (params: %v)\n", err.FieldPath, err.ErrorCode, err.Params)
		}
	}

	fmt.Println("\n=== Dynamic Schema Examples Completed ===")
}
