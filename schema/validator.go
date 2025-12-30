package schema

type Validator interface {
	Name() string
	Params() []any
	Validate(ctx *Context) error
}
