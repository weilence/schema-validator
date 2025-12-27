package validators

import (
	"cmp"
	"fmt"

	"github.com/spf13/cast"
	"github.com/weilence/schema-validator/data"
)

func init() {
	RegisterDefault(defaultRegistry)
}

func RegisterDefault(r *Registry) {
	registerField(r)
	registerFormat(r)
	registerNetwork(r)
	registerOther(r)
}

type compareType int

func (ct compareType) String() string {
	switch ct {
	case LessThan:
		return "lt"
	case LessThanOrEqual:
		return "lte"
	case GreaterThan:
		return "gt"
	case GreaterThanOrEqual:
		return "gte"
	case Equal:
		return "eq"
	case NotEqual:
		return "neq"
	default:
		return "unknown"
	}
}

const (
	LessThan           compareType = iota // <
	LessThanOrEqual                       // <=
	GreaterThan                           // >
	GreaterThanOrEqual                    // >=
	Equal                                 // ==
	NotEqual                              // !=
)

func compareFn[T cmp.Ordered](t compareType, a, b T) bool {
	switch t {
	case LessThan:
		return a < b
	case LessThanOrEqual:
		return a <= b
	case GreaterThan:
		return a > b
	case GreaterThanOrEqual:
		return a >= b
	case Equal:
		return a == b
	case NotEqual:
		return a != b
	default:
		panic("unknown compare type")
	}
}

func compareValue(ct compareType, currentValue, otherValue *data.Value) (bool, error) {
	switch v := currentValue.Raw().(type) {
	case int, int8, int16, int32, int64:
		a, err := cast.ToE[int64](v)
		if err != nil {
			return false, err
		}

		b, err := cast.ToE[int64](otherValue.Raw())
		if err != nil {
			return false, err
		}

		return compareFn(ct, a, b), nil
	case uint, uint8, uint16, uint32, uint64:
		a, err := cast.ToE[uint64](v)
		if err != nil {
			return false, err
		}

		b, err := cast.ToE[uint64](otherValue.Raw())
		if err != nil {
			return false, err
		}

		return compareFn(ct, a, b), nil
	case float32, float64:
		a, err := cast.ToE[float64](v)
		if err != nil {
			return false, err
		}

		b, err := cast.ToE[float64](otherValue.Raw())
		if err != nil {
			return false, err
		}

		return compareFn(ct, a, b), nil
	case string:
		a, err := cast.ToE[string](v)
		if err != nil {
			return false, err
		}

		b, err := cast.ToE[string](otherValue.Raw())
		if err != nil {
			return false, err
		}

		return compareFn(ct, a, b), nil
	default:
		if currentValue.IsSliceOrArray()  {
			b := cast.ToInt(otherValue.Raw())
			return compareFn(ct, currentValue.Len(), b), nil
		}

		return false, fmt.Errorf("unsupported type for comparison")
	}
}
