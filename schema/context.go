package schema

import (
	"strings"

	"github.com/weilence/schema-validator/data"
)

// Context 封装验证的所有上下文信息
type Context struct {
	schema   Schema
	accessor data.Accessor

	// 上下文信息
	parent *Context
	path   string
}

// NewContext 创建根 context
func NewContext(schema Schema, accessor data.Accessor) *Context {
	ctx := &Context{
		schema:   schema,
		accessor: accessor,

		path: "",
	}

	return ctx
}

// WithChild 创建子 context（用于字段/元素验证）
func (c *Context) WithChild(segment string, childSchema Schema, childAccessor data.Accessor) *Context {
	newPath := segment
	if c.path != "" {
		if strings.HasPrefix(segment, "[") {
			newPath = c.path + segment
		} else {
			newPath = c.path + "." + segment
		}
	}

	return &Context{
		schema:   childSchema,
		accessor: childAccessor,

		parent: c,
		path:   newPath,
	}
}

// Schema 返回当前 schema
func (c *Context) Schema() Schema {
	return c.schema
}

// Accessor 返回当前 accessor
func (c *Context) Accessor() data.Accessor {
	return c.accessor
}

// Path 返回当前路径
func (c *Context) Path() string {
	return c.path
}

func (c *Context) Value() *data.Value {
	v, err := c.accessor.GetValue("")
	if err != nil {
		panic(err)
	}

	return v
}

func (c *Context) GetValue(path string) (*data.Value, error) {
	return c.accessor.GetValue(path)
}

func (c *Context) Parent() *Context {
	return c.parent
}
