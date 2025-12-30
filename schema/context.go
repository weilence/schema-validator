package schema

import (
	"strings"

	"github.com/weilence/schema-validator/data"
)

// Context 封装验证的所有上下文信息
type Context struct {
	schema   Schema
	accessor data.Accessor
	skipRest bool

	// 上下文信息
	parent *Context
	path   contextPath

	// 收集的错误
	errs *ValidationErrors
}

type contextPath []string

func newContextPath(p contextPath, field string) contextPath {
	newPath := make(contextPath, len(p)+1)
	copy(newPath, p)
	newPath[len(p)] = field
	return newPath
}

func (p contextPath) String() string {
	var sb strings.Builder
	for i, segment := range p {
		if i > 0 && segment[0] != '[' {
			sb.WriteString(".")
		}

		sb.WriteString(segment)
	}

	return sb.String()
}

// NewContext 创建根 context
func NewContext(schema Schema, accessor data.Accessor) *Context {
	ctx := &Context{
		schema:   schema,
		accessor: accessor,
		errs:     &ValidationErrors{},
	}

	return ctx
}

// WithChild 创建子 context（用于字段/元素验证）
func (c *Context) WithChild(field string, childSchema Schema, childAccessor data.Accessor) *Context {
	newPath := newContextPath(c.path, field)

	return &Context{
		schema:   childSchema,
		accessor: childAccessor,

		parent: c,
		path:   newPath,
		errs:   c.errs,
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
	return contextPath(c.path).String()
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

func (c *Context) SkipRest() {
	c.skipRest = true
}

func (c *Context) AddError(err ValidationError) {
	c.errs.AddError(err)
}

func (c *Context) Errors() ValidationErrors {
	if c.errs == nil {
		return nil
	}

	return *c.errs
}
