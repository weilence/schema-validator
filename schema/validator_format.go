package schema

import (
	"fmt"
	"regexp"
)

func init() {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	Register("email", func(ctx *Context, params []string) error {
		field, err := ctx.Value()
		if err != nil {
			return nil
		}
		str := field.String()
		if !emailRegex.MatchString(str) {
			return NewValidationError(ctx.Path(), "email", nil)
		}
		return nil
	})

	Register("pattern", func(ctx *Context, params []string) error {
		if len(params) == 0 {
			return nil
		}
		pattern := params[0]
		regex, err := regexp.Compile(pattern)
		if err != nil {
			panic(fmt.Sprintf("invalid regex pattern: %s", pattern))
		}
		field, err := ctx.Value()
		if err != nil {
			return nil
		}
		str := field.String()
		if !regex.MatchString(str) {
			return NewValidationError(ctx.Path(), "pattern", map[string]any{"pattern": regex.String()})
		}
		return nil
	})
}
