package validators

import (
	"net"
	"regexp"

	"github.com/miekg/dns"
	"github.com/weilence/schema-validator/schema"
)

func registerNetwork(r *Registry)  {
	r.Register("ip", func(ctx *schema.Context) error {
		field, err := ctx.Value()
		if err != nil {
			return nil
		}

		val := field.String()
		if net.ParseIP(val) != nil {
			return nil
		}

		return schema.NewValidationError(ctx.Path(), "ip", map[string]any{
			"actual": val,
		})
	})

	r.Register("port", func(ctx *schema.Context) error {
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

		return schema.NewValidationError(ctx.Path(), "port", map[string]any{
			"actual": val,
		})
	})

	r.Register("domain", func(ctx *schema.Context) error {
		field, err := ctx.Value()
		if err != nil {
			return nil
		}

		val := field.String()
		_, ok := dns.IsDomainName(val)
		if ok {
			return nil
		}

		return schema.NewValidationError(ctx.Path(), "domain", map[string]any{
			"actual": val,
		})
	})

	var urlRegex = regexp.MustCompile(`^https?://[^\s]+$`)
	r.Register("url", func(ctx *schema.Context) error {
		field, err := ctx.Value()
		if err != nil {
			return nil
		}
		str := field.String()
		if !urlRegex.MatchString(str) {
			return schema.NewValidationError(ctx.Path(), "url", nil)
		}
		return nil
	})
}
