package i18n

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/weilence/schema-validator/schema"
	"golang.org/x/text/language"
)

func TestTranslator_Localize(t *testing.T) {
	translator, err := NewTranslator(language.English)
	assert.NoError(t, err)
	translator.Bundle().RegisterUnmarshalFunc("toml", toml.Unmarshal)

	translator.Bundle().MustAddMessages(language.English,
		&i18n.Message{ID: "greeting", Other: "Hello {{.Name}}"},
	)

	msg := translator.Localize("en", "greeting", map[string]any{"Name": "World"})
	if msg != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", msg)
	}
}

func TestTranslator_TranslateError(t *testing.T) {
	translator, err := NewTranslator(language.English)
	assert.NoError(t, err)

	translator.Bundle().MustAddMessages(language.English,
		&i18n.Message{ID: "required", Other: "This field is required"},
		&i18n.Message{ID: "min", Other: "Must be at least {{.Arg1}}"},
		&i18n.Message{ID: "eqfield", Other: "Must be equal to {{.Arg1}}"},
	)

	translator.Bundle().MustAddMessages(language.SimplifiedChinese,
		&i18n.Message{ID: "required", Other: "该字段为必填项"},
		&i18n.Message{ID: "min", Other: "最小值为 {{.Arg1}}"},
		&i18n.Message{ID: "eqfield", Other: "必须等于 {{.Arg1}}"},
	)

	tests := []struct {
		name     string
		lang     string
		err      schema.ValidationError
		expected string
	}{
		{
			name: "required in English",
			lang: "en",
			err: schema.ValidationError{
				Path:   "username",
				Code:   "required",
				Params: []any{},
			},
			expected: "This field is required",
		},
		{
			name: "required in Chinese",
			lang: "zh-CN",
			err: schema.ValidationError{
				Path:   "username",
				Code:   "required",
				Params: []any{},
			},
			expected: "该字段为必填项",
		},
		{
			name: "min with param in English",
			lang: "en",
			err: schema.ValidationError{
				Path:   "age",
				Code:   "min",
				Params: []any{18},
			},
			expected: "Must be at least 18",
		},
		{
			name: "min with param in Chinese",
			lang: "zh-CN",
			err: schema.ValidationError{
				Path:   "age",
				Code:   "min",
				Params: []any{18},
			},
			expected: "最小值为 18",
		},
		{
			name: "eqfield in English",
			lang: "en",
			err: schema.ValidationError{
				Path:   "confirmPassword",
				Code:   "eqfield",
				Params: []any{"password"},
			},
			expected: "Must be equal to password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := translator.TranslateError(tt.lang, tt.err)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
