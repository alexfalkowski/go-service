package mvc_test

import (
	"io/fs"
	"log/slog"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/mime"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/mvc"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/stretchr/testify/require"
)

func TestStaticPathValueRejectsTraversal(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  test.FileSystem,
		Layout:      test.Layout,
	})
	require.True(t, mvc.StaticPathValue("/{file...}", "file", "static"))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/robots.txt", http.NoBody)
	handler, _ := mux.Handler(req)
	req.SetPathValue("file", "../views/hello.tmpl")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)
}

func TestViewRenderReturnsContextError(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  test.FileSystem,
		Layout:      test.Layout,
	})

	view := mvc.NewFullView("views/hello.tmpl")
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	err := view.Render(ctx, &test.Model)

	require.ErrorIs(t, err, context.Canceled)
}

func TestRouteErrorIncludesMetaInTemplate(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem: fstest.MapFS{
			"views/full.tmpl":    &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
			"views/partial.tmpl": &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
			"views/error.tmpl":   &fstest.MapFile{Data: []byte(`{{ define "content" }}{{ index .Meta "mvcModelError" }}{{ end }}`)},
		},
		Layout: mvc.NewLayout("views/full.tmpl", "views/partial.tmpl"),
	})

	require.True(t, mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return mvc.NewFullView("views/error.tmpl"), &test.Model, status.BadRequestError(fs.ErrInvalid)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.TypeKey, mime.HTMLMediaType)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)
	require.Contains(t, res.Body.String(), fs.ErrInvalid.Error())
}

func TestRouteWritesStatusWhenRenderFails(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem: fstest.MapFS{
			"views/full.tmpl":    &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
			"views/partial.tmpl": &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
			"views/bad.tmpl":     &fstest.MapFile{Data: []byte(`{{ define "content" }}{{ index .Model 0 }}{{ end }}`)},
		},
		Layout: mvc.NewLayout("views/full.tmpl", "views/partial.tmpl"),
	})

	require.True(t, mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return mvc.NewFullView("views/bad.tmpl"), &test.Model, nil
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.TypeKey, mime.HTMLMediaType)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
}
