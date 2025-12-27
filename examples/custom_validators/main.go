package main

import (
	"fmt"
	"reflect"

	validator "github.com/weilence/schema-validator"
	"github.com/weilence/schema-validator/builder"
	"github.com/weilence/schema-validator/schema"
	"github.com/weilence/schema-validator/validators"
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
	fmt.Println("=== Multi-Parameter Validators Examples ===")

	// 创建自定义registry
	registry := validators.NewRegistry()

	// Example 1: between validator (2 parameters)
	fmt.Println("Example 1: Between Validator (min:max)")
	fmt.Println("---------------------------------------")

	registry.Register("between", func(ctx *schema.Context, params []string) error {
		if len(params) < 2 {
			return nil
		}

		min := parseInt(params[0])
		max := parseInt(params[1])

		fieldVal := ctx.Value()
		val := fieldVal.Int()
		if val < int64(min) || val > int64(max) {
			return schema.NewValidationError(ctx.Path(), "between", map[string]any{
				"min": min, "max": max, "actual": val,
			})
		}
		return nil
	})

	type Product struct {
		Price int `json:"price" validate:"between=10:100"`
	}

	typ := reflect.TypeOf(Product{})
	objSchema, _ := builder.Parse(typ, builder.WithRegistry(registry))
	v := validator.NewFromSchema(objSchema)

	// Valid price
	product1 := Product{Price: 50}
	err := v.Validate(product1)
	fmt.Printf("Product with price=50 (valid): %v\n", err)

	// Too low
	product2 := Product{Price: 5}
	err = v.Validate(product2)
	fmt.Printf("Product with price=5 (too low): %v\n", err)
	if err != nil {
		for _, err := range err.(*schema.ValidationResult).Errors() {
			fmt.Printf("  - %s: %s (params: %v)\n", err.Path, err.Name, err.Params)
		}
	}

	// Too high
	product3 := Product{Price: 150}
	err = v.Validate(product3)
	fmt.Printf("Product with price=150 (too high): %v\n\n", err)

	// Example 2: enum validator (multiple parameters)
	fmt.Println("Example 2: Enum Validator (value1:value2:...)")
	fmt.Println("----------------------------------------------")

	registry.Register("enum", func(ctx *schema.Context, params []string) error {
		allowedValues := make(map[string]bool)
		for _, p := range params {
			allowedValues[p] = true
		}

		fieldVal := ctx.Value()
		val := fieldVal.String()
		if !allowedValues[val] {
			return schema.NewValidationError(ctx.Path(), "enum", map[string]any{
				"allowed": params,
				"actual":  val,
			})
		}
		return nil
	})

	type Settings struct {
		Theme string `json:"theme" validate:"enum=light:dark:auto"`
	}

	typ2 := reflect.TypeOf(Settings{})
	objSchema2, _ := builder.Parse(typ2, builder.WithRegistry(registry))
	v2 := validator.NewFromSchema(objSchema2)

	// Valid theme
	settings1 := Settings{Theme: "dark"}
	err = v2.Validate(settings1)
	fmt.Printf("Settings with theme=dark (valid): %v\n", err)

	// Invalid theme
	settings2 := Settings{Theme: "blue"}
	err = v2.Validate(settings2)
	fmt.Printf("Settings with theme=blue (invalid): %v\n", err)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()

	// Example 3: range validator (3 parameters: min:max:step)
	fmt.Println("Example 3: Range Validator (min:max:step)")
	fmt.Println("-----------------------------------------")

	registry.Register("range", func(ctx *schema.Context, params []string) error {
		if len(params) < 3 {
			return nil
		}

		min := parseInt(params[0])
		max := parseInt(params[1])
		step := parseInt(params[2])

		fieldVal := ctx.Value()
		val := fieldVal.Int()

		if val < int64(min) || val > int64(max) {
			return schema.NewValidationError(ctx.Path(), "out_of_range", map[string]any{
				"min": min, "max": max, "actual": val,
			})
		}

		if (int(val)-min)%step != 0 {
			return schema.NewValidationError(ctx.Path(), "invalid_step", map[string]any{
				"step": step, "actual": val,
			})
		}

		return nil
	})

	type SliderValue struct {
		Volume int `json:"volume" validate:"range=0:100:5"`
	}

	typ3 := reflect.TypeOf(SliderValue{})
	objSchema3, _ := builder.Parse(typ3, builder.WithRegistry(registry))
	v3 := validator.NewFromSchema(objSchema3)

	// Valid values (multiples of 5)
	slider1 := SliderValue{Volume: 50}
	err = v3.Validate(slider1)
	fmt.Printf("Volume=50 (valid, multiple of 5): %v\n", err)

	slider2 := SliderValue{Volume: 0}
	err = v3.Validate(slider2)
	fmt.Printf("Volume=0 (valid, at min): %v\n", err)

	// Invalid step
	slider3 := SliderValue{Volume: 47}
	err = v3.Validate(slider3)
	fmt.Printf("Volume=47 (invalid, not multiple of 5): %v\n", err)
	if err != nil {
		for _, err := range err.(*schema.ValidationResult).Errors() {
			fmt.Printf("  - %s: %s (params: %v)\n", err.Path, err.Name, err.Params)
		}
	}

	// Out of range
	slider4 := SliderValue{Volume: 105}
	err = v3.Validate(slider4)
	fmt.Printf("Volume=105 (out of range): %v\n", err)
	if err != nil {
		for _, err := range err.(*schema.ValidationResult).Errors() {
			fmt.Printf("  - %s: %s (params: %v)\n", err.Path, err.Name, err.Params)
		}
	}

	fmt.Println("\n=== Multi-Parameter Validators Examples Completed ===")
}
