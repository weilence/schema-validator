package validators

import (
	"net"
	"regexp"

	"github.com/miekg/dns"
	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/schema"
)

func registerNetwork(r *Registry) {
	r.Register("ip", func(ctx *schema.Context) error {
		field := ctx.Value()
		val := field.String()
		if net.ParseIP(val) != nil {
			return nil
		}

		return schema.ErrCheckFailed
	})

	r.Register("port", func(ctx *schema.Context) error {
		field := ctx.Value()

		ok, err := compareValue(GreaterThanOrEqual, field, data.NewValue(1))
		if err != nil {
			return err
		}
		if !ok {
			return schema.ErrCheckFailed
		}

		ok, err = compareValue(LessThanOrEqual, field, data.NewValue(65535))
		if err != nil {
			return err
		}
		if !ok {
			return schema.ErrCheckFailed
		}

		return nil
	})

	r.Register("domain", func(ctx *schema.Context) error {
		field := ctx.Value()
		val := field.String()
		_, ok := dns.IsDomainName(val)
		if ok {
			return nil
		}

		return schema.ErrCheckFailed
	})

	var urlRegex = regexp.MustCompile(`^https?://[^\s]+$`)
	r.Register("url", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if !urlRegex.MatchString(str) {
			return schema.ErrCheckFailed
		}

		return nil
	})
}
