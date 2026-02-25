package rule

import (
	"strings"
	"unicode"

	"github.com/weilence/schema-validator/schema"
)

func registerString(r *Registry) {
	// ------------------------ workaround from go-playground/validator ------------------------
	r.Register("alpha", stringValidator(unicode.IsLetter))

	r.Register("alphaspace", stringValidator(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsSpace(r)
	}))

	r.Register("alphanum", stringValidator(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	}))

	r.Register("alphanumspace", stringValidator(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r)
	}))

	r.Register("alphanumunicode", stringValidator(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	}))

	r.Register("alphaunicode", stringValidator(unicode.IsLetter))

	r.Register("ascii", stringValidator(func(r rune) bool {
		return r <= unicode.MaxASCII
	}))

	r.Register("boolean", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if str == "true" || str == "false" || str == "1" || str == "0" {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("contains", func(ctx *schema.Context, substr string) error {
		str := ctx.Value().String()
		if strings.Contains(str, substr) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("containsany", func(ctx *schema.Context, chars string) error {
		str := ctx.Value().String()
		if strings.ContainsAny(str, chars) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("containsrune", func(ctx *schema.Context, runeStr string) error {
		str := ctx.Value().String()
		if len(runeStr) == 0 {
			return schema.ErrCheckFailed
		}
		r := []rune(runeStr)[0]
		if strings.ContainsRune(str, r) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("endsnotwith", func(ctx *schema.Context, suffix string) error {
		str := ctx.Value().String()
		if !strings.HasSuffix(str, suffix) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("endswith", func(ctx *schema.Context, suffix string) error {
		str := ctx.Value().String()
		if strings.HasSuffix(str, suffix) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("excludes", func(ctx *schema.Context, substr string) error {
		str := ctx.Value().String()
		if !strings.Contains(str, substr) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("excludesall", func(ctx *schema.Context, chars string) error {
		str := ctx.Value().String()
		for _, c := range chars {
			if strings.ContainsRune(str, c) {
				return schema.ErrCheckFailed
			}
		}
		return nil
	})

	r.Register("excludesrune", func(ctx *schema.Context, runeStr string) error {
		str := ctx.Value().String()
		if len(runeStr) == 0 {
			return schema.ErrCheckFailed
		}
		r := []rune(runeStr)[0]
		if !strings.ContainsRune(str, r) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("lowercase", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if str == strings.ToLower(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("multibyte", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		for _, r := range str {
			if r > unicode.MaxASCII {
				return nil
			}
		}
		return schema.ErrCheckFailed
	})

	r.Register("number", stringValidator(unicode.IsDigit))

	r.Register("numeric", stringValidator(func(r rune) bool {
		return unicode.IsDigit(r) || r == '.' || r == '-' || r == '+'
	}))

	r.Register("printascii", stringValidator(func(r rune) bool {
		return r <= unicode.MaxASCII && unicode.IsPrint(r)
	}))

	r.Register("startsnotwith", func(ctx *schema.Context, prefix string) error {
		str := ctx.Value().String()
		if !strings.HasPrefix(str, prefix) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("startswith", func(ctx *schema.Context, prefix string) error {
		str := ctx.Value().String()
		if strings.HasPrefix(str, prefix) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("uppercase", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if str == strings.ToUpper(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})
	// ------------------------ end of workaround ------------------------
}

// runeCheckFunc defines a function that checks if a rune meets a condition
type runeCheckFunc func(rune) bool

// stringValidator validates a string by applying a rune check function to all characters
func stringValidator(checkFunc runeCheckFunc) func(*schema.Context) error {
	return func(ctx *schema.Context) error {
		str := ctx.Value().String()
		for _, r := range str {
			if !checkFunc(r) {
				return schema.ErrCheckFailed
			}
		}
		return nil
	}
}
