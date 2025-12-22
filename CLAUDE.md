# Schema Validator - 完整功能文档

## 目录

- [项目概述](#项目概述)
- [快速开始](#快速开始)
- [核心概念](#核心概念)
- [基础功能](#基础功能)
- [高级功能](#高级功能)
- [内置验证器](#内置验证器)
- [自定义验证器](#自定义验证器)
- [API参考](#api参考)

---

## 项目概述

Schema Validator 是一个灵活、强大的 Go 数据验证库，提供以下核心特性：

✅ **统一Schema格式** - 支持Field、Array、Object三种类型
✅ **多数据源支持** - primitives、arrays、maps、structs
✅ **双API模式** - struct tags 和 code-based builders
✅ **跨字段验证** - tag语法和代码两种方式
✅ **嵌入结构体** - 支持访问嵌入结构体的私有字段
✅ **动态Schema** - 运行时根据数据修改验证规则
✅ **多参数验证器** - 支持复杂自定义验证器
✅ **Schema可视化** - ToString()方法输出JSON格式
✅ **清晰错误报告** - 字段路径 + 错误码

---

## 快速开始

### 安装

```bash
go get github.com/weilence/schema-validator
```

### 基本使用

#### 方式1: Struct Tags

```go
import "github.com/weilence/schema-validator/validator"

type User struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min_length=8"`
    Age      int    `json:"age" validate:"min=18,max=120"`
}

v, _ := validator.NewFromStruct(User{})
user := User{Email: "test@example.com", Password: "password123", Age: 25}
result, _ := v.Validate(user)

if !result.IsValid() {
    for _, err := range result.Errors() {
        fmt.Printf("%s: %s\n", err.FieldPath, err.ErrorCode)
    }
}
```

#### 方式2: Code-based Builder

```go
import (
    validator "github.com/weilence/schema-validator"
    "github.com/weilence/schema-validator/schema"
    "github.com/weilence/schema-validator/validation"
)

userSchema := schema.Object().
    Field("email", schema.Field().
        AddValidator(validation.Required()).
        AddValidator(validation.Email()).
        Build()).
    Field("age", schema.Field().
        AddValidator(validation.Min(18)).
        Build()).
    Build()

v := validator.New(userSchema)
result, _ := v.Validate(map[string]interface{}{
    "email": "test@example.com",
    "age":   25,
})
```

---

## 核心概念

### 1. Schema 系统

Schema定义了数据的验证规则，分为三种类型：

- **FieldSchema**: 验证单个字段（primitive类型）
- **ArraySchema**: 验证数组/切片
- **ObjectSchema**: 验证对象（struct/map）

Schema可以嵌套组合，构建复杂的验证规则。

### 2. Data Accessor 抽象层

统一的数据访问接口，屏蔽不同数据源的差异：

```go
type Accessor interface {
    Kind() DataKind                           // 数据类型
    IsNil() bool                              // 是否为nil
    AsField() (FieldAccessor, error)          // 转为字段访问器
    AsObject() (ObjectAccessor, error)        // 转为对象访问器
    AsArray() (ArrayAccessor, error)          // 转为数组访问器
}
```

### 3. Validation Context

验证上下文携带验证状态，支持跨字段验证：

```go
type Context struct {
    Root   data.Accessor      // 根对象
    Path   string              // 当前路径 (e.g., "user.email")
    Parent data.ObjectAccessor // 父对象
}
```

### 4. Error 格式

只返回错误码和字段路径，不包含消息（支持国际化）：

```go
type ValidationError struct {
    FieldPath string                   // "user.email", "items[0].name"
    ErrorCode string                   // "required", "min", "email"
    Params    map[string]interface{}   // 附加参数
}
```

---

## 基础功能

### Tag语法

在struct字段上使用`validate`tag定义验证规则：

```go
type Product struct {
    Name  string  `json:"name" validate:"required,min_length=3,max_length=100"`
    Price float64 `json:"price" validate:"min=0"`
    Tags  []string `json:"tags" validate:"min_items=1,max_items=5"`
}
```

**Tag规则**：
- 多个验证器用逗号分隔: `required,min=5,max=100`
- 参数用等号传递: `min_length=8`
- 多参数用冒号分隔: `between=10:100`

### 嵌入结构体

支持嵌入结构体，包括访问私有字段：

```go
type Address struct {
    Street  string `validate:"required"`
    private string // 当嵌入时可访问
}

type Person struct {
    Name string `validate:"required"`
    Address     // 嵌入结构体
}

v, _ := validator.NewFromStruct(Person{})
```

**重要**：
- ✅ 嵌入结构体的私有字段可访问
- ❌ 非嵌入结构体的私有字段不可访问

### 跨字段验证

#### Tag方式

```go
type PasswordForm struct {
    Password string `json:"password" validate:"required,min_length=8"`
    Confirm  string `json:"confirm" validate:"required,eqfield=password"`
}
```

#### Code方式

```go
passwordValidator := validation.ObjectValidatorFunc(
    func(ctx *validation.Context, obj data.ObjectAccessor) error {
        pwd, _ := obj.GetField("password")
        confirm, _ := obj.GetField("confirm")

        if pwd.AsField().String() != confirm.AsField().String() {
            return errors.NewValidationError(
                ctx.Path + ".confirm",
                "password_mismatch",
                nil,
            )
        }
        return nil
    },
)

schema := schema.Object().
    Field("password", schema.Field().AddValidator(validation.Required()).Build()).
    Field("confirm", schema.Field().AddValidator(validation.Required()).Build()).
    CrossField(passwordValidator).
    Build()
```

### 数组验证

```go
type TodoList struct {
    Items []string `json:"items" validate:"min_items=1,max_items=10"`
}

// 或使用代码方式
itemSchema := schema.Field().AddValidator(validation.MinLength(1)).Build()
arraySchema := schema.Array(itemSchema).MinItems(1).MaxItems(10).Build()
```

---

## 高级功能

### 动态Schema修改 (SchemaModifier)

**功能说明**：允许struct在验证前根据运行时数据动态修改验证规则。

#### 接口定义

```go
type SchemaModifier interface {
    ModifySchema(ctx *validation.Context, accessor data.ObjectAccessor, schema *ObjectSchema)
}
```

#### 示例1: 条件必填

根据标志字段动态设置其他字段是否必填：

```go
type DynamicForm struct {
    Type     string `json:"type" validate:"required"`
    Value    string `json:"value"`
    Required bool   `json:"required"`
}

func (f DynamicForm) ModifySchema(ctx *validation.Context, accessor data.ObjectAccessor, s *schema.ObjectSchema) {
    // 读取required字段的值
    requiredField, _ := accessor.GetField("required")
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

// 使用
v, _ := validator.NewFromStruct(DynamicForm{})

// required=true时，value必填
form1 := DynamicForm{Type: "text", Value: "", Required: true}
result, _ := v.Validate(form1) // 失败

// required=false时，value可选
form2 := DynamicForm{Type: "text", Value: "", Required: false}
result, _ := v.Validate(form2) // 通过
```

#### 示例2: 访问嵌套对象

根据嵌套对象的值添加验证规则：

```go
type User struct {
    Name    string  `json:"name" validate:"required"`
    Address Address `json:"address"`
}

type Address struct {
    Country string `json:"country"`
    ZipCode string `json:"zipCode"`
}

func (u User) ModifySchema(ctx *validation.Context, accessor data.ObjectAccessor, s *schema.ObjectSchema) {
    // 访问嵌套的address.country
    addressField, _ := accessor.GetField("address")
    addressObj, _ := addressField.AsObject()
    countryField, _ := addressObj.GetField("country")
    countryAcc, _ := countryField.AsField()
    country := countryAcc.String()

    // 根据国家添加不同的邮编验证
    if country == "US" {
        addressSchema := schema.Object().
            Field("country", schema.Field().Build()).
            Field("zipCode", schema.Field().
                AddValidator(validation.MinLength(5)).
                AddValidator(validation.MaxLength(5)).
                Build()).
            Build()
        s.AddField("address", addressSchema)
    }
}
```

#### 示例3: 访问Parent和Root

```go
func (obj MyObject) ModifySchema(ctx *validation.Context, accessor data.ObjectAccessor, s *schema.ObjectSchema) {
    // 访问父对象
    if ctx.Parent != nil {
        parentField, _ := ctx.Parent.GetField("someField")
        // 基于父对象值修改schema
    }

    // 访问根对象
    if ctx.Root != nil {
        rootObj, _ := ctx.Root.AsObject()
        rootField, _ := rootObj.GetField("someRootField")
        // 基于根对象值修改schema
    }
}
```

#### 访问数组

```go
func (obj MyObject) ModifySchema(ctx *validation.Context, accessor data.ObjectAccessor, s *schema.ObjectSchema) {
    // 访问数组字段
    itemsField, _ := accessor.GetField("items")
    itemsArray, _ := itemsField.AsArray()

    // 遍历数组元素
    itemsArray.Iterate(func(idx int, elem data.Accessor) error {
        elemObj, _ := elem.AsObject()
        // 访问数组元素的字段
        field, _ := elemObj.GetField("fieldName")
        // 基于数组内容修改schema
        return nil
    })
}
```

**适用场景**：
- ✅ 条件必填字段
- ✅ 基于类型的不同验证规则
- ✅ 复杂的业务逻辑验证
- ✅ 动态表单验证
- ✅ 多态数据验证

### 多参数验证器

**功能说明**：ValidatorFactory支持接收多个参数，使用冒号(:)分隔。

#### 参数格式

```go
validate:"between=10:100"        // ["10", "100"]
validate:"enum=red:green:blue"   // ["red", "green", "blue"]
validate:"range=0:100:5"         // ["0", "100", "5"] (min, max, step)
```

⚠️ **重要**：使用冒号(:)而不是逗号(,)分隔参数，因为逗号用于分隔validator。

#### 示例1: between验证器 (2参数)

```go
registry := tags.NewRegistry()

registry.RegisterField("between", func(params []string) (validation.FieldValidator, error) {
    if len(params) < 2 {
        return nil, nil
    }

    min := parseInt(params[0])
    max := parseInt(params[1])

    return validation.FieldValidatorFunc(func(ctx *validation.Context, value data.FieldAccessor) error {
        val, _ := value.Int()
        if val < int64(min) || val > int64(max) {
            return errors.NewValidationError(ctx.Path, "between", map[string]interface{}{
                "min": min, "max": max, "actual": val,
            })
        }
        return nil
    }), nil
})

// 使用
type Product struct {
    Price int `json:"price" validate:"between=10:100"`
}

typ := reflect.TypeOf(Product{})
objSchema, _ := tags.ParseStructTagsWithRegistry(typ, registry)
v := validator.New(objSchema)
```

#### 示例2: enum验证器 (多参数)

```go
registry.RegisterField("enum", func(params []string) (validation.FieldValidator, error) {
    allowedValues := make(map[string]bool)
    for _, p := range params {
        allowedValues[p] = true
    }

    return validation.FieldValidatorFunc(func(ctx *validation.Context, value data.FieldAccessor) error {
        val := value.String()
        if !allowedValues[val] {
            return errors.NewValidationError(ctx.Path, "enum", map[string]interface{}{
                "allowed": params,
                "actual":  val,
            })
        }
        return nil
    }), nil
})

// 使用
type Settings struct {
    Theme string `json:"theme" validate:"enum=light:dark:auto"`
}
```

#### 示例3: range验证器 (3参数)

```go
// validate:"range=0:100:5" - 0到100，步长为5
registry.RegisterField("range", func(params []string) (validation.FieldValidator, error) {
    if len(params) < 3 {
        return nil, nil
    }

    min := parseInt(params[0])
    max := parseInt(params[1])
    step := parseInt(params[2])

    return validation.FieldValidatorFunc(func(ctx *validation.Context, value data.FieldAccessor) error {
        val, _ := value.Int()

        if val < int64(min) || val > int64(max) {
            return errors.NewValidationError(ctx.Path, "out_of_range", nil)
        }

        if (int(val)-min)%step != 0 {
            return errors.NewValidationError(ctx.Path, "invalid_step", nil)
        }

        return nil
    }), nil
})
```

### Schema的JSON表示 (ToString)

**功能说明**：所有Schema都提供`ToString()`方法，可以将schema结构输出为格式化的JSON字符串，方便调试和文档生成。

#### 特性

- ✅ 输出合法的JSON格式
- ✅ 自动格式化（缩进）
- ✅ 支持嵌套结构（Object嵌套、Array元素等）
- ✅ 包含所有验证器信息

#### 示例1: Field Schema

```go
fieldSchema := schema.Field().
    AddValidator(validation.Required()).
    AddValidator(validation.MinLength(5)).
    AddValidator(validation.Email()).
    Build()

fmt.Println(fieldSchema.ToString())
```

输出：
```json
{
  "type": "field",
  "optional": false,
  "validators": [
    {
      "name": "required"
    },
    {
      "name": "min_length",
      "value": 5
    },
    {
      "name": "email"
    }
  ]
}
```

#### 示例2: Object Schema with Nested Fields

```go
userSchema := schema.Object().
    Field("name", schema.Field().AddValidator(validation.Required()).Build()).
    Field("email", schema.Field().AddValidator(validation.Email()).Build()).
    Field("age", schema.Field().AddValidator(validation.Min(18)).Build()).
    Build()

fmt.Println(userSchema.ToString())
```

输出：
```json
{
  "type": "object",
  "strict": false,
  "fields": {
    "name": {
      "type": "field",
      "optional": false,
      "validators": [
        {
          "name": "required"
        }
      ]
    },
    "email": {
      "type": "field",
      "optional": false,
      "validators": [
        {
          "name": "email"
        }
      ]
    },
    "age": {
      "type": "field",
      "optional": false,
      "validators": [
        {
          "name": "min",
          "value": 18
        }
      ]
    }
  }
}
```

#### 示例3: Complex Nested Schema

```go
addressSchema := schema.Object().
    Field("street", schema.Field().AddValidator(validation.Required()).Build()).
    Field("city", schema.Field().AddValidator(validation.Required()).Build()).
    Build()

phoneSchema := schema.Object().
    Field("type", schema.Field().AddValidator(validation.Required()).Build()).
    Field("number", schema.Field().AddValidator(validation.Required()).Build()).
    Build()

userSchema := schema.Object().
    Field("name", schema.Field().AddValidator(validation.Required()).Build()).
    Field("phones", schema.Array(phoneSchema).MinItems(1).Build()).
    Field("address", addressSchema).
    Build()

fmt.Println(userSchema.ToString())
```

输出：
```json
{
  "type": "object",
  "strict": false,
  "fields": {
    "name": {
      "type": "field",
      "optional": false,
      "validators": [
        {
          "name": "required"
        }
      ]
    },
    "phones": {
      "type": "array",
      "minItems": 1,
      "element": {
        "type": "object",
        "strict": false,
        "fields": {
          "type": {
            "type": "field",
            "optional": false,
            "validators": [
              {
                "name": "required"
              }
            ]
          },
          "number": {
            "type": "field",
            "optional": false,
            "validators": [
              {
                "name": "required"
              }
            ]
          }
        }
      },
      "validators": [
        {
          "name": "min_items",
          "value": 1
        }
      ]
    },
    "address": {
      "type": "object",
      "strict": false,
      "fields": {
        "street": {
          "type": "field",
          "optional": false,
          "validators": [
            {
              "name": "required"
            }
          ]
        },
        "city": {
          "type": "field",
          "optional": false,
          "validators": [
            {
              "name": "required"
            }
          ]
        }
      }
    }
  }
}
```

#### JSON结构说明

**FieldSchema输出格式**：
```json
{
  "type": "field",
  "optional": true/false,
  "validators": [
    {
      "name": "validator_name",
      "value": "parameter_value (如果有)"
    }
  ]
}
```

**ArraySchema输出格式**：
```json
{
  "type": "array",
  "element": { /* 嵌套的element schema */ },
  "minItems": number,
  "maxItems": number,
  "validators": [ /* array validators */ ]
}
```

**ObjectSchema输出格式**：
```json
{
  "type": "object",
  "strict": true/false,
  "fields": {
    "field_name": { /* 嵌套的field schema */ }
  },
  "validators": [ /* object validators */ ]
}
```

#### 适用场景

- ✅ **调试** - 查看schema结构是否符合预期
- ✅ **文档生成** - 自动生成API文档
- ✅ **Schema导出** - 将schema保存为配置文件
- ✅ **测试** - 验证schema构建是否正确
- ✅ **日志记录** - 记录使用的验证规则

---

## 内置验证器

### 字段验证器

| Tag | 说明 | 示例 | 参数 |
|-----|------|------|------|
| `required` | 必填 | `validate:"required"` | 无 |
| `min` | 最小值 | `validate:"min=18"` | 数值 |
| `max` | 最大值 | `validate:"max=120"` | 数值 |
| `min_length` | 最小长度 | `validate:"min_length=8"` | 数值 |
| `max_length` | 最大长度 | `validate:"max_length=100"` | 数值 |
| `email` | 邮箱格式 | `validate:"email"` | 无 |
| `url` | URL格式 | `validate:"url"` | 无 |

### 跨字段验证器

| Tag | 说明 | 示例 | 参数 |
|-----|------|------|------|
| `eqfield` | 等于另一字段 | `validate:"eqfield=password"` | 字段名 |
| `nefield` | 不等于另一字段 | `validate:"nefield=oldPassword"` | 字段名 |
| `gtfield` | 大于另一字段 | `validate:"gtfield=startDate"` | 字段名 |
| `ltfield` | 小于另一字段 | `validate:"ltfield=endDate"` | 字段名 |

### 数组验证器

| Tag | 说明 | 示例 | 参数 |
|-----|------|------|------|
| `min_items` | 最小元素数 | `validate:"min_items=1"` | 数值 |
| `max_items` | 最大元素数 | `validate:"max_items=10"` | 数值 |

---

## 自定义验证器

### 字段验证器

```go
// 自定义validator函数
customValidator := validation.FieldValidatorFunc(
    func(ctx *validation.Context, value data.FieldAccessor) error {
        val := value.String()
        if !isValid(val) {
            return errors.NewValidationError(ctx.Path, "custom_error", nil)
        }
        return nil
    },
)

// 使用
schema := schema.Field().AddValidator(customValidator).Build()
```

### 对象验证器 (跨字段)

```go
crossFieldValidator := validation.ObjectValidatorFunc(
    func(ctx *validation.Context, obj data.ObjectAccessor) error {
        field1, _ := obj.GetField("field1")
        field2, _ := obj.GetField("field2")

        // 跨字段验证逻辑
        if !isValid(field1, field2) {
            return errors.NewValidationError(ctx.Path, "cross_field_error", nil)
        }
        return nil
    },
)

// 使用
schema := schema.Object().
    Field("field1", schema.Field().Build()).
    Field("field2", schema.Field().Build()).
    CrossField(crossFieldValidator).
    Build()
```

### 注册到Registry

```go
registry := tags.NewRegistry()

// 单参数
registry.RegisterField("custom", func(params []string) (validation.FieldValidator, error) {
    param := params[0]
    return validation.FieldValidatorFunc(func(ctx *validation.Context, value data.FieldAccessor) error {
        // 验证逻辑
        return nil
    }), nil
})

// 多参数
registry.RegisterField("between", func(params []string) (validation.FieldValidator, error) {
    min, max := parseInt(params[0]), parseInt(params[1])
    // 返回validator
})

// 使用自定义registry
typ := reflect.TypeOf(MyStruct{})
objSchema, _ := tags.ParseStructTagsWithRegistry(typ, registry)
v := validator.New(objSchema)
```

---

## API参考

### Validator

```go
// 从code-based schema创建
func New(s schema.Schema) *Validator

// 从struct tags创建
func NewFromStruct(prototype interface{}) (*Validator, error)

// 验证数据
func (v *Validator) Validate(data interface{}) (*errors.ValidationResult, error)

// 检查是否有效
func (v *Validator) IsValid(data interface{}) bool

// 验证，panic on error
func (v *Validator) MustValidate(data interface{}) *errors.ValidationResult
```

### Schema Builders

```go
// Field schema
schema.Field().
    Required().
    AddValidator(validator).
    Build()

// Array schema
schema.Array(elementSchema).
    MinItems(1).
    MaxItems(10).
    Build()

// Object schema
schema.Object().
    Field("name", fieldSchema).
    CrossField(objectValidator).
    Strict().
    Build()
```

### ValidationResult

```go
// 是否有效
func (r *ValidationResult) IsValid() bool

// 获取所有错误
func (r *ValidationResult) Errors() []*ValidationError

// 获取第一个错误
func (r *ValidationResult) FirstError() *ValidationError

// 按字段分组错误
func (r *ValidationResult) ErrorsByField() map[string][]*ValidationError

// 检查特定字段是否有错误
func (r *ValidationResult) HasFieldError(fieldPath string) bool
```

### ValidationError

```go
type ValidationError struct {
    FieldPath string                   // 字段路径
    ErrorCode string                   // 错误码
    Params    map[string]interface{}   // 参数
}

func (e *ValidationError) Error() string
```

---

## 完整示例

### 综合使用示例

```go
package main

import (
    "fmt"
    "reflect"

    validator "github.com/weilence/schema-validator"
    "github.com/weilence/schema-validator/data"
    "github.com/weilence/schema-validator/errors"
    "github.com/weilence/schema-validator/schema"
    "github.com/weilence/schema-validator/tags"
    "github.com/weilence/schema-validator/validation"
)

// 订单表单，综合使用多个高级功能
type OrderForm struct {
    OrderType   string `json:"orderType" validate:"enum=standard:express:same_day"`
    Price       int    `json:"price" validate:"between=1:10000"`
    MinQuantity int    `json:"minQuantity"`
    MaxQuantity int    `json:"maxQuantity"`
    Quantity    int    `json:"quantity"`
}

// 实现SchemaModifier接口
func (f OrderForm) ModifySchema(ctx *validation.Context, accessor data.ObjectAccessor, s *schema.ObjectSchema) {
    // 1. 根据min/max动态设置quantity范围
    minField, _ := accessor.GetField("minQuantity")
    maxField, _ := accessor.GetField("maxQuantity")

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
    typeField, _ := accessor.GetField("orderType")
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
    // 创建自定义registry注册enum和between validator
    registry := tags.NewRegistry()

    // 注册enum validator
    registry.RegisterField("enum", func(params []string) (validation.FieldValidator, error) {
        allowedValues := make(map[string]bool)
        for _, p := range params {
            allowedValues[p] = true
        }
        return validation.FieldValidatorFunc(func(ctx *validation.Context, value data.FieldAccessor) error {
            if !allowedValues[value.String()] {
                return errors.NewValidationError(ctx.Path, "enum", map[string]interface{}{
                    "allowed": params,
                })
            }
            return nil
        }), nil
    })

    // 注册between validator
    registry.RegisterField("between", func(params []string) (validation.FieldValidator, error) {
        min, max := parseInt(params[0]), parseInt(params[1])
        return validation.FieldValidatorFunc(func(ctx *validation.Context, value data.FieldAccessor) error {
            val, _ := value.Int()
            if val < int64(min) || val > int64(max) {
                return errors.NewValidationError(ctx.Path, "between", map[string]interface{}{
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

    // 测试1: 有效订单
    validOrder := OrderForm{
        OrderType:   "standard",
        Price:       50,
        MinQuantity: 1,
        MaxQuantity: 10,
        Quantity:    5,
    }
    result, _ := v.Validate(validOrder)
    fmt.Printf("Valid order: %v\n", result.IsValid())

    // 测试2: 当日达但价格太低
    invalidOrder := OrderForm{
        OrderType:   "same_day",
        Price:       50, // 太低，当日达最低100
        MinQuantity: 1,
        MaxQuantity: 10,
        Quantity:    5,
    }
    result, _ = v.Validate(invalidOrder)
    fmt.Printf("Invalid order (low price): %v\n", result.IsValid())
    for _, err := range result.Errors() {
        fmt.Printf("  - %s: %s %v\n", err.FieldPath, err.ErrorCode, err.Params)
    }

    // 测试3: 数量超出范围
    invalidOrder2 := OrderForm{
        OrderType:   "standard",
        Price:       50,
        MinQuantity: 1,
        MaxQuantity: 10,
        Quantity:    20, // 超过maxQuantity
    }
    result, _ = v.Validate(invalidOrder2)
    fmt.Printf("Invalid order (quantity): %v\n", result.IsValid())
    for _, err := range result.Errors() {
        fmt.Printf("  - %s: %s %v\n", err.FieldPath, err.ErrorCode, err.Params)
    }
}

func parseInt(s string) int {
    var result int
    for _, ch := range s {
        if ch >= '0' && ch <= '9' {
            result = result*10 + int(ch-'0')
        }
    }
    return result
}
```

---

## 设计原则

### 1. 错误码优先

只返回错误码，不返回错误消息，便于国际化：

```go
// ✅ 好的设计
err := ValidationError{
    FieldPath: "user.email",
    ErrorCode: "invalid_email",
    Params:    map[string]interface{}{"pattern": emailRegex},
}

// ❌ 避免
err := ValidationError{
    Message: "Email format is invalid", // 硬编码消息
}
```

### 2. 数据抽象

统一的Accessor接口抽象不同数据源：

- Struct → StructAccessor
- Map → MapAccessor
- Slice → SliceAccessor
- Primitive → FieldAccessor

### 3. Schema组合

Schema可以嵌套组合，构建复杂验证规则：

```go
userSchema := Object().
    Field("profile", Object().
        Field("name", Field().Required().Build()).
        Field("tags", Array(Field()).MinItems(1).Build()).
        Build()).
    Build()
```

### 4. 性能优化

- 反射缓存：`StructInfo`缓存结构体元数据
- 字段索引：使用`FieldByIndex`快速访问嵌入字段
- 延迟计算：只在需要时转换数据类型

---

## 最佳实践

### 1. 选择合适的API

- **简单场景**：使用struct tags
- **复杂逻辑**：使用code-based builders
- **运行时规则**：使用SchemaModifier

### 2. 错误处理

```go
result, err := v.Validate(data)
if err != nil {
    // 系统错误
    log.Error(err)
    return
}

if !result.IsValid() {
    // 验证错误
    for _, verr := range result.Errors() {
        // 根据ErrorCode显示本地化消息
        msg := i18n.Get(verr.ErrorCode, verr.Params)
        fmt.Println(msg)
    }
}
```

### 3. 自定义Validator

优先使用已有的组合方式，必要时才创建自定义validator：

```go
// ✅ 优先：组合已有validator
schema.Field().
    AddValidator(validation.MinLength(8)).
    AddValidator(validation.MaxLength(20)).
    Build()

// ✅ 必要时：自定义validator
customValidator := validation.FieldValidatorFunc(...)
```

### 4. SchemaModifier使用建议

- 只在需要运行时决策时使用
- 保持逻辑简单清晰
- 避免在ModifySchema中进行复杂计算
- 考虑性能影响

---

## 版本历史

### v1.2.0 (Latest)
- ✅ 新增Schema.ToString()方法，支持将schema输出为JSON格式
- ✅ 支持嵌套schema的JSON序列化
- ✅ 包含所有validator信息的完整JSON表示

### v1.1.0
- ✅ 新增SchemaModifier接口，支持动态schema修改
- ✅ ValidatorFactory支持多参数（使用冒号分隔）
- ✅ 改进tag解析器，智能处理参数中的逗号
- ✅ 添加SetOptional方法到FieldSchemaBuilder

### v1.0.0
- ✅ 基础验证功能
- ✅ Struct tags支持
- ✅ Code-based builders
- ✅ 跨字段验证
- ✅ 嵌入结构体支持（包括私有字段）
- ✅ 内置验证器

---

## 贡献指南

### 添加新功能

所有新功能必须更新本文档的相应章节：

1. 在[高级功能](#高级功能)章节添加功能说明
2. 在[API参考](#api参考)更新API文档
3. 在[完整示例](#完整示例)添加使用示例
4. 更新[版本历史](#版本历史)

### 文档结构

- **快速开始**: 5分钟上手
- **核心概念**: 理解架构
- **基础功能**: 常用功能
- **高级功能**: 复杂场景
- **API参考**: 完整API
- **完整示例**: 实战代码

---

## 许可证

MIT License
