package schema

import "strconv"

type Validator interface {
	Name() string
	Params() []string
	Validate(ctx *Context) error
}

// FieldValidator validates a single field value
type validator struct {
	name   string
	params []string
	fn     func(ctx *Context, params []string) error
}

func (v validator) Name() string {
	return v.name
}

func (v validator) Params() []string {
	return v.params
}

// Validate implements FieldValidator
func (v validator) Validate(ctx *Context) error {
	return v.fn(ctx, v.params)
}

// validatorFactory creates a validator from tag parameters
// params is a slice of parameter strings extracted from the tag
// For example: validate:"between=10,20" -> params = ["10", "20"]
type validatorFactory struct {
	name string
	fn   func(ctx *Context, params []string) error
}

func (vf validatorFactory) Build(params []string) Validator {
	return &validator{
		name:   vf.name,
		params: params,
		fn:     vf.fn,
	}
}

// validatorToMap converts a validator to a map representation
func validatorToMap(v Validator) map[string]any {
	result := map[string]any{}
	result["name"] = v.Name()
	if len(v.Params()) > 0 {
		result["value"] = convertParams(v.Params())
	}

	return result
}

// convertParams converts string params to numbers when appropriate for JSON output
func convertParams(params []string) any {
	if len(params) == 0 {
		return nil
	}

	if len(params) == 1 {
		s := params[0]
		isNum := true
		for i := 0; i < len(s); i++ {
			if s[i] < '0' || s[i] > '9' {
				isNum = false
				break
			}
		}
		if isNum && len(s) > 0 {
			if n, err := strconv.Atoi(s); err == nil {
				return n
			}
		}
		return s
	}

	out := make([]any, len(params))
	for i, s := range params {
		isNum := true
		for j := 0; j < len(s); j++ {
			if s[j] < '0' || s[j] > '9' {
				isNum = false
				break
			}
		}
		if isNum && len(s) > 0 {
			if n, err := strconv.Atoi(s); err == nil {
				out[i] = n
				continue
			}
		}
		out[i] = s
	}
	return out
}
