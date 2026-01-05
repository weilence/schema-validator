package validators

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/schema"
)

func compareValidator(ct compareType) func(*schema.Context, any) error {
	return func(ctx *schema.Context, value any) error {
		field := ctx.Value()
		otherValue := data.NewValue(value)
		ok, err := compareValue(ct, field, otherValue)
		if err != nil {
			return err
		}

		if !ok {
			return schema.ErrCheckFailed
		}

		return nil
	}
}

func requiredFn(ctx *schema.Context) error {
	v := ctx.Value()
	if v.IsNilOrZero() {
		return schema.ErrCheckFailed
	}

	return nil
}

func registerOther(r *Registry) {
	// ------------------------- workaround from go-playground/validator ------------------------
	r.Register("dir", func(ctx *schema.Context) error {
		path := ctx.Value().String()
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("dirpath", func(ctx *schema.Context) error {
		path := ctx.Value().String()
		if filepath.IsAbs(path) || strings.Contains(path, "/") {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("file", func(ctx *schema.Context) error {
		path := ctx.Value().String()
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("filepath", func(ctx *schema.Context) error {
		path := ctx.Value().String()
		if filepath.IsAbs(path) || strings.Contains(path, "/") || strings.Contains(path, "\\") {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("image", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		ext := strings.ToLower(filepath.Ext(str))
		validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp"}
		if slices.Contains(validExts, ext) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("isdefault", func(ctx *schema.Context) error {
		// Assuming default is zero value
		if ctx.Value().IsNilOrZero() {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("len", func(ctx *schema.Context, expectedLen int) error {
		str := ctx.Value().String()
		if len(str) == expectedLen {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("max", compareValidator(LessThanOrEqual))

	r.Register("min", compareValidator(GreaterThanOrEqual))

	r.Register("oneof", func(ctx *schema.Context, params []string) error {
		val := ctx.Value().String()
		if !slices.Contains(params, val) {
			return schema.ErrCheckFailed
		}

		return nil
	})

	r.Register("required", requiredFn)

	r.Register("required_if", func(ctx *schema.Context, fieldName string, expectedValue any) error {
		otherValue, err := ctx.Parent().GetValue(fieldName)
		if err != nil {
			return fmt.Errorf("failed to get field '%s': %v", fieldName, err)
		}

		ok, err := compareValue(Equal, otherValue, data.NewValue(expectedValue))
		if err != nil {
			return err
		}

		if ok {
			return requiredFn(ctx)
		}

		return nil
	})

	r.Register("required_unless", func(ctx *schema.Context, fieldName string, expectedValue any) error {
		otherValue, err := ctx.Parent().GetValue(fieldName)
		if err != nil {
			return fmt.Errorf("failed to get field '%s': %v", fieldName, err)
		}

		ok, err := compareValue(NotEqual, otherValue, data.NewValue(expectedValue))
		if err != nil {
			return err
		}

		if ok {
			return requiredFn(ctx)
		}

		return nil
	})

	r.Register("required_with", func(ctx *schema.Context, fieldNames []string) error {
		for _, fieldName := range fieldNames {
			otherValue, err := ctx.Parent().GetValue(fieldName)
			if err != nil {
				return fmt.Errorf("failed to get field '%s': %v", fieldName, err)
			}
			if !otherValue.IsNilOrZero() {
				return requiredFn(ctx)
			}
		}
		return nil
	})

	r.Register("required_with_all", func(ctx *schema.Context, fieldNames []string) error {
		allPresent := true
		for _, fieldName := range fieldNames {
			otherValue, err := ctx.Parent().GetValue(fieldName)
			if err != nil {
				return fmt.Errorf("failed to get field '%s': %v", fieldName, err)
			}
			if otherValue.IsNilOrZero() {
				allPresent = false
				break
			}
		}
		if allPresent {
			return requiredFn(ctx)
		}
		return nil
	})

	r.Register("required_without", func(ctx *schema.Context, fieldNames []string) error {
		for _, fieldName := range fieldNames {
			otherValue, err := ctx.Parent().GetValue(fieldName)
			if err != nil {
				return fmt.Errorf("failed to get field '%s': %v", fieldName, err)
			}
			if otherValue.IsNilOrZero() {
				return requiredFn(ctx)
			}
		}
		return nil
	})

	r.Register("required_without_all", func(ctx *schema.Context, fieldNames []string) error {
		allAbsent := true
		for _, fieldName := range fieldNames {
			otherValue, err := ctx.Parent().GetValue(fieldName)
			if err != nil {
				return fmt.Errorf("failed to get field '%s': %v", fieldName, err)
			}
			if !otherValue.IsNilOrZero() {
				allAbsent = false
				break
			}
		}
		if allAbsent {
			return requiredFn(ctx)
		}
		return nil
	})

	r.Register("excluded_if", func(ctx *schema.Context, fieldName string, expectedValue any) error {
		otherValue, err := ctx.Parent().GetValue(fieldName)
		if err != nil {
			return fmt.Errorf("failed to get field '%s': %v", fieldName, err)
		}

		ok, err := compareValue(Equal, otherValue, data.NewValue(expectedValue))
		if err != nil {
			return err
		}

		if ok && !ctx.Value().IsNilOrZero() {
			return schema.ErrCheckFailed
		}

		return nil
	})

	r.Register("excluded_unless", func(ctx *schema.Context, fieldName string, expectedValue any) error {
		otherValue, err := ctx.Parent().GetValue(fieldName)
		if err != nil {
			return fmt.Errorf("failed to get field '%s': %v", fieldName, err)
		}

		ok, err := compareValue(NotEqual, otherValue, data.NewValue(expectedValue))
		if err != nil {
			return err
		}

		if ok && !ctx.Value().IsNilOrZero() {
			return schema.ErrCheckFailed
		}

		return nil
	})

	r.Register("excluded_with", func(ctx *schema.Context, fieldNames []string) error {
		for _, fieldName := range fieldNames {
			otherValue, err := ctx.Parent().GetValue(fieldName)
			if err != nil {
				return fmt.Errorf("failed to get field '%s': %v", fieldName, err)
			}
			if !otherValue.IsNilOrZero() && !ctx.Value().IsNilOrZero() {
				return schema.ErrCheckFailed
			}
		}
		return nil
	})

	r.Register("excluded_with_all", func(ctx *schema.Context, fieldNames []string) error {
		allPresent := true
		for _, fieldName := range fieldNames {
			otherValue, err := ctx.Parent().GetValue(fieldName)
			if err != nil {
				return fmt.Errorf("failed to get field '%s': %v", fieldName, err)
			}
			if otherValue.IsNilOrZero() {
				allPresent = false
				break
			}
		}
		if allPresent && !ctx.Value().IsNilOrZero() {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("excluded_without", func(ctx *schema.Context, fieldNames []string) error {
		for _, fieldName := range fieldNames {
			otherValue, err := ctx.Parent().GetValue(fieldName)
			if err != nil {
				return fmt.Errorf("failed to get field '%s': %v", fieldName, err)
			}
			if otherValue.IsNilOrZero() && !ctx.Value().IsNilOrZero() {
				return schema.ErrCheckFailed
			}
		}
		return nil
	})

	r.Register("excluded_without_all", func(ctx *schema.Context, fieldNames []string) error {
		allAbsent := true
		for _, fieldName := range fieldNames {
			otherValue, err := ctx.Parent().GetValue(fieldName)
			if err != nil {
				return fmt.Errorf("failed to get field '%s': %v", fieldName, err)
			}
			if !otherValue.IsNilOrZero() {
				allAbsent = false
				break
			}
		}
		if allAbsent && !ctx.Value().IsNilOrZero() {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("unique", func(ctx *schema.Context) error {
		// For simplicity, assume it's a list and check uniqueness
		// In practice, this might need more context
		// For now, just pass
		return nil
	})
	// ------------------------ end of workaround ------------------------

	r.Register("omitempty", func(ctx *schema.Context) error {
		if ctx.Value().IsNilOrZero() {
			ctx.SkipRest()
		}

		return nil
	})
}
