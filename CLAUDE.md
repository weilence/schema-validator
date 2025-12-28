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
- [完整示例](#完整示例)
- [设计原则](#设计原则)
- [最佳实践](#最佳实践)

---

## 项目概述

Schema Validator 是一个灵活、强大的 Go 数据验证库，提供以下核心特性：

- ✅ **统一 Schema 格式** - 支持 Field、Array、Object 三种类型
- ✅ **多数据源支持** - primitives、arrays、maps、structs
- ✅ **双 API 模式** - struct tags 和 code-based builders
- ✅ **跨字段验证** - tag 语法和代码两种方式
- ✅ **嵌入结构体** - 支持访问嵌入结构体的私有字段
- ✅ **动态 Schema** - 运行时根据数据修改验证规则 (SchemaModifier)
- ✅ **自定义验证器** - 支持多参数自定义验证器
- ✅ **清晰错误报告** - 字段路径 + 错误码（支持国际化）

---

## 快速开始

### 安装

```bash
go get github.com/weilence/schema-validator
```

### 基本使用

#### 方式1: Struct Tags

```go
import validator "github.com/weilence/schema-validator"

type User struct {
    Email    string `json:"email" validate:"required|email"`
    Password string `json:"password" validate:"required|min=8"`
    Age      int    `json:"age" validate:"min=18|max=120"`
}

v, _ := validator.New(User{})
user := User{Email: "test@example.com", Password: "password123", Age: 25}
err := v.Validate(user)

if err != nil {
    fmt.Println(err)
}
```

#### 方式2: Code-based Builder

```go
import (
    validator "github.com/weilence/schema-validator"
    "github.com/weilence/schema-validator/builder"
)

userSchema := builder.Object().
    Field("email", builder.Field().
        AddValidator("required").
        AddValidator("email").
        Build()).
    Field("age", builder.Field().
        AddValidator("min", 18).
        AddValidator("max", 120).
        Build()).
    Build()

v := validator.NewFromSchema(userSchema)
err := v.Validate(map[string]interface{}{
    "email": "test@example.com",
    "age":   25,
})
```

---

## 核心概念

### 1. Schema 系统

Schema 定义了数据的验证规则，分为三种类型：

- **FieldSchema**: 验证单个字段（primitive 类型）
- **ArraySchema**: 验证数组/切片
- **ObjectSchema**: 验证对象（struct/map）

Schema 可以嵌套组合，构建复杂的验证规则。

**示例**：

```go
// Field Schema
fieldSchema := builder.Field().
    AddValidator("required").
    AddValidator("email").
    Build()

// Array Schema
arraySchema := builder.Array(fieldSchema).
    AddValidator("min", 1).
    AddValidator("max", 10).
    Build()

// Object Schema
objectSchema := builder.Object().
    Field("name", builder.Field().AddValidator("required").Build()).
    Field("tags", arraySchema).
    Build()
```

### 2. Data Accessor 抽象层

统一的数据访问接口，屏蔽不同数据源的差异：

```go
type Accessor interface {
    GetValue(path string) (*Value, error)  // 获取路径上的值
    Raw() any                                // 获取原始数据
}
```

**实现类型**：
- `Value` - 原始值（int, string, bool等）
- `ArrayAccessor` - 数组/切片访问器
- `MapAccessor` - Map访问器
- `StructAccessor` - 结构体访问器

**优势**：
- 统一的接口处理不同数据类型
- 支持路径访问（如 "user.address.city"）
- 自动处理类型转换

### 3. Validation Context

验证上下文携带验证状态，支持跨字段验证：

```go
type Context struct {
    schema   Schema         // 当前 schema
    accessor Accessor       // 当前数据访问器
    parent   *Context       // 父上下文
    path     string         // 当前路径 (e.g., "user.email", "items[0]")
}
```

**主要方法**：
- `Path()` - 获取当前字段路径
- `Value()` - 获取当前字段值
- `GetValue(path)` - 获取相对路径的值
- `Parent()` - 获取父上下文（用于跨字段验证）
- `Schema()` - 获取当前 schema
- `Accessor()` - 获取当前数据访问器

**Context 的作用**：
- 提供当前验证位置的完整上下文信息
- 支持访问父对象和根对象进行跨字段验证
- 构建准确的错误路径

### 4. Error 格式

只返回错误码和字段路径，不包含消息（支持国际化）：

```go
type ValidationError struct {
    Path   string   // 字段路径: "user.email", "items[0].name"
    Name   string   // 错误码: "required", "min", "email"
    Params []any    // 附加参数: ["min=8", "actual=5"]
}
```

**错误处理**：

```go
err := v.Validate(data)
if err != nil {
    // 方式1: 单个验证错误
    if validationErr, ok := err.(*schema.ValidationError); ok {
        fmt.Printf("%s: %s\n", validationErr.Path, validationErr.Name)
    }
    
    // 方式2: 多个验证错误
    if validationResult, ok := err.(*schema.ValidationResult); ok {
        for _, e := range validationResult.Errors() {
            fmt.Printf("%s: %s %v\n", e.Path, e.Name, e.Params)
        }
    }
}
```

**国际化支持**：

```go
// 根据错误码和参数显示本地化消息
err := v.Validate(data)
if validationErr, ok := err.(*schema.ValidationError); ok {
    msg := i18n.Get(validationErr.Name, validationErr.Params)
    fmt.Println(msg)
}
```

---

## 基础功能

### Tag语法

在struct字段上使用`validate`tag定义验证规则：

```go
type Product struct {
    Name  string  `json:"name" validate:"required|min=3,max=100"`
    Price float64 `json:"price" validate:"min=0"`
    Tags  []string `json:"tags" validate:"min_items=1,max_items=5"`
}
```

**Tag规则**：
- 多个验证器用逗号分隔: `required,min=5,max=100`
- 参数用等号传递: `min=8`
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
    Password string `json:"password" validate:"required|min=8"`
    Confirm  string `json:"confirm" validate:"required|eqfield=password"`
}
```

#### Code方式

```go
passwordValidator := validation.ObjectValidatorFunc(
    func(ctx *validation.Context) error {
        obj, _ := ctx.AsObject()
        pwd, _ := obj.GetField("password")
        confirm, _ := obj.GetField("confirm")

        pwdField, _ := pwd.AsField()
        confirmField, _ := confirm.AsField()

        if pwdField.String() != confirmField.String() {
            return schema.ErrCheckFailed
        }
        return nil
    },
)

schema := tags.Object().
    Field("password", tags.Field().AddValidator(validation.Required()).Build()).
    Field("confirm", tags.Field().AddValidator(validation.Required()).Build()).
    CrossField(passwordValidator).
    Build()
```

### 数组验证

```go
type TodoList struct {
    Items []string `json:"items" validate:"dive|required|min=3"`
}

// 或使用代码方式
itemSchema := builder.Field().AddValidator("required").AddValidator("min", 3).Build()
arraySchema := builder.Array(itemSchema).AddValidator("min", 1).AddValidator("max", 10).Build()
```

**说明**：
- 使用 `dive` 关键字进入数组，验证每个元素
- 数组本身可以使用 `min`/`max` 验证器限制长度

---

## 高级功能

### 动态Schema修改 (SchemaModifier)

**功能说明**：允许struct在验证前根据运行时数据动态修改验证规则。

#### 接口定义

```go
type SchemaModifier interface {
    ModifySchema(ctx *schema.Context)
}
```

**核心概念**：

- `ctx` 提供完整的验证上下文信息：
  - 当前验证路径（通过 `ctx.Path()`）
  - 父上下文（通过 `ctx.Parent()`）
  - 当前 Schema（通过 `ctx.Schema()`）
  - 当前数据访问器（通过 `ctx.Accessor()`）

#### 示例代码

```go
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

    // 根据 required 标志动态修改 value 字段的验证
    if f.Required {
        s.AddField("value", builder.Field().
            AddValidator("required").
            Build())
    } else {
        s.RemoveField("value")
    }
}
```

#### 高级用法

- **访问当前对象数据** - 结构体可以直接访问自身字段（如 `f.Required`）
- **访问父上下文** - 使用 `ctx.Parent()` 访问父对象上下文
- **修改 Schema** - 通过 `AddField()` 或 `RemoveField()` 动态修改验证规则
- **获取路径信息** - 使用 `ctx.Path()` 获取当前验证路径

**完整示例**: 查看 [examples/dynamic_schema/main.go](examples/dynamic_schema/main.go) 了解更多用法

**适用场景**：
- ✅ 条件必填字段
- ✅ 基于类型的不同验证规则
- ✅ 复杂的业务逻辑验证
- ✅ 动态表单验证
- ✅ 多态数据验证

### 自定义验证器注册

**功能说明**：通过 validators.Registry 注册自定义验证器，支持多参数。

#### 参数格式

```go
validate:"between=10:100"        // params: ["10", "100"]
validate:"enum=red:green:blue"   // params: ["red", "green", "blue"]
validate:"range=0:100:5"         // params: ["0", "100", "5"]
```

⚠️ **重要**：多个参数使用冒号 `:` 分隔。

#### 示例: between验证器

```go
registry := validators.NewRegistry()

registry.Register("between", func(ctx *schema.Context, params []string) error {
    min := parseInt(params[0])
    max := parseInt(params[1])

    fieldVal, _ := ctx.Value()
    val, _ := fieldVal.Int()
    
    if val < int64(min) || val > int64(max) {
        return schema.ErrCheckFailed
    }
    return nil
})

// 使用自定义 registry
type Product struct {
    Price int `json:"price" validate:"between=10:100"`
}

v, _ := validator.New(Product{}, builder.WithRegistry(registry))
```

**完整示例**: 查看 [examples/custom_validators/main.go](examples/custom_validators/main.go) 了解 between、enum、range 等多参数验证器的实现

---

## 内置验证器

### 字段验证器

| Tag | 说明 | 示例 | 适用类型 |
|-----|------|------|----------|
| `required` | 必填 | `validate:"required"` | 所有类型 |
| `min` | 最小值/长度 | `validate:"min=18"` | 数值/字符串/数组 |
| `max` | 最大值/长度 | `validate:"max=120"` | 数值/字符串/数组 |
| `email` | 邮箱格式 | `validate:"email"` | 字符串 |
| `ip` | IP 地址 | `validate:"ip"` | 字符串 |
| `port` | 端口号 (1-65535) | `validate:"port"` | 数值 |
| `domain` | 域名 | `validate:"domain"` | 字符串 |
| `fqdn` | 完全限定域名 | `validate:"fqdn"` | 字符串 |
| `url` | URL 格式 | `validate:"url"` | 字符串 |
| `pattern` | 正则表达式 | `validate:"pattern=^[A-Z]+$"` | 字符串 |
| `oneof` | 枚举值 | `validate:"oneof=red:green:blue"` | 字符串 |

### 跨字段验证器

| Tag | 说明 | 示例 | 参数说明 |
|-----|------|------|----------|
| `eqfield` | 等于另一字段 | `validate:"eqfield=Password"` | 字段名（大小写敏感）|
| `nefield` | 不等于另一字段 | `validate:"nefield=OldPassword"` | 字段名 |
| `gtfield` | 大于另一字段 | `validate:"gtfield=StartDate"` | 字段名 |
| `ltfield` | 小于另一字段 | `validate:"ltfield=EndDate"` | 字段名 |
| `required_if` | 条件必填 | `validate:"required_if=Type:premium"` | 字段名:期望值 |

### 数组验证

使用 `dive` 关键字验证数组元素：

```go
type List struct {
    Items []string `validate:"dive|required|min=3"`
}
```

数组本身可以使用 `min`/`max` 验证器限制长度。

---

## 自定义验证器

### 注册自定义验证器

```go
registry := validators.NewRegistry()

// 注册一个自定义验证器
registry.Register("between", func(ctx *schema.Context, params []string) error {
    min := parseInt(params[0])
    max := parseInt(params[1])

    fieldVal, _ := ctx.Value()
    val, _ := fieldVal.Int()
    
    if val < int64(min) || val > int64(max) {
        return schema.ErrCheckFailed
    }
    return nil
})

// 使用自定义 registry
type Product struct {
    Price int `json:"price" validate:"between=10:100"`
}

v, _ := validator.New(Product{}, builder.WithRegistry(registry))
```

### 验证器类型

验证器函数签名支持多种参数组合：

```go
// 无参数
func(ctx *schema.Context) error

// 单个参数
func(ctx *schema.Context, param string) error
func(ctx *schema.Context, param int) error

// 多个参数
func(ctx *schema.Context, params []string) error
func(ctx *schema.Context, min int, max int) error
```

Registry 会自动识别函数签名并进行参数类型转换。

### 完整示例

查看 [examples/custom_validators/main.go](examples/custom_validators/main.go) 了解：
- between 验证器（范围验证）
- enum 验证器（枚举值验证）
- range 验证器（带步长的范围验证）

---

## API参考

### Validator

```go
// 从 struct tags 创建
func New(prototype any, opts ...builder.ParseOption) (*Validator, error)

// 从 code-based schema 创建
func NewFromSchema(s schema.Schema) *Validator

// 验证数据
func (v *Validator) Validate(value any) error
```

**返回值**：
- `nil` - 验证成功
- `*schema.ValidationError` - 单个验证错误
- `*schema.ValidationResult` - 多个验证错误

### Schema Builders

```go
// Field schema
builder.Field().
    AddValidator("required").
    AddValidator("email").
    Build()

// Array schema
builder.Array(elementSchema).
    AddValidator("min", 1).
    AddValidator("max", 10).
    Build()

// Object schema
builder.Object().
    Field("name", fieldSchema).
    Field("email", emailSchema).
    Build()
```

### Builder Options

```go
// 使用自定义 registry
builder.WithRegistry(registry)

// 设置规则分隔符（默认 '|'）
builder.WithRuleSplitter('|')

// 设置参数分隔符（默认 ':'）
builder.WithParamsSeparator(':')

// 设置 dive 关键字（默认 "dive"）
builder.WithDiveTag("dive")
```

### ValidationResult

```go
type ValidationResult struct {
    // ...
}

// 获取所有错误
func (r *ValidationResult) Errors() []*ValidationError

// 是否有效（无错误）
func (r *ValidationResult) IsValid() bool

// 获取第一个错误
func (r *ValidationResult) FirstError() *ValidationError

// 按字段分组错误
func (r *ValidationResult) ErrorsByField() map[string][]*ValidationError

// 检查特定字段是否有错误
func (r *ValidationResult) HasFieldError(fieldPath string) bool

// 实现 error 接口
func (r *ValidationResult) Error() string
```

### ValidationError

```go
type ValidationError struct {
    Path   string   // 字段路径: "user.email", "items[0].name"
    Name   string   // 错误码: "required", "min", "email"
    Params []any    // 附加参数: ["min=8", "actual=5"]
    Err    error    // 内部错误（可选）
}

// 实现 error 接口
func (e *ValidationError) Error() string

// 创建验证错误
func NewValidationError(path, name string, params map[string]any) *ValidationError
```

**错误处理示例**：

```go
err := v.Validate(data)
if err != nil {
    // 方式1: 单个验证错误
    if validationErr, ok := err.(*schema.ValidationError); ok {
        fmt.Printf("%s: %s\n", validationErr.Path, validationErr.Name)
    }
    
    // 方式2: 多个验证错误
    if validationResult, ok := err.(*schema.ValidationResult); ok {
        for _, e := range validationResult.Errors() {
            fmt.Printf("%s: %s %v\n", e.Path, e.Name, e.Params)
        }
    }
}
```

---

## 完整示例

完整的可运行示例代码已移至 `examples/` 目录：

### 基础示例

**[examples/basic/main.go](examples/basic/main.go)** - 涵盖基础功能：
- Tag-based validation（标签验证）
- Code-based validation（代码方式验证）
- Embedded struct with private fields（嵌入结构体及私有字段）
- Array validation（数组验证）
- Cross-field validation（跨字段验证）
- Error handling patterns（错误处理模式）

### 高级功能示例

**[examples/dynamic_schema/main.go](examples/dynamic_schema/main.go)** - 动态Schema（SchemaModifier）：
- 条件必填字段
- 基于类型的验证（如国家特定的邮编验证）

**[examples/custom_validators/main.go](examples/custom_validators/main.go)** - 多参数自定义验证器：
- between validator（范围验证）
- enum validator（枚举验证）
- range validator with step（步长范围验证）

**[examples/comprehensive/main.go](examples/comprehensive/main.go)** - 综合示例：
- 结合 SchemaModifier + 多参数验证器
- 动态订单表单验证
- 复杂业务逻辑验证

### 运行示例

```bash
# 基础示例
go run examples/basic/main.go

# 动态Schema示例
go run examples/dynamic_schema/main.go

# 自定义验证器示例
go run examples/custom_validators/main.go

# 综合示例
go run examples/comprehensive/main.go
```

---

## 设计原则

### 1. 错误码优先

只返回错误码，不返回错误消息，便于国际化：

```go
// ✅ 好的设计
err := schema.ErrCheckFailed

// ❌ 避免
err := ValidationError{
    Message: "Email format is invalid", // 硬编码消息
}
```

**国际化支持**：
```go
err := v.Validate(data)
if validationErr, ok := err.(*schema.ValidationError); ok {
    msg := i18n.Get(validationErr.Name, validationErr.Params)
    fmt.Println(msg)
}
```

### 2. 数据抽象

统一的 Accessor 接口抽象不同数据源：

- Struct → StructAccessor
- Map → MapAccessor
- Array/Slice → ArrayAccessor
- Primitive → Value

**优势**：
- 统一的数据访问接口
- 支持路径访问（如 "user.address.city"）
- 自动处理类型转换

### 3. Schema 组合

Schema 可以嵌套组合，构建复杂验证规则：

```go
userSchema := builder.Object().
    Field("profile", builder.Object().
        Field("name", builder.Field().AddValidator("required").Build()).
        Field("tags", builder.Array(builder.Field().AddValidator("required").Build()).
            AddValidator("min", 1).Build()).
        Build()).
    Build()
```

### 4. 性能优化

- **反射缓存**：缓存结构体元数据
- **字段索引**：使用 `FieldByIndex` 快速访问嵌入字段
- **延迟计算**：只在需要时转换数据类型

### 5. 可扩展性

- **自定义验证器**：通过 Registry 注册
- **动态 Schema**：通过 SchemaModifier 接口
- **多种数据源**：通过 Accessor 接口

---

## 最佳实践

### 1. 选择合适的 API

- **简单场景**：使用 struct tags
- **复杂逻辑**：使用 code-based builders
- **运行时规则**：使用 SchemaModifier

```go
// 简单场景 - 使用 tags
type User struct {
    Email string `validate:"required|email"`
    Age   int    `validate:"min=18|max=120"`
}

// 复杂逻辑 - 使用 builder
schema := builder.Object().
    Field("email", builder.Field().AddValidator("required").AddValidator("email").Build()).
    Build()

// 运行时 - 使用 SchemaModifier
func (u User) ModifySchema(ctx *schema.Context) {
    // 动态修改验证规则
}
```

### 2. 错误处理

```go
err := v.Validate(data)
if err != nil {
    // 检查验证错误
    if validationResult, ok := err.(*schema.ValidationResult); ok {
        for _, verr := range validationResult.Errors() {
            // 根据 Name 显示本地化消息
            msg := i18n.Get(verr.Name, verr.Params)
            fmt.Println(msg)
        }
    } else if validationErr, ok := err.(*schema.ValidationError); ok {
        msg := i18n.Get(validationErr.Name, validationErr.Params)
        fmt.Println(msg)
    } else {
        // 系统错误
        log.Error(err)
    }
}
```

### 3. 自定义 Validator

优先使用已有的组合方式，必要时才创建自定义 validator：

```go
// ✅ 优先：组合已有 validator
builder.Field().
    AddValidator("min", 8).
    AddValidator("max", 20).
    Build()

// ✅ 必要时：自定义 validator
registry.Register("custom", func(ctx *schema.Context) error {
    // 自定义逻辑
    return nil
})
```

### 4. SchemaModifier 使用建议

- 只在需要运行时决策时使用
- 保持逻辑简单清晰
- 避免在 ModifySchema 中进行复杂计算
- 考虑性能影响

### 5. Tag 语法注意事项

- 使用 `|` 分隔验证器：`"required|email|min=8"`
- 多参数使用 `:` 分隔：`"between=10:100"`
- 数组元素验证使用 `dive`：`"dive|required"`
- 字段名大小写敏感：`"eqfield=Password"`

---

## 版本历史

### v1.0.0 (Current)
- ✅ 基础验证功能
- ✅ Struct tags 支持
- ✅ Code-based builders
- ✅ 跨字段验证
- ✅ 嵌入结构体支持（包括私有字段）
- ✅ SchemaModifier 接口，支持动态 schema 修改
- ✅ 自定义验证器支持多参数
- ✅ 内置验证器（required, min, max, email, ip, port, domain, url, pattern, oneof 等）
- ✅ 跨字段验证器（eqfield, nefield, gtfield, ltfield, required_if）
- ✅ 数组验证（dive 关键字）

---

## 贡献指南

欢迎贡献！请遵循以下指南：

### 添加新功能

所有新功能必须：

1. 在相应章节添加功能说明
2. 更新 API 参考文档
3. 添加使用示例到 `examples/` 目录
4. 更新版本历史

### 文档结构

本文档按以下结构组织：

- **项目概述** - 核心特性一览
- **快速开始** - 5分钟上手指南  
- **核心概念** - 理解架构和设计
- **基础功能** - 常用功能说明
- **高级功能** - 复杂场景处理
- **内置验证器** - 完整的验证器列表
- **自定义验证器** - 扩展验证功能
- **API参考** - 完整API文档
- **完整示例** - 实战代码
- **设计原则** - 架构设计理念
- **最佳实践** - 使用建议

### 代码贡献

1. Fork 本仓库
2. 创建功能分支: `git checkout -b feature/my-feature`
3. 提交更改: `git commit -am 'Add some feature'`
4. 推送到分支: `git push origin feature/my-feature`
5. 提交 Pull Request

### 报告问题

使用 GitHub Issues 报告问题时，请包含：

- 问题描述
- 复现步骤
- 期望行为
- 实际行为
- Go 版本和操作系统信息

---

## 许可证

MIT License

---

## 相关链接

- **GitHub**: https://github.com/weilence/schema-validator
- **Issues**: https://github.com/weilence/schema-validator/issues
- **示例代码**: [examples/](examples/)
- **主文档**: [README.md](README.md)
