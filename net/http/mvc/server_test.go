package mvc_test

import (
	"io/fs"
	"log/slog"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
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
	ctx := meta.WithContent(t.Context(), nil, &test.ErrResponseWriter{}, nil)

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
	req.Header.Set(content.TypeKey, media.HTML)
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
	req.Header.Set(content.TypeKey, media.HTML)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.Empty(t, res.Body.String())
}

func TestRouteRenderErrorDoesNotUseNotFoundController(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem: fstest.MapFS{
			"views/full.tmpl":    &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
			"views/partial.tmpl": &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
			"views/bad.tmpl":     &fstest.MapFile{Data: []byte(`{{ define "content" }}hello {{ index .Model 0 }}{{ end }}`)},
			"views/error.tmpl":   &fstest.MapFile{Data: []byte(`{{ define "content" }}error {{ index .Meta "mvcModelError" }}{{ end }}`)},
		},
		Pool:   test.Pool,
		Layout: mvc.NewLayout("views/full.tmpl", "views/partial.tmpl"),
	})

	require.True(t, mvc.NotFound(func(_ context.Context) (*mvc.View, *test.Page) {
		return mvc.NewFullView("views/error.tmpl"), nil
	}))
	require.True(t, mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return mvc.NewFullView("views/bad.tmpl"), &test.Model, nil
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.TypeKey, media.HTML)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.Empty(t, res.Body.String())
	require.Equal(t, media.WithUTF8(media.HTML), res.Header().Get(content.TypeKey))
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
	req.Header.Set(content.TypeKey, media.HTML)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.Empty(t, res.Body.String())
}

func TestNotFoundHandlesNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem: fstest.MapFS{
			"views/full.tmpl":    &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
			"views/partial.tmpl": &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
			"views/error.tmpl":   &fstest.MapFile{Data: []byte(`{{ define "content" }}{{ .Model.Code }} {{ .Model.Err }}{{ end }}`)},
		},
		Pool:   test.Pool,
		Layout: mvc.NewLayout("views/full.tmpl", "views/partial.tmpl"),
	})

	require.True(t, mvc.NotFound(func(_ context.Context) (*mvc.View, *errorModel) {
		return mvc.NewFullView("views/error.tmpl"), &errorModel{
			Code: http.StatusNotFound,
			Err:  status.Error(http.StatusNotFound, http.StatusText(http.StatusNotFound)),
		}
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/missing", http.NoBody)
	res := httptest.NewRecorder()

	mvc.NewHandler(mux).ServeHTTP(res, req)

	require.Equal(t, http.StatusNotFound, res.Code)
	require.Equal(t, media.WithUTF8(media.HTML), res.Header().Get(content.TypeKey))
	require.Contains(t, res.Body.String(), "404 Not Found")
}

func TestNotFoundIncludesRequestMeta(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem: fstest.MapFS{
			"views/full.tmpl":    &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
			"views/partial.tmpl": &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
			"views/error.tmpl":   &fstest.MapFile{Data: []byte(`{{ define "content" }}{{ .Model.Method }} {{ .Model.Path }}{{ end }}`)},
		},
		Pool:   test.Pool,
		Layout: mvc.NewLayout("views/full.tmpl", "views/partial.tmpl"),
	})

	require.True(t, mvc.NotFound(func(ctx context.Context) (*mvc.View, *requestModel) {
		req := meta.Request(ctx)
		return mvc.NewFullView("views/error.tmpl"), &requestModel{
			Method: req.Method,
			Path:   req.URL.Path,
		}
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/missing", http.NoBody)
	res := httptest.NewRecorder()

	mvc.NewHandler(mux).ServeHTTP(res, req)

	require.Equal(t, http.StatusNotFound, res.Code)
	require.Contains(t, res.Body.String(), "PUT /missing")
}

func TestNotFoundUsesDefaultWhenControllerMissing(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  test.FileSystem,
		Pool:        test.Pool,
		Layout:      test.Layout,
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/missing", http.NoBody)
	res := httptest.NewRecorder()

	mvc.NewHandler(mux).ServeHTTP(res, req)

	require.Equal(t, http.StatusNotFound, res.Code)
	require.Contains(t, res.Body.String(), "404 page not found")
}

func TestNotFoundWritesRenderStatusWhenViewMissing(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  test.FileSystem,
		Pool:        test.Pool,
		Layout:      test.Layout,
	})

	require.True(t, mvc.NotFound(func(_ context.Context) (*mvc.View, *test.Page) {
		return nil, &test.Model
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/missing", http.NoBody)
	res := httptest.NewRecorder()

	mvc.NewHandler(mux).ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.Empty(t, res.Body.String())
}

func TestNotFoundDoesNotReplaceMethodNotAllowed(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem: fstest.MapFS{
			"views/full.tmpl":    &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
			"views/partial.tmpl": &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
			"views/error.tmpl":   &fstest.MapFile{Data: []byte(`{{ define "content" }}custom {{ .Model.Code }}{{ end }}`)},
			"views/hello.tmpl":   &fstest.MapFile{Data: []byte(`{{ define "content" }}hello{{ end }}`)},
		},
		Pool:   test.Pool,
		Layout: mvc.NewLayout("views/full.tmpl", "views/partial.tmpl"),
	})

	require.True(t, mvc.NotFound(func(_ context.Context) (*mvc.View, *errorModel) {
		return mvc.NewFullView("views/error.tmpl"), &errorModel{Code: http.StatusNotFound}
	}))
	require.True(t, mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return mvc.NewFullView("views/hello.tmpl"), &test.Model, nil
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", http.NoBody)
	res := httptest.NewRecorder()

	mvc.NewHandler(mux).ServeHTTP(res, req)

	require.Equal(t, http.StatusMethodNotAllowed, res.Code)
	require.NotContains(t, res.Body.String(), "custom")
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

type errorModel struct {
	Err  error
	Code int
}

type requestModel struct {
	Method string
	Path   string
}
