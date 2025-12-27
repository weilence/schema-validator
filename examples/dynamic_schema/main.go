package main

import (
	"fmt"

	validator "github.com/weilence/schema-validator"
	"github.com/weilence/schema-validator/builder"
	"github.com/weilence/schema-validator/schema"
)

// Example 1: Conditional required fields
type DynamicForm struct {
	Type     string `json:"type" validate:"required"`
	Value    string `json:"value"`
	Required bool   `json:"required"`
}

func (f DynamicForm) ModifySchema(ctx *schema.Context) {
	// 从 ctx 获取 ObjectSchema
	s, ok := ctx.Schema().(*schema.ObjectSchema)
	if !ok || s == nil {
		return
	}

	// 根据required标志动态修改value字段的验证
	if f.Required {
		s.AddField("value", builder.Field().
			AddValidator("required").
			Build())
	} else {
		s.RemoveField("value")
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

func (u User) ModifySchema(ctx *schema.Context) {
	s := ctx.Schema().(*schema.ObjectSchema)

	// 根据国家添加不同的邮编验证
	switch u.Address.Country {
	case "US":
		addressSchema := builder.Object().
			Field("zipCode", builder.Field().
				AddValidator("min", 5).
				AddValidator("max", 5).
				Build()).
			Build()
		s.AddField("address", addressSchema)
	case "CN":
		addressSchema := builder.Object().
			Field("zipCode", builder.Field().
				AddValidator("min", 6).
				AddValidator("max", 6).
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

	v1, _ := validator.New(DynamicForm{})

	// required=true时，value必填
	form1 := DynamicForm{Type: "text", Value: "", Required: true}
	err := v1.Validate(form1)
	fmt.Printf("Form with required=true, empty value: %v\n", err)
	if err != nil {
		newErr := err.(*schema.ValidationError)
		fmt.Printf("  - %s: %s\n", newErr.Path, newErr.Name)
	}

	// required=false时，value可选
	form2 := DynamicForm{Type: "text", Value: "", Required: false}
	err = v1.Validate(form2)
	fmt.Printf("Form with required=false, empty value: %v\n", err)

	// required=true时，value有值
	form3 := DynamicForm{Type: "text", Value: "hello", Required: true}
	err = v1.Validate(form3)
	fmt.Printf("Form with required=true, non-empty value: %v\n\n", err)

	// Example 2: Type-based validation
	fmt.Println("Example 2: Type-based Validation (Country-specific ZipCode)")
	fmt.Println("-------------------------------------------------------------")

	v2, _ := validator.New(User{})

	// US zipcode (5 digits)
	user1 := User{
		Name: "John",
		Address: Address{
			Country: "US",
			ZipCode: "12345",
		},
	}
	err = v2.Validate(user1)
	fmt.Printf("US user with valid zipcode (12345): %v\n", err)

	// US zipcode invalid (6 digits)
	user2 := User{
		Name: "John",
		Address: Address{
			Country: "US",
			ZipCode: "123456",
		},
	}
	err = v2.Validate(user2)
	fmt.Printf("US user with invalid zipcode (123456): %v\n", err)
	if err != nil {
		fmt.Println(err)
	}

	// China zipcode (6 digits)
	user3 := User{
		Name: "李明",
		Address: Address{
			Country: "CN",
			ZipCode: "100000",
		},
	}
	err = v2.Validate(user3)
	fmt.Printf("CN user with valid zipcode (100000): %v\n", err)

	// China zipcode invalid (5 digits)
	user4 := User{
		Name: "李明",
		Address: Address{
			Country: "CN",
			ZipCode: "10000",
		},
	}
	err = v2.Validate(user4)
	fmt.Printf("CN user with invalid zipcode (10000): %v\n", err)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("\n=== Dynamic Schema Examples Completed ===")
}
