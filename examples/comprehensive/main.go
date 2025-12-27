package main

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"
	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/validators"

	validator "github.com/weilence/schema-validator"
	"github.com/weilence/schema-validator/builder"
	"github.com/weilence/schema-validator/schema"
)

// OrderForm 订单表单，综合使用多个高级功能
type OrderForm struct {
	OrderType   string `json:"orderType" validate:"enum=standard,express,same_day"`
	Price       int    `json:"price" validate:"between=1,10000"`
	MinQuantity int    `json:"minQuantity"`
	MaxQuantity int    `json:"maxQuantity"`
	Quantity    int    `json:"quantity"`
}

// 实现SchemaModifier接口
func (f OrderForm) ModifySchema(ctx *schema.Context) {
	s, ok := ctx.Schema().(*schema.ObjectSchema)
	if !ok || s == nil {
		return
	}

	obj, _ := ctx.Accessor().(data.ObjectAccessor)

	// 1. 根据min/max动态设置quantity范围
	minField, _ := obj.GetField("MinQuantity")
	maxField, _ := obj.GetField("MaxQuantity")

	if minField != nil && maxField != nil {
		minVal, _ := minField.GetValue("")
		maxVal, _ := maxField.GetValue("")
		min, _ := minVal.Int()
		max, _ := maxVal.Int()

		quantitySchema := builder.Field().
			AddValidator("min", min).
			AddValidator("max", max).
			Build()
		s.AddField("quantity", quantitySchema)
	}

	// 2. 根据订单类型调整价格要求
	typeField, _ := obj.GetField("OrderType")
	if typeField != nil {
		typeVal, _ := typeField.GetValue("")
		orderType := typeVal.String()

		if orderType == "same_day" {
			// 当日达最低100元
			priceSchema := builder.Field().
				AddValidator("min", 100).
				Build()
			s.AddField("price", priceSchema)
		}
	}
}

func main() {
	fmt.Println("=== Comprehensive Example: Dynamic Order Form ===")

	// 创建自定义registry注册enum和between validator
	registry := validators.NewRegistry()

	// 注册enum validator
	registry.Register("enum", func(ctx *schema.Context, params []string) error {
		allowedValues := make(map[string]bool)
		for _, p := range params {
			allowedValues[p] = true
		}
		val := ctx.Value()
		if !allowedValues[val.String()] {
			return schema.NewValidationError(ctx.Path(), "enum", map[string]any{
				"allowed": params,
			})
		}
		return nil
	})

	// 注册between validator
	registry.Register("between", func(ctx *schema.Context, before, after any) error {
		min, max := cast.To[int](before), cast.To[int](after)
		valAcc := ctx.Value()
		val, _ := valAcc.Int()
		if val < int64(min) || val > int64(max) {
			return schema.NewValidationError(ctx.Path(), "between", map[string]any{
				"min": min, "max": max,
			})
		}
		return nil
	})

	// 创建validator
	typ := reflect.TypeOf(OrderForm{})
	objSchema, _ := builder.Parse(typ, builder.WithRegistry(registry))
	v := validator.NewFromSchema(objSchema)

	// 测试1: 有效标准订单
	fmt.Println("Test 1: Valid Standard Order")
	fmt.Println("----------------------------")
	validOrder := OrderForm{
		OrderType:   "standard",
		Price:       50,
		MinQuantity: 1,
		MaxQuantity: 10,
		Quantity:    5,
	}
	err := v.Validate(validOrder)
	fmt.Printf("Result: %v\n", err == nil)
	if err != nil {
		if res, ok := err.(*schema.ValidationResult); ok {
			for _, e := range res.Errors() {
				fmt.Printf("  - %s: %s %v\n", e.Path, e.Name, e.Params)
			}
		} else {
			fmt.Printf("  - %v\n", err)
		}
	}
	fmt.Println()

	// 测试2: 当日达但价格太低
	fmt.Println("Test 2: Same-Day Order with Low Price")
	fmt.Println("--------------------------------------")
	invalidOrder := OrderForm{
		OrderType:   "same_day",
		Price:       50, // 太低，当日达最低100
		MinQuantity: 1,
		MaxQuantity: 10,
		Quantity:    5,
	}
	err = v.Validate(invalidOrder)
	fmt.Printf("Result: %v\n", err == nil)
	if err != nil {
		if res, ok := err.(*schema.ValidationResult); ok {
			for _, e := range res.Errors() {
				fmt.Printf("  - %s: %s %v\n", e.Path, e.Name, e.Params)
			}
		} else {
			fmt.Printf("  - %v\n", err)
		}
	}
	fmt.Println()

	// 测试3: 数量超出范围
	fmt.Println("Test 3: Quantity Exceeds Max")
	fmt.Println("----------------------------")
	invalidOrder2 := OrderForm{
		OrderType:   "standard",
		Price:       50,
		MinQuantity: 1,
		MaxQuantity: 10,
		Quantity:    20, // 超过maxQuantity
	}
	err = v.Validate(invalidOrder2)
	fmt.Printf("Result: %v\n", err == nil)
	if err != nil {
		if res, ok := err.(*schema.ValidationResult); ok {
			for _, e := range res.Errors() {
				fmt.Printf("  - %s: %s %v\n", e.Path, e.Name, e.Params)
			}
		} else {
			fmt.Printf("  - %v\n", err)
		}
	}
	fmt.Println()

	// 测试4: 有效当日达订单
	fmt.Println("Test 4: Valid Same-Day Order")
	fmt.Println("----------------------------")
	validSameDayOrder := OrderForm{
		OrderType:   "same_day",
		Price:       150, // 满足当日达最低价格
		MinQuantity: 1,
		MaxQuantity: 5,
		Quantity:    3,
	}
	err = v.Validate(validSameDayOrder)
	fmt.Printf("Result: %v\n", err == nil)
	fmt.Println()

	// 测试5: 无效订单类型
	fmt.Println("Test 5: Invalid Order Type")
	fmt.Println("--------------------------")
	invalidOrderType := OrderForm{
		OrderType:   "overnight", // 不在允许的枚举值中
		Price:       50,
		MinQuantity: 1,
		MaxQuantity: 10,
		Quantity:    5,
	}
	err = v.Validate(invalidOrderType)
	fmt.Printf("Result: %v\n", err == nil)
	if err != nil {
		if res, ok := err.(*schema.ValidationResult); ok {
			for _, e := range res.Errors() {
				fmt.Printf("  - %s: %s %v\n", e.Path, e.Name, e.Params)
			}
		} else {
			fmt.Printf("  - %v\n", err)
		}
	}
	fmt.Println()

	// 测试6: 数量低于最小值
	fmt.Println("Test 6: Quantity Below Min")
	fmt.Println("--------------------------")
	belowMinQuantity := OrderForm{
		OrderType:   "express",
		Price:       80,
		MinQuantity: 5,
		MaxQuantity: 20,
		Quantity:    3, // 低于minQuantity
	}
	err = v.Validate(belowMinQuantity)
	fmt.Printf("Result: %v\n", err == nil)
	if err != nil {
		if res, ok := err.(*schema.ValidationResult); ok {
			for _, e := range res.Errors() {
				fmt.Printf("  - %s: %s %v\n", e.Path, e.Name, e.Params)
			}
		} else {
			fmt.Printf("  - %v\n", err)
		}
	}

	fmt.Println("\n=== Comprehensive Example Completed ===")
}
