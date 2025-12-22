package validation

import (
	"strings"

	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/errors"
)

// Context 封装验证的所有上下文信息
type Context struct {
	// 数据访问
	root     data.Accessor
	accessor data.Accessor // 当前正在验证的数据

	// 上下文信息
	path   string
	parent *Context // 父级 Context（自动追踪）
	values map[string]interface{}

	// 当前 Schema（用于 ModifySchema）
	schema interface{} // 可以是 *ObjectSchema、*FieldSchema、*ArraySchema

	// 懒加载缓存的类型化 accessors
	cachedField  data.FieldAccessor
	cachedObject data.ObjectAccessor
	cachedArray  data.ArrayAccessor
	cacheErr     error
	initialized  bool
}

// NewContext 创建根 context
func NewContext(accessor data.Accessor) *Context {
	return &Context{
		root:        accessor,
		accessor:    accessor,
		path:        "",
		parent:      nil,
		values:      make(map[string]interface{}),
		initialized: false,
	}
}

// WithChild 创建子 context（用于字段/元素验证）
func (c *Context) WithChild(segment string, childAccessor data.Accessor, childSchema interface{}) *Context {
	newPath := segment
	if c.path != "" {
		if strings.HasPrefix(segment, "[") {
			newPath = c.path + segment
		} else {
			newPath = c.path + "." + segment
		}
	}

	return &Context{
		root:        c.root,
		accessor:    childAccessor,
		path:        newPath,
		parent:      c, // 自动追踪父级
		values:      c.values, // 共享 values map
		schema:      childSchema,
		initialized: false,
	}
}

// --- 数据访问方法 ---

// Root 返回根 accessor
func (c *Context) Root() data.Accessor {
	return c.root
}

// Accessor 返回当前 accessor
func (c *Context) Accessor() data.Accessor {
	return c.accessor
}

// Path 返回当前路径
func (c *Context) Path() string {
	return c.path
}

// Parent 返回父级 context
func (c *Context) Parent() *Context {
	return c.parent
}

// Kind 返回 accessor 类型
func (c *Context) Kind() data.DataKind {
	return c.accessor.Kind()
}

// IsNil 检查 accessor 是否为 nil
func (c *Context) IsNil() bool {
	return c.accessor.IsNil()
}

// --- Schema 访问 ---

// Schema 返回当前 schema
func (c *Context) Schema() interface{} {
	return c.schema
}

// SetSchema 设置当前 schema（用于初始化）
func (c *Context) SetSchema(schema interface{}) {
	c.schema = schema
}

// ObjectSchema 返回类型化的 ObjectSchema（用于 ModifySchema）
// 注意：这里使用 interface{} 以避免循环依赖，调用者需要类型断言
func (c *Context) ObjectSchema() interface{} {
	return c.schema
}

// --- 懒加载缓存的类型化访问 ---

func (c *Context) initCache() {
	if c.initialized {
		return
	}

	kind := c.accessor.Kind()
	switch kind {
	case data.KindObject:
		c.cachedObject, c.cacheErr = c.accessor.AsObject()
	case data.KindArray:
		c.cachedArray, c.cacheErr = c.accessor.AsArray()
	default:
		c.cachedField, c.cacheErr = c.accessor.AsField()
	}

	c.initialized = true
}

func (c *Context) AsField() (data.FieldAccessor, error) {
	c.initCache()
	return c.cachedField, c.cacheErr
}

func (c *Context) AsObject() (data.ObjectAccessor, error) {
	c.initCache()
	return c.cachedObject, c.cacheErr
}

func (c *Context) AsArray() (data.ArrayAccessor, error) {
	c.initCache()
	return c.cachedArray, c.cacheErr
}

// --- 跨字段访问辅助方法 ---

// GetParentObject 获取父级 ObjectAccessor（用于跨字段验证）
func (c *Context) GetParentObject() (data.ObjectAccessor, error) {
	if c.parent == nil {
		return nil, errors.NewValidationError(c.path, "no_parent", nil)
	}
	return c.parent.AsObject()
}

// GetField 从父级或根获取字段
func (c *Context) GetField(path string) (data.Accessor, bool) {
	if c.parent != nil {
		if obj, err := c.parent.AsObject(); err == nil {
			return obj.GetField(path)
		}
	}

	if obj, err := c.root.AsObject(); err == nil {
		return obj.GetField(path)
	}

	return nil, false
}

// --- 自定义值 ---

func (c *Context) Set(key string, value interface{}) {
	c.values[key] = value
}

func (c *Context) Get(key string) (interface{}, bool) {
	val, ok := c.values[key]
	return val, ok
}
