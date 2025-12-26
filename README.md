# Schema Validator

一个灵活的Go数据校验库，支持schema定义和多种数据格式校验。

## 特性

- ✅ **统一的Schema格式** 支持字段、数组、对象三种格式
- ✅ **多种数据类型** 支持primitives、arrays、maps、structs
- ✅ **双API设计** 同时支持struct tags和代码构建器
- ✅ **跨字段校验** tag语法和代码两种方式
- ✅ **嵌入结构体支持** 支持访问嵌入结构体的私有字段
- ✅ **清晰的错误报告** 返回字段路径 + 错误码
- ✅ **可扩展** 易于添加自定义校验器

## 安装

```bash
go get github.com/weilence/schema-validator
```

## 快速开始

### 方式1: 使用Struct Tags

```go
package main

import (
    "fmt"
    "github.com/weilence/schema-validator/validator"
)

type User struct {
    Email    string `json:"email" validate:"required|email"`
    Password string `json:"password" validate:"required|min_length=8"`
    Confirm  string `json:"confirm" validate:"required|eqfield=Password"`
    Age      int    `json:"age" validate:"min=18|max=120"`
}

func main() {
    v, _ := validator.NewFromStruct(User{})

    user := User{
        Email:    "test@example.com",
        Password: "password123",
        Confirm:  "password123",
        Age:      25,
    }

    result, _ := v.Validate(user)

    if result.IsValid() {
        fmt.Println("用户数据有效！")
    } else {
        for _, err := range result.Errors() {
            fmt.Printf("%s: %s\n", err.FieldPath, err.ErrorCode)
        }
    }
}
```

### 方式2: 使用代码构建器

```go
package main

import (
    "fmt"
    "github.com/weilence/schema-validator/schema"
    "github.com/weilence/schema-validator/validation"
    "github.com/weilence/schema-validator/validator"
)

func main() {
    userSchema := schema.Object().
        Field("email", schema.Field().
            AddValidator(validation.Required()).
            AddValidator(validation.Email())).
        Field("password", schema.Field().
            AddValidator(validation.Required()).
            AddValidator(validation.MinLength(8))).
        Field("age", schema.Field().
            AddValidator(validation.Min(18)).
            AddValidator(validation.Max(120))).
        Build()

    v := validator.New(userSchema)

    data := map[string]interface{}{
        "email":    "test@example.com",
        "password": "password123",
        "age":      25,
    }

    result, _ := v.Validate(data)

    if result.IsValid() {
        fmt.Println("数据有效！")
    } else {
        for _, err := range result.Errors() {
            fmt.Printf("%s: %s\n", err.FieldPath, err.ErrorCode)
        }
    }
}
```

## 高级用法

### 嵌入结构体

```go
type Address struct {
    Street  string `json:"street" validate:"required"`
    City    string `json:"city" validate:"required"`
}

type Person struct {
    Name string `json:"name" validate:"required"`
    Address // 嵌入结构体
}

v, _ := validator.NewFromStruct(Person{})
person := Person{
    Name: "John Doe",
    Address: Address{
        Street: "123 Main St",
        City:   "New York",
    },
}

result, _ := v.Validate(person)
```

### 数组验证

```go
type TodoList struct {
    Items []string `json:"items" validate:"min_items=1|max_items=10"`
}

v, _ := validator.NewFromStruct(TodoList{})
list := TodoList{
    Items: []string{"task1", "task2"},
}

result, _ := v.Validate(list)
```

### 自定义跨字段验证

```go
// 使用代码构建器添加自定义验证
passwordMatchValidator := validation.ObjectValidatorFunc(
    func(ctx *validation.Context, obj data.ObjectAccessor) error {
        password, _ := obj.GetField("password")
        confirm, _ := obj.GetField("confirmPassword")

        pwdField, _ := password.AsField()
        confField, _ := confirm.AsField()

        if pwdField.String() != confField.String() {
            return schema.NewValidationError(
                ctx.Path + ".confirmPassword",
                "password_mismatch",
                nil,
            )
        }
        return nil
    },
)

schema := schema.Object().
    Field("password", schema.Field().AddValidator(validation.Required())).
    Field("confirmPassword", schema.Field().AddValidator(validation.Required())).
    CrossField(passwordMatchValidator).
    Build()
```

## 错误处理

```go
result, err := v.Validate(data)
if err != nil {
    // 系统错误
    panic(err)
}

if !result.IsValid() {
    // 获取所有错误
    for _, validationErr := range result.Errors() {
        fmt.Printf("字段: %s, 错误: %s, 参数: %v\n",
            validationErr.FieldPath,
            validationErr.ErrorCode,
            validationErr.Params)
    }

    // 按字段分组错误
    errorsByField := result.ErrorsByField()
    for field, errs := range errorsByField {
        fmt.Printf("字段 %s 有 %d 个错误\n", field, len(errs))
    }
}
```

## 设计原则

1. **数据抽象层** - 统一的Accessor接口抽象不同数据类型
2. **Schema系统** - 可组合的Schema (Field/Array/Object)
3. **验证引擎** - 基于Context的验证，支持跨字段访问
4. **错误码优先** - 只返回错误码，消息由外部处理（支持i18n）
5. **反射缓存** - 缓存结构体元数据以提升性能

## 架构

```
validator/
├── data/         # 数据抽象层
├── schema/       # Schema定义
├── validation/   # 验证引擎和验证器
├── tags/         # Tag解析
├── errors/       # 错误类型
└── reflect/      # 反射工具
```

## License

MIT
