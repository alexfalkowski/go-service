package mvc_test

import (
	"log/slog"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/i18n"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/mvc"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestRouteLocalizesViews(t *testing.T) {
	t.Parallel()

	catalog, err := i18n.NewCatalog(language.English, test.FileSystem, "locales/*.yaml")
	require.NoError(t, err)

	mux := http.NewServeMux()
	server := mvc.NewServer(mvc.ServerParams{
		Router:      newTestRouter(mux),
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  test.FileSystem,
		Pool:        test.Pool,
		Layout:      test.Layout,
		Catalog:     catalog,
	})
	full, partial := server.NewViewPair("views/localized.tmpl")
	for path, view := range map[string]*mvc.View{"/full": full, "/partial": partial} {
		require.True(t, server.Get(path, func(_ context.Context) (*mvc.View, *localizedPage, error) {
			return view, &localizedPage{Name: "<b>Alex</b>", Count: 2}, nil
		}))
	}

	for _, tt := range []struct {
		name             string
		accept           string
		preference       string
		language         string
		title            string
		greeting         string
		tasks            string
		additionalAccept []string
	}{
		{name: "default", language: "en", title: "My task list", greeting: "Hello", tasks: "2 tasks"},
		{name: "German", accept: "de", language: "de", title: "Meine Aufgabenliste", greeting: "Hallo", tasks: "2 Aufgaben"},
		{name: "regional", accept: "de-CH", language: "de", title: "Meine Aufgabenliste", greeting: "Hallo", tasks: "2 Aufgaben"},
		{name: "weighted", accept: "en;q=0.2, de;q=0.9", language: "de", title: "Meine Aufgabenliste", greeting: "Hallo", tasks: "2 Aufgaben"},
		{
			name: "multiple fields", accept: "en;q=0.1", additionalAccept: []string{"de;q=1"},
			language: "de", title: "Meine Aufgabenliste", greeting: "Hallo", tasks: "2 Aufgaben",
		},
		{
			name: "multiple fields reversed", accept: "de;q=1", additionalAccept: []string{"en;q=0.1"},
			language: "de", title: "Meine Aufgabenliste", greeting: "Hallo", tasks: "2 Aufgaben",
		},
		{name: "unsupported", accept: "ja", language: "en", title: "My task list", greeting: "Hello", tasks: "2 tasks"},
		{name: "invalid", accept: "!", language: "en", title: "My task list", greeting: "Hello", tasks: "2 tasks"},
		{name: "override", accept: "de", preference: "en", language: "en", title: "My task list", greeting: "Hello", tasks: "2 tasks"},
		{name: "regional override", accept: "en", preference: "de-CH", language: "de", title: "Meine Aufgabenliste", greeting: "Hallo", tasks: "2 Aufgaben"},
		{
			name: "regional override with secondary match", accept: "en, de;q=0.5", preference: "de-CH",
			language: "de", title: "Meine Aufgabenliste", greeting: "Hallo", tasks: "2 Aufgaben",
		},
		{name: "unsupported override", accept: "de", preference: "ja", language: "de", title: "Meine Aufgabenliste", greeting: "Hallo", tasks: "2 Aufgaben"},
		{name: "invalid override", accept: "de", preference: "!", language: "de", title: "Meine Aufgabenliste", greeting: "Hallo", tasks: "2 Aufgaben"},
		{name: "excluded language", accept: "de;q=0, en", language: "en", title: "My task list", greeting: "Hello", tasks: "2 tasks"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			for _, path := range []string{"/full", "/partial"} {
				ctx := mvc.WithLanguage(t.Context(), tt.preference)
				req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, http.NoBody)
				req.Header.Set(http.AcceptLanguageKey, tt.accept)
				for _, accept := range tt.additionalAccept {
					req.Header.Add(http.AcceptLanguageKey, accept)
				}
				res := httptest.NewRecorder()
				res.Header().Set(http.VaryKey, http.AcceptEncodingKey)

				mux.ServeHTTP(res, req)

				require.Equal(t, http.StatusOK, res.Code)
				require.Equal(t, tt.language, res.Header().Get(http.ContentLanguageKey))
				require.Equal(t, []string{http.AcceptEncodingKey, http.AcceptLanguageKey}, res.Header().Values(http.VaryKey))
				test.RequireResponseBodyContains(t, res, `lang="`+tt.language+`"`)
				test.RequireResponseBodyContains(t, res, tt.title)
				test.RequireResponseBodyContains(t, res, tt.greeting+", &lt;b&gt;Alex&lt;/b&gt;!")
				test.RequireResponseBodyContains(t, res, tt.tasks)
				test.RequireResponseBodyContains(t, res, "Available in English")
			}
		})
	}
}

