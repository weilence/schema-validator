package validators

import (
	"fmt"
	"regexp"

	"github.com/weilence/schema-validator/schema"
)

func registerFormat(r *Registry)  {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	r.Register("email", func(ctx *schema.Context) error {
		field, err := ctx.Value()
		if err != nil {
			return nil
		}
		str := field.String()
		if !emailRegex.MatchString(str) {
			return schema.NewValidationError(ctx.Path(), "email", nil)
		}
		return nil
	})

	r.Register("pattern", func(ctx *schema.Context, params []string) error {
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
			return schema.NewValidationError(ctx.Path(), "pattern", map[string]any{"pattern": regex.String()})
		}
		return nil
	})
}
