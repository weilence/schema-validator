package i18n

import (
	"embed"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/weilence/schema-validator/schema"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

//go:embed messages/*.yaml
var messageFiles embed.FS

// Translator translates validation errors to localized messages
type Translator struct {
	bundle     *i18n.Bundle
	localizers map[string]*i18n.Localizer
}

// NewTranslator creates a new Translator with the given default language
// It automatically loads embedded English and Chinese message files
func NewTranslator(defaultLang language.Tag) (*Translator, error) {
	bundle := i18n.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	
	// Load embedded message files
	_, err := bundle.LoadMessageFileFS(messageFiles, "messages/active.en.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to load message file: %w", err)
	}

	_, err = bundle.LoadMessageFileFS(messageFiles, "messages/active.zh-CN.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to load message file: %w", err)
	}
	
	return &Translator{
		bundle:     bundle,
		localizers: make(map[string]*i18n.Localizer),
	}, nil
}

// Bundle returns the underlying i18n.Bundle for loading message files
func (t *Translator) Bundle() *i18n.Bundle {
	return t.bundle
}

// Localize translates a message with the given ID and template data
func (t *Translator) Localize(lang, messageID string, templateData map[string]any) string {
	if loc, ok := t.localizers[lang]; ok {
		msg, _ := loc.Localize(&i18n.LocalizeConfig{
			MessageID:    messageID,
			TemplateData: templateData,
		})
		return msg
	}

	loc := i18n.NewLocalizer(t.bundle, lang)
	t.localizers[lang] = loc

	msg, _ := loc.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
	return msg
}

// TranslateError converts a ValidationError to a localized string
// Positional params are automatically mapped to Arg1, Arg2, Arg3, etc.
func (t *Translator) TranslateError(lang string, err schema.ValidationError) string {
	templateData := make(map[string]any)
	templateData["Path"] = err.Path

	for i, param := range err.Params {
		templateData[fmt.Sprintf("Arg%d", i+1)] = param
	}

	msg := t.Localize(lang, err.Code, templateData)
	if msg == "" {
		return err.Error()
	}
	return msg
}
