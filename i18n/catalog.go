package i18n

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/codec"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/encoding/toml"
	"github.com/alexfalkowski/go-service/v2/encoding/yaml"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// NewCatalog loads go-i18n JSON, TOML, or YAML message files from fileSystem.
// File names identify their language, for example active.en.yaml and active.de.yaml.
// Each io.Glob pattern must match at least one file, and the default language must
// have messages. File, catalog, and message template syntax errors prevent construction.
// Construct the catalog during startup and share it across requests.
func NewCatalog(lang language.Tag, fs io.FS, patterns ...string) (*Catalog, error) {
	bundle := i18n.NewBundle(lang)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	// go-i18n decodes into an untyped map, for which TOML cannot track consumed keys.
	bundle.RegisterUnmarshalFunc("toml", func(data []byte, value any) error {
		return toml.NewEncoder().Decode(bytes.NewReader(data), value, codec.WithDiscardUnknown())
	})
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	bundle.RegisterUnmarshalFunc("yml", yaml.Unmarshal)

	if len(patterns) == 0 {
		return nil, fmt.Errorf("i18n: no message files: %w", io.ErrNotExist)
	}

	defaultCount := 0
	for _, pattern := range patterns {
		names, err := io.Glob(fs, pattern)
		if err != nil {
			return nil, fmt.Errorf("i18n: message files %q: %w", pattern, err)
		}
		if len(names) == 0 {
			return nil, fmt.Errorf("i18n: message files %q: %w", pattern, io.ErrNotExist)
		}

		for _, name := range names {
			file, err := bundle.LoadMessageFileFS(fs, name)
			if err != nil {
				return nil, fmt.Errorf("i18n: load %q: %w", name, err)
			}

			if err := validateMessages(file); err != nil {
				return nil, err
			}

			if file.Tag == lang {
				defaultCount += len(file.Messages)
			}
		}
	}
	if defaultCount == 0 {
		return nil, fmt.Errorf("i18n: no messages for default language %q", lang)
	}

	tags := bundle.LanguageTags()

	return &Catalog{
		bundle:   bundle,
		fallback: i18n.NewLocalizer(bundle, lang.String()),
		tags:     tags,
		matcher:  language.NewMatcher(tags),
	}, nil
}

// Catalog owns immutable messages and language matching for an application.
// Applications supply their catalogs; no HTTP or process-global state is required.
type Catalog struct {
	bundle   *i18n.Bundle
	fallback *i18n.Localizer
	matcher  language.Matcher
	tags     []language.Tag
}

// Localizer selects a supported language from ordered preference sources, each
// of which may contain an Accept-Language value. It uses the first source with
// a supported match, ignoring invalid or unsupported sources. If no source
// matches, it selects the catalog's default language.
func (c *Catalog) Localizer(preferences ...string) *Localizer {
	// The bundle lists its default language first.
	tag := c.tags[0]
	for _, preference := range preferences {
		tags, _, err := language.ParseAcceptLanguage(preference)
		if err != nil {
			continue
		}

		_, index, confidence := c.matcher.Match(tags...)
		if confidence != language.No {
			tag = c.tags[index]
			break
		}
	}

	return &Localizer{
		localizer: i18n.NewLocalizer(c.bundle, tag.String()),
		fallback:  c.fallback,
		tag:       tag,
	}
}

func validateMessages(file *i18n.MessageFile) error {
	for _, message := range file.Messages {
		for _, text := range []string{message.Zero, message.One, message.Two, message.Few, message.Many, message.Other} {
			if _, err := messageParser.Parse(text, message.LeftDelim, message.RightDelim); err != nil {
				return fmt.Errorf("i18n: message %q in %q: %w", message.ID, file.Path, err)
			}
		}
	}

	return nil
}
