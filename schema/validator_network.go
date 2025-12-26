package schema

import (
	"net"
	"regexp"

	"github.com/miekg/dns"
)

func init() {
	Register("ip", func(ctx *Context, params []string) error {
		field, err := ctx.Value()
		if err != nil {
			return nil
		}

		val := field.String()
		if net.ParseIP(val) != nil {
			return nil
		}

		return NewValidationError(ctx.Path(), "ip", map[string]any{
			"actual": val,
		})
	})

	Register("port", func(ctx *Context, params []string) error {
		field, err := ctx.Value()
		if err != nil {
			return nil
		}

		val, err := field.Int()
		if err != nil {
			return nil
		}

		if val >= 1 && val <= 65535 {
			return nil
		}

		return NewValidationError(ctx.Path(), "port", map[string]any{
			"actual": val,
		})
	})

	Register("domain", func(ctx *Context, params []string) error {
		field, err := ctx.Value()
		if err != nil {
			return nil
		}

		val := field.String()
		_, ok := dns.IsDomainName(val)
		if ok {
			return nil
		}

		return NewValidationError(ctx.Path(), "domain", map[string]any{
			"actual": val,
		})
	})

	var urlRegex = regexp.MustCompile(`^https?://[^\s]+$`)
	Register("url", func(ctx *Context, params []string) error {
		field, err := ctx.Value()
		if err != nil {
			return nil
		}
		str := field.String()
		if !urlRegex.MatchString(str) {
			return NewValidationError(ctx.Path(), "url", nil)
		}
		return nil
	})
}
