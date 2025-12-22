package main

import (
	"fmt"
	"reflect"

	validator "github.com/weilence/schema-validator"
	"github.com/weilence/schema-validator/errors"
	"github.com/weilence/schema-validator/schema"
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

// OrderForm 订单表单，综合使用多个高级功能
type OrderForm struct {
	OrderType   string `json:"orderType" validate:"enum=standard:express:same_day"`
	Price       int    `json:"price" validate:"between=1:10000"`
	MinQuantity int    `json:"minQuantity"`
	MaxQuantity int    `json:"maxQuantity"`
	Quantity    int    `json:"quantity"`
}

// 实现SchemaModifier接口
func (f OrderForm) ModifySchema(ctx *validation.Context) {
	s, ok := ctx.ObjectSchema().(*schema.ObjectSchema)
	if !ok || s == nil {
		return
	}

	obj, _ := ctx.AsObject()

	// 1. 根据min/max动态设置quantity范围
	minField, _ := obj.GetField("minQuantity")
	maxField, _ := obj.GetField("maxQuantity")

	if minField != nil && maxField != nil {
		minAcc, _ := minField.AsField()
		maxAcc, _ := maxField.AsField()
		min, _ := minAcc.Int()
		max, _ := maxAcc.Int()

		quantitySchema := schema.Field().
			AddValidator(validation.Min(int(min))).
			AddValidator(validation.Max(int(max))).
			Build()
		s.AddField("quantity", quantitySchema)
	}

	// 2. 根据订单类型调整价格要求
	typeField, _ := obj.GetField("orderType")
	if typeField != nil {
		typeAcc, _ := typeField.AsField()
		orderType := typeAcc.String()

		if orderType == "same_day" {
			// 当日达最低100元
			priceSchema := schema.Field().
				AddValidator(validation.Min(100)).
				Build()
			s.AddField("price", priceSchema)
		}
	}
}

func main() {
	fmt.Println("=== Comprehensive Example: Dynamic Order Form ===")

	// 创建自定义registry注册enum和between validator
	registry := tags.NewRegistry()

	// 注册enum validator
	registry.RegisterField("enum", func(params []string) (validation.FieldValidator, error) {
		allowedValues := make(map[string]bool)
		for _, p := range params {
			allowedValues[p] = true
		}
		return validation.FieldValidatorFunc(func(ctx *validation.Context) error {
			field, _ := ctx.AsField()
			if !allowedValues[field.String()] {
				return errors.NewValidationError(ctx.Path(), "enum", map[string]interface{}{
					"allowed": params,
				})
			}
			return nil
		}), nil
	})

	// 注册between validator
	registry.RegisterField("between", func(params []string) (validation.FieldValidator, error) {
		min, max := parseInt(params[0]), parseInt(params[1])
		return validation.FieldValidatorFunc(func(ctx *validation.Context) error {
			field, _ := ctx.AsField()
			val, _ := field.Int()
			if val < int64(min) || val > int64(max) {
				return errors.NewValidationError(ctx.Path(), "between", map[string]interface{}{
					"min": min, "max": max,
				})
			}
			return nil
		}), nil
	})

	// 创建validator
	typ := reflect.TypeOf(OrderForm{})
	objSchema, _ := tags.ParseStructTagsWithRegistry(typ, registry)
	v := validator.New(objSchema)

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
	result, _ := v.Validate(validOrder)
	fmt.Printf("Result: %v\n", result.IsValid())
	if !result.IsValid() {
		for _, err := range result.Errors() {
			fmt.Printf("  - %s: %s %v\n", err.FieldPath, err.ErrorCode, err.Params)
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
	result, _ = v.Validate(invalidOrder)
	fmt.Printf("Result: %v\n", result.IsValid())
	for _, err := range result.Errors() {
		fmt.Printf("  - %s: %s %v\n", err.FieldPath, err.ErrorCode, err.Params)
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
	result, _ = v.Validate(invalidOrder2)
	fmt.Printf("Result: %v\n", result.IsValid())
	for _, err := range result.Errors() {
		fmt.Printf("  - %s: %s %v\n", err.FieldPath, err.ErrorCode, err.Params)
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
	result, _ = v.Validate(validSameDayOrder)
	fmt.Printf("Result: %v\n", result.IsValid())
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
	result, _ = v.Validate(invalidOrderType)
	fmt.Printf("Result: %v\n", result.IsValid())
	for _, err := range result.Errors() {
		fmt.Printf("  - %s: %s %v\n", err.FieldPath, err.ErrorCode, err.Params)
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
	result, _ = v.Validate(belowMinQuantity)
	fmt.Printf("Result: %v\n", result.IsValid())
	for _, err := range result.Errors() {
		fmt.Printf("  - %s: %s %v\n", err.FieldPath, err.ErrorCode, err.Params)
	}

	fmt.Println("\n=== Comprehensive Example Completed ===")
}
