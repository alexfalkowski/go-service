package i18n_test

import (
	"testing"
	"testing/fstest"

	"github.com/alexfalkowski/go-service/v3/i18n"
	"github.com/alexfalkowski/go-service/v3/internal/test"
	"github.com/alexfalkowski/go-service/v3/io"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestCatalogLocalizesMessages(t *testing.T) {
	catalog, err := i18n.NewCatalog(language.English, test.FileSystem, "locales/*.yaml")
	require.NoError(t, err)

	localizer := catalog.Localizer("de-CH")
	require.Equal(t, language.German, localizer.Tag())

	text, err := localizer.Text("page_title")
	require.NoError(t, err)
	require.Equal(t, "Meine Aufgabenliste", text)

	text, err = localizer.Message("greeting", map[string]string{"Name": "Alex"})
	require.NoError(t, err)
	require.Equal(t, "Hallo, Alex!", text)

	for _, tt := range []struct {
		want  string
		count int
	}{
		{want: "1 Aufgabe", count: 1},
		{want: "2 Aufgaben", count: 2},
	} {
		text, err := localizer.Plural("tasks", tt.count, nil)
		require.NoError(t, err)
		require.Equal(t, tt.want, text)
	}

	text, err = localizer.Text("fallback")
	require.NoError(t, err)
	require.Equal(t, "Available in English", text)

	text, err = localizer.Plural("assigned_tasks", 2, map[string]any{"Name": "Alex", "Count": 2})
	require.NoError(t, err)
	require.Equal(t, "Alex has 2 tasks", text)

	_, err = localizer.Text("missing")
	require.ErrorContains(t, err, "missing")

	_, err = localizer.Message("greeting", map[string]string{})
	require.ErrorContains(t, err, "Name")

	_, err = localizer.Message("fallback_greeting", map[string]string{})
	require.ErrorContains(t, err, "Name")

	_, err = localizer.Plural("tasks", "invalid", nil)
	require.Error(t, err)
}

func TestCatalogLoadsSupportedFormats(t *testing.T) {
	for _, tt := range []struct {
		name string
		data string
	}{
		{name: "active.en.json", data: `{"greeting": "Hello"}`},
		{name: "active.en.toml", data: `greeting = "Hello"`},
		{name: "active.en.yaml", data: `greeting: Hello`},
		{name: "active.en.yml", data: `greeting: Hello`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fileSystem := fstest.MapFS{tt.name: &fstest.MapFile{Data: []byte(tt.data)}}
			catalog, err := i18n.NewCatalog(language.English, fileSystem, tt.name)
			require.NoError(t, err)

			text, err := catalog.Localizer().Text("greeting")
			require.NoError(t, err)
			require.Equal(t, "Hello", text)
		})
	}
}

func TestCatalogUsesDefaultForUnsupportedPreferences(t *testing.T) {
	fileSystem := fstest.MapFS{
		"active.en.yaml": &fstest.MapFile{Data: []byte("greeting: Hello")},
		"active.ru.yaml": &fstest.MapFile{Data: []byte("greeting: Привет")},
	}
	catalog, err := i18n.NewCatalog(language.Russian, fileSystem, "*.yaml")
	require.NoError(t, err)

	for _, tt := range []struct {
		language    language.Tag
		name        string
		greeting    string
		preferences []string
	}{
		{name: "no preferences", language: language.Russian, greeting: "Привет"},
		{name: "unsupported same script", preferences: []string{"fr"}, language: language.Russian, greeting: "Привет"},
		{name: "invalid preference", preferences: []string{"!"}, language: language.Russian, greeting: "Привет"},
		{name: "supported later source", preferences: []string{"fr", "en"}, language: language.English, greeting: "Hello"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			localizer := catalog.Localizer(tt.preferences...)
			require.Equal(t, tt.language, localizer.Tag())

			text, err := localizer.Text("greeting")
			require.NoError(t, err)
			require.Equal(t, tt.greeting, text)
		})
	}
}

func TestCatalogRejectsInvalidFiles(t *testing.T) {
	for _, tt := range []struct {
		fileSystem io.FS
		name       string
		want       string
		patterns   []string
	}{
		{fileSystem: test.FileSystem, name: "no patterns", patterns: nil, want: "no message files"},
		{fileSystem: test.FileSystem, name: "missing file", patterns: []string{"locales/missing.en.yaml"}, want: "file does not exist"},
		{fileSystem: test.FileSystem, name: "missing second pattern", patterns: []string{"locales/active.en.yaml", "missing.*"}, want: "file does not exist"},
		{fileSystem: test.FileSystem, name: "invalid glob", patterns: []string{"["}, want: "syntax error in pattern"},
		{fileSystem: test.FileSystem, name: "missing default", patterns: []string{"locales/active.de.yaml"}, want: "default language"},
		{
			fileSystem: fstest.MapFS{"active.en.yaml": &fstest.MapFile{Data: []byte("greeting: [")}},
			name:       "malformed catalog", patterns: []string{"*.yaml"}, want: "load",
		},
		{
			fileSystem: fstest.MapFS{"active.en.yaml": &fstest.MapFile{Data: []byte(`greeting: '{{.Name'`)}},
			name:       "malformed template", patterns: []string{"*.yaml"}, want: `message "greeting"`,
		},
		{
			fileSystem: fstest.MapFS{"active.en.yaml": &fstest.MapFile{Data: []byte("tasks:\n  one: '{{.Count'\n  other: tasks")}},
			name:       "malformed plural template", patterns: []string{"*.yaml"}, want: `message "tasks"`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			catalog, err := i18n.NewCatalog(language.English, tt.fileSystem, tt.patterns...)
			require.ErrorContains(t, err, tt.want)
			require.Nil(t, catalog)
		})
	}
}
