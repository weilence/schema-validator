package validators

import (
	"regexp"

	"github.com/weilence/schema-validator/schema"
)

func registerFormat(r *Registry) {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	r.Register("email", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if !emailRegex.MatchString(str) {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("pattern", func(ctx *schema.Context, pattern string) error {
		field := ctx.Value().String()

		matched, err := regexp.MatchString(pattern, field)
		if err != nil {
			return err
		}

		if !matched {
			return schema.ErrCheckFailed
		}

		return nil
	})
}
