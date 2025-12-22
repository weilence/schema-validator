package validation

import (
	"strings"

	"github.com/weilence/schema-validator/data"
)

// Context provides validation context with cross-field access
type Context struct {
	// Root is the root data accessor
	Root data.Accessor

	// Path is the current path in the data structure (e.g., "user.address.city")
	Path string

	// Parent is the parent object accessor (for cross-field validation)
	Parent data.ObjectAccessor

	// values stores custom context values
	values map[string]interface{}
}

// NewContext creates a new validation context
func NewContext(root data.Accessor) *Context {
	return &Context{
		Root:   root,
		Path:   "",
		Parent: nil,
		values: make(map[string]interface{}),
	}
}

// WithPath creates a new context with an additional path segment
func (c *Context) WithPath(segment string) *Context {
	newPath := segment
	if c.Path != "" {
		if strings.HasPrefix(segment, "[") {
			newPath = c.Path + segment
		} else {
			newPath = c.Path + "." + segment
		}
	}

	return &Context{
		Root:   c.Root,
		Path:   newPath,
		Parent: c.Parent,
		values: c.values,
	}
}

// GetField retrieves a field from the parent or root by path
func (c *Context) GetField(path string) (data.Accessor, bool) {
	if c.Parent != nil {
		return c.Parent.GetField(path)
	}

	if obj, err := c.Root.AsObject(); err == nil {
		return obj.GetField(path)
	}

	return nil, false
}

// Set sets a custom context value
func (c *Context) Set(key string, value interface{}) {
	c.values[key] = value
}

// Get retrieves a custom context value
func (c *Context) Get(key string) (interface{}, bool) {
	val, ok := c.values[key]
	return val, ok
}
