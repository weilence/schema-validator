# Schema Validator

一个灵活强大的 Go 数据验证库，支持 schema 定义和多种数据格式验证。

## 特性

- ✅ **统一的 Schema 格式** - 支持 Field、Array、Object 三种类型
- ✅ **多数据源支持** - 支持 primitives、arrays、maps、structs
- ✅ **双 API 设计** - 同时支持 struct tags 和代码构建器
- ✅ **跨字段验证** - 支持 tag 语法和代码两种方式
- ✅ **嵌入结构体** - 支持访问嵌入结构体的私有字段
- ✅ **动态 Schema** - 运行时根据数据修改验证规则 (SchemaModifier)
- ✅ **自定义验证器** - 易于注册和扩展验证器
- ✅ **清晰的错误报告** - 返回字段路径 + 错误码 (支持国际化)

## 安装

```bash
go get github.com/weilence/schema-validator
```

## 快速开始

### 方式1: 使用 Struct Tags

```go
package main

import (
    "fmt"
    validator "github.com/weilence/schema-validator"
)

type User struct {
    Email    string `json:"email" validate:"required|email"`
    Password string `json:"password" validate:"required|min=8"`
    Confirm  string `json:"confirm" validate:"required|eqfield=Password"`
    Age      int    `json:"age" validate:"min=18|max=120"`
}

func main() {
    v, _ := validator.New(User{})

    user := User{
        Email:    "test@example.com",
        Password: "password123",
        Confirm:  "password123",
        Age:      25,
    }

    err := v.Validate(user)

    if err == nil {
        fmt.Println("用户数据有效！")
    } else {
        fmt.Println(err)
    }
}
```

### 方式2: 使用代码构建器

```go
package main

import (
    "fmt"
    validator "github.com/weilence/schema-validator"
    "github.com/weilence/schema-validator/builder"
)

func main() {
    userSchema := builder.Object().
        Field("email", builder.Field().
            AddValidator("required").
            AddValidator("email").
            Build()).
        Field("password", builder.Field().
            AddValidator("required").
            AddValidator("min", 8).
            Build()).
        Field("age", builder.Field().
            AddValidator("min", 18).
            AddValidator("max", 120).
            Build()).
        Build()

    v := validator.NewFromSchema(userSchema)

    data := map[string]interface{}{
        "email":    "test@example.com",
        "password": "password123",
        "age":      25,
    }

    err := v.Validate(data)

    if err == nil {
        fmt.Println("数据有效！")
    } else {
        fmt.Println(err)
    }
}
```

## 核心概念

### Schema 系统

Schema 定义了数据的验证规则，分为三种类型：

- **FieldSchema**: 验证单个字段（primitive 类型）
- **ArraySchema**: 验证数组/切片
- **ObjectSchema**: 验证对象（struct/map）

Schema 可以嵌套组合，构建复杂的验证规则。

### Data Accessor 抽象层

统一的数据访问接口，屏蔽不同数据源的差异：

- Struct → StructAccessor
- Map → MapAccessor
- Slice/Array → ArrayAccessor
- Primitive → Value (FieldAccessor)

### 验证上下文 (Context)

验证上下文携带验证状态，支持跨字段验证：

```go
type Context struct {
    schema   Schema         // 当前 schema
    accessor Accessor       // 当前数据访问器
    parent   *Context       // 父上下文
    path     string         // 当前路径 (e.g., "user.email")
}
```

### 错误格式

只返回错误码和字段路径，不包含消息（支持国际化）：

```go
type ValidationError struct {
    Path   string   // "user.email", "items[0].name"
    Name   string   // "required", "min", "email"
    Params []any    // 附加参数
}
```

## 主要功能

### Tag 语法

在 struct 字段上使用 `validate` tag 定义验证规则：

```go
type Product struct {
    Name  string  `json:"name" validate:"required|min=3|max=100"`
    Price float64 `json:"price" validate:"min=0"`
    Tags  []string `json:"tags" validate:"dive|required"`
}
```

**Tag 规则**：
- 多个验证器用竖线分隔: `required|min=5|max=100`
- 参数用等号传递: `min=8`
- 多参数用冒号分隔: `between=10:100`
- 数组元素验证使用 `dive` 关键字

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

v, _ := validator.New(Person{})
```

**重要**：
- ✅ 嵌入结构体的私有字段可访问
- ❌ 非嵌入结构体的私有字段不可访问

### 跨字段验证

#### Tag 方式

```go
type PasswordForm struct {
    Password string `json:"password" validate:"required|min=8"`
    Confirm  string `json:"confirm" validate:"required|eqfield=Password"`
}
```

#### 代码方式

可以通过 Context 访问父对象和其他字段进行跨字段验证。

### 数组验证

```go
type TodoList struct {
    Items []string `json:"items" validate:"dive|required"`
}

// 或使用代码方式
itemSchema := builder.Field().AddValidator("required").Build()
arraySchema := builder.Array(itemSchema).
    AddValidator("min", 1).
    AddValidator("max", 10).
    Build()
