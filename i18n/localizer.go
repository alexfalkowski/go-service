package i18n

import (
	"github.com/alexfalkowski/go-service/v3/errors"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/nicksnyder/go-i18n/v2/i18n/template"
	"golang.org/x/text/language"
)

var messageParser = &template.TextParser{Option: "missingkey=error"}

// Localizer renders messages in one selected language, falling back to the
// catalog's default language for missing translations. Its results are plain
// text; the caller's HTML or Markdown renderer owns output escaping.
// Localizers and their catalog may be used concurrently.
type Localizer struct {
	localizer *i18n.Localizer
	fallback  *i18n.Localizer
	tag       language.Tag
}

// Tag returns the selected catalog language. Individual messages may use the
// default language when their translation is missing.
func (l *Localizer) Tag() language.Tag {
	return l.tag
}

// Text looks up a message without variables.
func (l *Localizer) Text(id string) (string, error) {
	return l.localize(&i18n.LocalizeConfig{MessageID: id})
}

// Message renders a message using data for its named template variables.
// Missing messages and template execution errors are returned to the caller.
func (l *Localizer) Message(id string, data any) (string, error) {
	return l.localize(&i18n.LocalizeConfig{MessageID: id, TemplateData: data})
}

// Plural selects a CLDR plural form for count and renders it with data.
// Count accepts go-i18n integer types and decimal strings. When data is nil,
// templates receive count as .PluralCount; otherwise data supplies all variables.
// Missing forms, invalid counts, and execution failures return errors.
func (l *Localizer) Plural(id string, count, data any) (string, error) {
	return l.localize(&i18n.LocalizeConfig{MessageID: id, PluralCount: count, TemplateData: data})
}

func (l *Localizer) localize(config *i18n.LocalizeConfig) (string, error) {
	config.TemplateParser = messageParser

	message, err := l.localizer.Localize(config)
	if _, missing := errors.AsType[*i18n.MessageNotFoundErr](err); missing {
		// Execute the default message again so a missing translation cannot mask
		// an execution error from the fallback template.
		return l.fallback.Localize(config)
	}

	return message, err
}
