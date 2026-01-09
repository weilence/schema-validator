package rule

import (
	"strings"
	"unicode"

	"github.com/weilence/schema-validator/schema"
)

func registerString(r *Registry) {
	// ------------------------ workaround from go-playground/validator ------------------------
	r.Register("alpha", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		for _, r := range str {
			if !unicode.IsLetter(r) {
				return schema.ErrCheckFailed
			}
		}
		return nil
	})

	r.Register("alphaspace", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		for _, r := range str {
			if !unicode.IsLetter(r) && !unicode.IsSpace(r) {
				return schema.ErrCheckFailed
			}
		}
		return nil
	})

	r.Register("alphanum", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		for _, r := range str {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				return schema.ErrCheckFailed
			}
		}
		return nil
	})

	r.Register("alphanumspace", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		for _, r := range str {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
				return schema.ErrCheckFailed
			}
		}
		return nil
	})

	r.Register("alphanumunicode", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		for _, r := range str {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				return schema.ErrCheckFailed
			}
		}
		return nil
	})

	r.Register("alphaunicode", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		for _, r := range str {
			if !unicode.IsLetter(r) {
				return schema.ErrCheckFailed
			}
		}
		return nil
	})

	r.Register("ascii", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		for _, r := range str {
			if r > 127 {
				return schema.ErrCheckFailed
			}
		}
		return nil
	})

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
		for _, sr := range str {
			if sr == r {
				return nil
			}
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
		for _, sr := range str {
			if sr == r {
				return schema.ErrCheckFailed
			}
		}
		return nil
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
			if r > 127 {
				return nil
			}
		}
		return schema.ErrCheckFailed
	})

	r.Register("number", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		for _, r := range str {
			if !unicode.IsDigit(r) {
				return schema.ErrCheckFailed
			}
		}
		return nil
	})

	r.Register("numeric", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		for _, r := range str {
			if !unicode.IsDigit(r) && r != '.' && r != '-' && r != '+' {
				return schema.ErrCheckFailed
			}
		}
		return nil
	})

	r.Register("printascii", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		for _, r := range str {
			if r > 127 || !unicode.IsPrint(r) {
				return schema.ErrCheckFailed
			}
		}
		return nil
	})

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
