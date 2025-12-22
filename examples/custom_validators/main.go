package main

import (
	"fmt"
	"reflect"

	validator "github.com/weilence/schema-validator"
	"github.com/weilence/schema-validator/errors"
	"github.com/weilence/schema-validator/tags"
	"github.com/weilence/schema-validator/validation"
)

// parseInt 辅助函数
func parseInt(s string) int {
	var result int
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			result = result*10 + int(ch-'0')
		}
	}
	return result
}

func main() {
	fmt.Println("=== Multi-Parameter Validators Examples ===\n")

	// 创建自定义registry
	registry := tags.NewRegistry()

	// Example 1: between validator (2 parameters)
	fmt.Println("Example 1: Between Validator (min:max)")
	fmt.Println("---------------------------------------")

	registry.RegisterField("between", func(params []string) (validation.FieldValidator, error) {
		if len(params) < 2 {
			return nil, nil
		}

		min := parseInt(params[0])
		max := parseInt(params[1])

		return validation.FieldValidatorFunc(func(ctx *validation.Context) error {
			field, _ := ctx.AsField()
			val, _ := field.Int()
			if val < int64(min) || val > int64(max) {
				return errors.NewValidationError(ctx.Path(), "between", map[string]interface{}{
					"min": min, "max": max, "actual": val,
				})
			}
			return nil
		}), nil
	})

	type Product struct {
		Price int `json:"price" validate:"between=10:100"`
	}

	typ := reflect.TypeOf(Product{})
	objSchema, _ := tags.ParseStructTagsWithRegistry(typ, registry)
	v := validator.New(objSchema)

	// Valid price
	product1 := Product{Price: 50}
	result, _ := v.Validate(product1)
	fmt.Printf("Product with price=50 (valid): %v\n", result.IsValid())

	// Too low
	product2 := Product{Price: 5}
	result, _ = v.Validate(product2)
	fmt.Printf("Product with price=5 (too low): %v\n", result.IsValid())
	if !result.IsValid() {
		for _, err := range result.Errors() {
			fmt.Printf("  - %s: %s (params: %v)\n", err.FieldPath, err.ErrorCode, err.Params)
		}
	}

	// Too high
	product3 := Product{Price: 150}
	result, _ = v.Validate(product3)
	fmt.Printf("Product with price=150 (too high): %v\n\n", result.IsValid())

	// Example 2: enum validator (multiple parameters)
	fmt.Println("Example 2: Enum Validator (value1:value2:...)")
	fmt.Println("----------------------------------------------")

	registry.RegisterField("enum", func(params []string) (validation.FieldValidator, error) {
		allowedValues := make(map[string]bool)
		for _, p := range params {
			allowedValues[p] = true
		}

		return validation.FieldValidatorFunc(func(ctx *validation.Context) error {
			field, _ := ctx.AsField()
			val := field.String()
			if !allowedValues[val] {
				return errors.NewValidationError(ctx.Path(), "enum", map[string]interface{}{
					"allowed": params,
					"actual":  val,
				})
			}
			return nil
		}), nil
	})

	type Settings struct {
		Theme string `json:"theme" validate:"enum=light:dark:auto"`
	}

	typ2 := reflect.TypeOf(Settings{})
	objSchema2, _ := tags.ParseStructTagsWithRegistry(typ2, registry)
	v2 := validator.New(objSchema2)

	// Valid theme
	settings1 := Settings{Theme: "dark"}
	result, _ = v2.Validate(settings1)
	fmt.Printf("Settings with theme=dark (valid): %v\n", result.IsValid())

	// Invalid theme
	settings2 := Settings{Theme: "blue"}
	result, _ = v2.Validate(settings2)
	fmt.Printf("Settings with theme=blue (invalid): %v\n", result.IsValid())
	if !result.IsValid() {
		for _, err := range result.Errors() {
			fmt.Printf("  - %s: %s (params: %v)\n", err.FieldPath, err.ErrorCode, err.Params)
		}
	}
	fmt.Println()

	// Example 3: range validator (3 parameters: min:max:step)
	fmt.Println("Example 3: Range Validator (min:max:step)")
	fmt.Println("-----------------------------------------")

	registry.RegisterField("range", func(params []string) (validation.FieldValidator, error) {
		if len(params) < 3 {
			return nil, nil
		}

		min := parseInt(params[0])
		max := parseInt(params[1])
		step := parseInt(params[2])

		return validation.FieldValidatorFunc(func(ctx *validation.Context) error {
			field, _ := ctx.AsField()
			val, _ := field.Int()

			if val < int64(min) || val > int64(max) {
				return errors.NewValidationError(ctx.Path(), "out_of_range", map[string]interface{}{
					"min": min, "max": max, "actual": val,
				})
			}

			if (int(val)-min)%step != 0 {
				return errors.NewValidationError(ctx.Path(), "invalid_step", map[string]interface{}{
					"step": step, "actual": val,
				})
			}

			return nil
		}), nil
	})

	type SliderValue struct {
		Volume int `json:"volume" validate:"range=0:100:5"`
	}

	typ3 := reflect.TypeOf(SliderValue{})
	objSchema3, _ := tags.ParseStructTagsWithRegistry(typ3, registry)
	v3 := validator.New(objSchema3)

	// Valid values (multiples of 5)
	slider1 := SliderValue{Volume: 50}
	result, _ = v3.Validate(slider1)
	fmt.Printf("Volume=50 (valid, multiple of 5): %v\n", result.IsValid())

	slider2 := SliderValue{Volume: 0}
	result, _ = v3.Validate(slider2)
	fmt.Printf("Volume=0 (valid, at min): %v\n", result.IsValid())

	// Invalid step
	slider3 := SliderValue{Volume: 47}
	result, _ = v3.Validate(slider3)
	fmt.Printf("Volume=47 (invalid, not multiple of 5): %v\n", result.IsValid())
	if !result.IsValid() {
		for _, err := range result.Errors() {
			fmt.Printf("  - %s: %s (params: %v)\n", err.FieldPath, err.ErrorCode, err.Params)
		}
	}

	// Out of range
	slider4 := SliderValue{Volume: 105}
	result, _ = v3.Validate(slider4)
	fmt.Printf("Volume=105 (out of range): %v\n", result.IsValid())
	if !result.IsValid() {
		for _, err := range result.Errors() {
			fmt.Printf("  - %s: %s (params: %v)\n", err.FieldPath, err.ErrorCode, err.Params)
		}
	}

	fmt.Println("\n=== Multi-Parameter Validators Examples Completed ===")
}
