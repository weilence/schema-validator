# 示例代码

本目录包含 schema-validator 库的各种功能示例。

## 运行示例

每个示例都在其自己的子目录中。运行示例：

```bash
# 运行基础验证示例
go run examples/basic/main.go

# 运行动态 Schema 示例
go run examples/dynamic_schema/main.go

# 运行自定义验证器示例
go run examples/custom_validators/main.go

# 运行综合示例
go run examples/comprehensive/main.go
```

## 可用示例

### 1. 基础验证 (`basic/`)

涵盖核心验证功能的综合示例：

- **示例 1**: 使用 struct tags 的基于标签的验证
- **示例 2**: 使用 maps 的基于代码的验证
- **示例 3**: 带私有字段的嵌入结构体
- **示例 4**: 数组验证
- **示例 5**: 跨字段验证
- **示例 6**: 错误处理模式

**运行：**
```bash
go run examples/basic/main.go
```

**输出示例：**
```
=== Example 1: Tag-based Validation ===
Valid user result: true
Invalid user result: false
  - email: email
  - password: min
  - confirm: eqfield
  - age: min
```

### 2. 动态 Schema (`dynamic_schema/`)

演示 SchemaModifier 接口的使用：

- **示例 1**: 条件必填字段 - 根据运行时标志动态添加必填验证
- **示例 2**: 基于类型的验证 - 根据国家代码应用不同的邮编验证规则

**运行：**
```bash
go run examples/dynamic_schema/main.go
```

**特性：**
- 运行时根据数据修改验证规则
- 访问当前对象的字段值
- 动态添加或删除字段验证

### 3. 自定义验证器 (`custom_validators/`)

演示如何注册和使用多参数自定义验证器：

- **示例 1**: `between` 验证器 - 范围验证（2个参数）
- **示例 2**: `enum` 验证器 - 枚举值验证（多个参数）
- **示例 3**: `range` 验证器 - 带步长的范围验证（3个参数）

**运行：**
```bash
go run examples/custom_validators/main.go
```

**要点：**
- 使用 `validators.Registry` 注册自定义验证器
- 多参数使用冒号 (`:`) 分隔
- 参数在 tag 中的格式: `validate:"between=10:100"`

### 4. 综合示例 (`comprehensive/`)

结合多个高级特性的实际应用场景：

- SchemaModifier + 自定义多参数验证器
- 动态订单表单验证
- 复杂业务逻辑验证

**运行：**
```bash
go run examples/comprehensive/main.go
```

**演示内容：**
- 根据订单类型动态修改验证规则
- 使用 `between` 验证器验证数量范围
- 组合使用多种验证器

## 示例说明

### Tag 语法

在示例中，你会看到以下 tag 格式：

```go
type User struct {
    Email string `validate:"required|email"`
    Age   int    `validate:"min=18|max=120"`
}
```

**要点：**
- 多个验证器用竖线 `|` 分隔
- 参数用等号 `=` 传递
- 多参数用冒号 `:` 分隔（如 `between=10:100`）

### 数组验证

使用 `dive` 关键字验证数组元素：

```go
type List struct {
    Items []string `validate:"dive|required|min=3"`
}
```

### 错误处理

示例展示了两种错误处理方式：

```go
err := v.Validate(data)
if err != nil {
    // 方式1: ValidationResult（多个错误）
    if res, ok := err.(*schema.ValidationResult); ok {
        for _, e := range res.Errors() {
            fmt.Printf("%s: %s\n", e.Path, e.Name)
        }
    }
    // 方式2: ValidationError（单个错误）
    if e, ok := err.(*schema.ValidationError); ok {
        fmt.Printf("%s: %s\n", e.Path, e.Name)
    }
}
```

## 了解更多

查看完整文档：
- [README.md](../README.md) - 快速开始和基础功能
- [CLAUDE.md](../CLAUDE.md) - 完整的 API 参考和高级功能