func TestRouteLocalizesErrorViews(t *testing.T) {
	catalog, err := i18n.NewCatalog(language.English, test.FileSystem, "locales/*.yaml")
	require.NoError(t, err)

	mux := http.NewServeMux()
	server := mvc.NewServer(mvc.ServerParams{
		Router:      newTestRouter(mux),
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  test.FileSystem,
		Pool:        test.Pool,
		Layout:      test.Layout,
		Catalog:     catalog,
	})
	view := server.NewFullView("views/localized_error.tmpl")
	require.True(t, server.Get("/error", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return view, nil, status.BadRequestError(test.ErrFailed)
	}))
	require.True(t, server.NotFound(func(_ context.Context) (*mvc.View, *mvc.Error) {
		return view, &mvc.Error{Code: http.StatusNotFound}
	}))

	for path, code := range map[string]int{"/error": http.StatusBadRequest, "/missing": http.StatusNotFound} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, http.NoBody)
			req.Header.Set(http.AcceptLanguageKey, "de")
			res := httptest.NewRecorder()

			server.NewHandler(mux).ServeHTTP(res, req)

			require.Equal(t, code, res.Code)
			require.Equal(t, "de", res.Header().Get(http.ContentLanguageKey))
			test.RequireResponseBodyContains(t, res, "Meine Aufgabenliste")
		})
	}
}

func TestRouteRejectsLocalizationErrors(t *testing.T) {
	catalog, err := i18n.NewCatalog(language.English, test.FileSystem, "locales/*.yaml")
	require.NoError(t, err)

	for _, tt := range []struct {
		name   string
		action string
	}{
		{name: "missing message", action: `{{ .Localizer.Text "missing" }}`},
		{name: "missing variable", action: `{{ .Localizer.Message "greeting" .Model }}`},
		{name: "fallback missing variable", action: `{{ .Localizer.Message "fallback_greeting" .Model }}`},
		{name: "invalid plural count", action: `{{ .Localizer.Plural "tasks" "invalid" nil }}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			server := mvc.NewServer(mvc.ServerParams{
				Router:      newTestRouter(mux),
				FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
				FileSystem: fstest.MapFS{
					"full.tmpl": &fstest.MapFile{Data: []byte("before " + tt.action)},
					"page.tmpl": &fstest.MapFile{Data: nil},
				},
				Pool:    test.Pool,
				Layout:  mvc.NewLayout("full.tmpl", "full.tmpl"),
				Catalog: catalog,
			})
			view := server.NewFullView("page.tmpl")
			require.True(t, server.Get("/hello", func(_ context.Context) (*mvc.View, *map[string]string, error) {
				return view, &map[string]string{}, nil
			}))
			ctx := status.WithRequestError(t.Context())
			req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/hello", http.NoBody)
			req.Header.Set(http.AcceptLanguageKey, "de")
			res := httptest.NewRecorder()

			mux.ServeHTTP(res, req)

			require.Equal(t, http.StatusInternalServerError, res.Code)
			test.RequireEmptyResponseBody(t, res)
			require.Error(t, status.RequestError(ctx))
		})
	}
}

func TestViewRenderLocalizesWithoutRequest(t *testing.T) {
	catalog, err := i18n.NewCatalog(language.English, test.FileSystem, "locales/*.yaml")
	require.NoError(t, err)

	server := mvc.NewServer(mvc.ServerParams{
		Router:      newTestRouter(http.NewServeMux()),
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  test.FileSystem,
		Pool:        test.Pool,
		Layout:      test.Layout,
		Catalog:     catalog,
	})
	view := server.NewPartialView("views/localized.tmpl")
	res := httptest.NewRecorder()
	ctx := meta.WithRequestResponse(mvc.WithLanguage(t.Context(), "de"), nil, res)

	err = view.Render(ctx, &localizedPage{Name: "Alex", Count: 1})

	require.NoError(t, err)
	require.Equal(t, "de", res.Header().Get(http.ContentLanguageKey))
	test.RequireResponseBodyContains(t, res, "1 Aufgabe")
}

type localizedPage struct {
	Name  string
	Count int
}
