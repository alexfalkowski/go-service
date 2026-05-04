package mvc_test

import (
	"io"
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
	hm "github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/mvc"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestStaticPathValueRejectsTraversal(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  test.FileSystem,
		Pool:        test.Pool,
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
		Pool:        test.Pool,
		Layout:      test.Layout,
	})

	view := mvc.NewFullView("views/hello.tmpl")
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	err := view.Render(ctx, &test.Model)

	require.ErrorIs(t, err, context.Canceled)
}

func TestViewRenderReturnsWriteError(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  test.FileSystem,
		Pool:        test.Pool,
		Layout:      test.Layout,
	})

	view := mvc.NewFullView("views/hello.tmpl")
	ctx := hm.WithContent(t.Context(), nil, &test.ErrResponseWriter{}, nil)

	err := view.Render(ctx, &test.Model)

	require.ErrorIs(t, err, test.ErrFailed)
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
		Pool:   test.Pool,
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
			"views/bad.tmpl":     &fstest.MapFile{Data: []byte(`{{ define "content" }}hello {{ index .Model 0 }}{{ end }}`)},
		},
		Pool:   test.Pool,
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
	require.Empty(t, res.Body.String())
}

func TestRouteErrorWritesRenderStatusWhenErrorViewFails(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem: fstest.MapFS{
			"views/full.tmpl":    &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
			"views/partial.tmpl": &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
			"views/bad.tmpl":     &fstest.MapFile{Data: []byte(`{{ define "content" }}hello {{ index .Model 0 }}{{ end }}`)},
		},
		Pool:   test.Pool,
		Layout: mvc.NewLayout("views/full.tmpl", "views/partial.tmpl"),
	})

	require.True(t, mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return mvc.NewFullView("views/bad.tmpl"), &test.Model, status.BadRequestError(fs.ErrInvalid)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.TypeKey, mime.HTMLMediaType)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.Empty(t, res.Body.String())
}

func TestStaticFileDoesNotWritePartialBodyWhenReadFails(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  errFileSystem{},
		Pool:        test.Pool,
		Layout:      test.Layout,
	})

	require.True(t, mvc.StaticFile("/asset", "asset.txt"))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/asset", http.NoBody)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.Empty(t, res.Body.String())
}

type errFileSystem struct{}

func (errFileSystem) Open(name string) (fs.File, error) {
	if name != "asset.txt" {
		return nil, fs.ErrNotExist
	}

	return &errFile{}, nil
}

type errFile struct {
	read bool
}

func (f *errFile) Stat() (fs.FileInfo, error) {
	return errFileInfo{}, nil
}

func (f *errFile) Read(p []byte) (int, error) {
	if f.read {
		return 0, io.EOF
	}

	f.read = true
	copy(p, "hello")

	return len("hello"), test.ErrFailed
}

func (f *errFile) Close() error {
	return nil
}

type errFileInfo struct{}

func (errFileInfo) Name() string {
	return "asset.txt"
}

func (errFileInfo) Size() int64 {
	return 5
}

func (errFileInfo) Mode() fs.FileMode {
	return 0
}

func (errFileInfo) ModTime() time.Time {
	return time.Time{}
}

func (errFileInfo) IsDir() bool {
	return false
}

func (errFileInfo) Sys() any {
	return nil
}