```

### 动态 Schema (SchemaModifier)

允许 struct 在验证前根据运行时数据动态修改验证规则：

```go
type DynamicForm struct {
    Type     string `json:"type" validate:"required"`
    Value    string `json:"value"`
    Required bool   `json:"required"`
}

func (f DynamicForm) ModifySchema(ctx *schema.Context) {
    s := ctx.Schema().(*schema.ObjectSchema)
    
    if f.Required {
        s.AddField("value", builder.Field().
            AddValidator("required").
            Build())
    }
}
```

**适用场景**：
- 条件必填字段
- 基于类型的不同验证规则
- 复杂的业务逻辑验证

## 内置验证器

### 字段验证器

| Tag | 说明 | 示例 |
|-----|------|------|
| `required` | 必填 | `validate:"required"` |
| `min` | 最小值/长度 | `validate:"min=8"` |
| `max` | 最大值/长度 | `validate:"max=100"` |
| `email` | 邮箱格式 | `validate:"email"` |
| `ip` | IP 地址 | `validate:"ip"` |
| `port` | 端口号 (1-65535) | `validate:"port"` |
| `domain` | 域名 | `validate:"domain"` |
| `fqdn` | 完全限定域名 | `validate:"fqdn"` |
| `url` | URL 格式 | `validate:"url"` |
| `pattern` | 正则表达式 | `validate:"pattern=^[A-Z]+$"` |
| `oneof` | 枚举值 | `validate:"oneof=red:green:blue"` |

### 跨字段验证器

| Tag | 说明 | 示例 |
|-----|------|------|
| `eqfield` | 等于另一字段 | `validate:"eqfield=Password"` |
| `nefield` | 不等于另一字段 | `validate:"nefield=OldPassword"` |
| `gtfield` | 大于另一字段 | `validate:"gtfield=StartDate"` |
| `ltfield` | 小于另一字段 | `validate:"ltfield=EndDate"` |
| `required_if` | 条件必填 | `validate:"required_if=Type:premium"` |

### 数组验证器

使用 `dive` 关键字验证数组元素：

```go
type List struct {
    Items []string `validate:"dive|required|min=3"`
}
```

## 自定义验证器

### 注册自定义验证器

```go
registry := validators.NewRegistry()

// 注册自定义验证器
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

## 错误处理

```go
err := v.Validate(data)
if err != nil {
    // 检查是否是验证错误
    if validationResult, ok := err.(*schema.ValidationResult); ok {
        // 获取所有错误
        for _, validationErr := range validationResult.Errors() {
            fmt.Printf("字段: %s, 错误: %s, 参数: %v\n",
                validationErr.Path,
                validationErr.Name,
                validationErr.Params)
        }

        // 按字段分组错误
        errorsByField := validationResult.ErrorsByField()
        for field, errs := range errorsByField {
            fmt.Printf("字段 %s 有 %d 个错误\n", field, len(errs))
        }
    } else if validationErr, ok := err.(*schema.ValidationError); ok {
        // 单个验证错误
        fmt.Printf("%s: %s\n", validationErr.Path, validationErr.Name)
    }
}
```

## 示例

完整的示例代码位于 `examples/` 目录：

- **[examples/basic/](examples/basic/)** - 基础用法示例
- **[examples/dynamic_schema/](examples/dynamic_schema/)** - 动态 Schema 示例
- **[examples/custom_validators/](examples/custom_validators/)** - 自定义验证器示例
- **[examples/comprehensive/](examples/comprehensive/)** - 综合示例

运行示例：

```bash
go run examples/basic/main.go
go run examples/dynamic_schema/main.go
go run examples/custom_validators/main.go
go run examples/comprehensive/main.go
```

## 设计原则

1. **数据抽象层** - 统一的 Accessor 接口抽象不同数据类型
2. **Schema 系统** - 可组合的 Schema (Field/Array/Object)
3. **验证引擎** - 基于 Context 的验证，支持跨字段访问
4. **错误码优先** - 只返回错误码，消息由外部处理（支持 i18n）
5. **反射缓存** - 缓存结构体元数据以提升性能

## 架构

```
schema-validator/
├── validator.go        # 验证器主入口
├── data/              # 数据抽象层
│   ├── accessor.go    # 统一访问接口
│   ├── value.go       # 原始值访问器
│   ├── array_accessor.go
│   ├── map_accessor.go
│   └── struct_accessor.go
├── schema/            # Schema 定义
│   ├── schema.go      # Schema 接口
│   ├── field.go       # 字段 Schema
│   ├── array.go       # 数组 Schema
│   ├── object.go      # 对象 Schema
│   ├── context.go     # 验证上下文
│   ├── error.go       # 错误定义
│   └── validator.go   # 验证器接口
├── builder/           # Schema 构建器
│   ├── builder.go     # 流式 API
│   └── parser.go      # Tag 解析器
└── validators/        # 内置验证器
    ├── registry.go    # 验证器注册表
    ├── field.go       # 跨字段验证器
    ├── format.go      # 格式验证器
    ├── network.go     # 网络验证器
    └── other.go       # 其他验证器
```

## API 参考

完整的 API 文档请参考 [CLAUDE.md](CLAUDE.md)。

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

