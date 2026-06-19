package mvc_test

import (
	"io/fs"
	"log/slog"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
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

func TestStaticPathValueSetsContentType(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem: fstest.MapFS{
			"static/icon.svg": &fstest.MapFile{Data: []byte(`<svg xmlns="http://www.w3.org/2000/svg"></svg>`)},
		},
		Pool:   test.Pool,
		Layout: test.Layout,
	})
	require.True(t, mvc.StaticPathValue("/{file...}", "file", "static"))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/icon.svg", http.NoBody)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "image/svg+xml", res.Header().Get(content.TypeKey))
}

func TestStaticPathValueRejectsDirectory(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem: fstest.MapFS{
			"static/assets": &fstest.MapFile{Mode: fs.ModeDir},
		},
		Pool:   test.Pool,
		Layout: test.Layout,
	})
	require.True(t, mvc.StaticPathValue("/{file...}", "file", "static"))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/assets", http.NoBody)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusNotFound, res.Code)
}

func TestStaticFileSetsContentType(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem: fstest.MapFS{
			"static/icon.svg": &fstest.MapFile{Data: []byte(`<svg xmlns="http://www.w3.org/2000/svg"></svg>`)},
		},
		Pool:   test.Pool,
		Layout: test.Layout,
	})
	require.True(t, mvc.StaticFile("/icon.svg", "static/icon.svg"))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/icon.svg", http.NoBody)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "image/svg+xml", res.Header().Get(content.TypeKey))
}

func TestStaticFileRejectsDirectory(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem: fstest.MapFS{
			"static/assets": &fstest.MapFile{Mode: fs.ModeDir},
		},
		Pool:   test.Pool,
		Layout: test.Layout,
	})
	require.True(t, mvc.StaticFile("/assets", "static/assets"))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/assets", http.NoBody)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusNotFound, res.Code)
}

func TestStaticFileRejectsPermissionDenied(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  permissionFileSystem{},
		Pool:        test.Pool,
		Layout:      test.Layout,
	})
	require.True(t, mvc.StaticFile("/asset", "asset.txt"))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/asset", http.NoBody)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusForbidden, res.Code)
}

func TestStaticFileSetsCacheControl(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem: fstest.MapFS{
			"static/asset.txt": &fstest.MapFile{Data: []byte("hello")},
		},
		Pool:   test.Pool,
		Layout: test.Layout,
	})
	require.True(t, mvc.StaticFile("/asset.txt", "static/asset.txt", mvc.WithCacheControl("public, max-age=60")))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/asset.txt", http.NoBody)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "public, max-age=60", res.Header().Get("Cache-Control"))
	require.Empty(t, res.Header().Get("ETag"))
	test.RequireResponseBody(t, res, "hello")
}

func TestStaticFileUsesETagValidator(t *testing.T) {
	modified := time.Now().UTC()
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem: fstest.MapFS{
			"static/asset.txt": &fstest.MapFile{Data: []byte("hello"), ModTime: modified},
		},
		Pool:   test.Pool,
		Layout: test.Layout,
	})
	require.True(t, mvc.StaticFile("/asset.txt", "static/asset.txt", mvc.WithCacheValidators()))

	firstReq := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/asset.txt", http.NoBody)
	firstRes := httptest.NewRecorder()

	mux.ServeHTTP(firstRes, firstReq)

	etag := firstRes.Header().Get("ETag")
	require.Equal(t, http.StatusOK, firstRes.Code)
	require.NotEmpty(t, etag)
	require.GreaterOrEqual(t, len(etag), 2)
	require.Equal(t, "W/", etag[:2])
	require.Equal(t, modified.Format(http.TimeFormat), firstRes.Header().Get("Last-Modified"))
	test.RequireResponseBody(t, firstRes, "hello")

	for _, tt := range []struct {
		headers map[string]string
		name    string
		body    string
		code    int
	}{
		{headers: map[string]string{"If-None-Match": etag}, name: "etag", code: http.StatusNotModified},
		{headers: map[string]string{"If-None-Match": etag[2:]}, name: "strong-etag", code: http.StatusNotModified},
		{headers: map[string]string{"If-Modified-Since": modified.Format(http.TimeFormat)}, name: "modified", code: http.StatusNotModified},
		{headers: map[string]string{"If-Modified-Since": modified.AddDate(0, 0, -1).Format(http.TimeFormat)}, name: "stale", body: "hello", code: http.StatusOK},
	} {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/asset.txt", http.NoBody)
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}
			res := httptest.NewRecorder()

			mux.ServeHTTP(res, req)

			require.Equal(t, tt.code, res.Code)
			if tt.body == "" {
				test.RequireEmptyResponseBody(t, res)
				return
			}
			test.RequireResponseBody(t, res, tt.body)
		})
	}
}

func TestStaticPathValueUsesETagValidator(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem: fstest.MapFS{
			"static/asset.txt": &fstest.MapFile{Data: []byte("hello")},
		},
		Pool:   test.Pool,
		Layout: test.Layout,
	})
	require.True(t, mvc.StaticPathValue("/{file...}", "file", "static", mvc.WithCacheValidators()))

	firstReq := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/asset.txt", http.NoBody)
	firstRes := httptest.NewRecorder()

	mux.ServeHTTP(firstRes, firstReq)

	etag := firstRes.Header().Get("ETag")
	require.Equal(t, http.StatusOK, firstRes.Code)
	require.NotEmpty(t, etag)
	test.RequireResponseBody(t, firstRes, "hello")

	secondReq := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/asset.txt", http.NoBody)
	secondReq.Header.Set("If-None-Match", etag)
	secondRes := httptest.NewRecorder()

	mux.ServeHTTP(secondRes, secondReq)

	require.Equal(t, http.StatusNotModified, secondRes.Code)
	test.RequireEmptyResponseBody(t, secondRes)
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

func TestNewViewUsesLayoutPathWhenBasenameCollides(t *testing.T) {
	for _, tt := range []struct {
		view func() *mvc.View
		name string
		want string
	}{
		{
			view: func() *mvc.View { return mvc.NewFullView("pages/full.tmpl") },
			name: "full",
			want: "full layout full page",
		},
		{
			view: func() *mvc.View { return mvc.NewPartialView("pages/partial.tmpl") },
			name: "partial",
			want: "partial layout partial page",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mvc.Register(mvc.RegisterParams{
				Mux:         mux,
				FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
				FileSystem: fstest.MapFS{
					"views/full.tmpl":    &fstest.MapFile{Data: []byte(`full layout {{ block "content" . }}default{{ end }}`)},
					"views/partial.tmpl": &fstest.MapFile{Data: []byte(`partial layout {{ block "content" . }}default{{ end }}`)},
					"pages/full.tmpl":    &fstest.MapFile{Data: []byte(`view-root{{ define "content" }}full page{{ end }}`)},
					"pages/partial.tmpl": &fstest.MapFile{Data: []byte(`view-root{{ define "content" }}partial page{{ end }}`)},
				},
				Pool:   test.Pool,
				Layout: mvc.NewLayout("views/full.tmpl", "views/partial.tmpl"),
			})

			res := httptest.NewRecorder()
			ctx := meta.WithContent(t.Context(), nil, res, nil)

			err := tt.view().Render(ctx, &test.Model)

			require.NoError(t, err)
			test.RequireResponseBody(t, res, tt.want)
		})
	}
}

func TestLayoutNames(t *testing.T) {
	layout := mvc.NewLayout("views/full.tmpl", "views/partial.tmpl")

	require.Equal(t, "full.tmpl", layout.FullName())
	require.Equal(t, "partial.tmpl", layout.PartialName())
}

func TestRouteErrorIncludesSafeModelAndRawMetaInTemplate(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem: fstest.MapFS{
			"views/full.tmpl":    &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
			"views/partial.tmpl": &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
			"views/error.tmpl":   &fstest.MapFile{Data: []byte(`{{ define "content" }}{{ .Model.Code }} {{ .Model.Message }} {{ index .Meta "mvcModelError" }}{{ end }}`)},
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
	test.RequireResponseBodyContains(t, res, "400 http: bad request invalid argument")
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
	test.RequireEmptyResponseBody(t, res)
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
	test.RequireEmptyResponseBody(t, res)
	require.Equal(t, "text/html; charset=utf-8", res.Header().Get(content.TypeKey))
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
	test.RequireEmptyResponseBody(t, res)
}

func TestNotFoundHandlesNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem: fstest.MapFS{
			"views/full.tmpl":    &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
			"views/partial.tmpl": &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
			"views/error.tmpl":   &fstest.MapFile{Data: []byte(`{{ define "content" }}{{ .Model.Code }} {{ .Model.Message }}{{ end }}`)},
		},
		Pool:   test.Pool,
		Layout: mvc.NewLayout("views/full.tmpl", "views/partial.tmpl"),
	})

	require.True(t, mvc.NotFound(func(_ context.Context) (*mvc.View, *mvc.Error) {
		return mvc.NewFullView("views/error.tmpl"), &mvc.Error{
			Code:    http.StatusNotFound,
			Message: http.StatusText(http.StatusNotFound),
		}
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/missing", http.NoBody)
	res := httptest.NewRecorder()

	mvc.NewHandler(mux).ServeHTTP(res, req)

	require.Equal(t, http.StatusNotFound, res.Code)
	require.Equal(t, "text/html; charset=utf-8", res.Header().Get(content.TypeKey))
	test.RequireResponseBodyContains(t, res, "404 Not Found")
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
	test.RequireResponseBodyContains(t, res, "PUT /missing")
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
	test.RequireResponseBodyContains(t, res, "404 page not found")
}

func TestFallback(t *testing.T) {
	tests := []struct {
		name    string
		accept  string
		hx      string
		handled bool
	}{
		{name: "accepts html", accept: media.HTML, handled: true},
		{name: "accepts htmx", accept: "*/*", hx: "true", handled: true},
		{name: "ignores api", accept: media.JSON},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mvc.Register(mvc.RegisterParams{
				Mux:         http.NewServeMux(),
				FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
				FileSystem: fstest.MapFS{
					"views/full.tmpl":    &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
					"views/partial.tmpl": &fstest.MapFile{Data: []byte(`{{ block "content" . }}{{ end }}`)},
					"views/error.tmpl":   &fstest.MapFile{Data: []byte(`{{ define "content" }}{{ .Model.Code }} {{ .Model.Message }}{{ end }}`)},
				},
				Pool:   test.Pool,
				Layout: mvc.NewLayout("views/full.tmpl", "views/partial.tmpl"),
			})

			require.True(t, mvc.NotFound(func(_ context.Context) (*mvc.View, *mvc.Error) {
				return mvc.NewPartialView("views/error.tmpl"), &mvc.Error{
					Code:    http.StatusNotFound,
					Message: http.StatusText(http.StatusNotFound),
				}
			}))

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/missing", http.NoBody)
			req.Header.Set("Accept", tt.accept)
			req.Header.Set("Hx-Request", tt.hx)
			res := httptest.NewRecorder()

			require.Equal(t, tt.handled, mvc.NotFoundHandler()(res, req))
			if !tt.handled {
				return
			}

			require.Equal(t, http.StatusNotFound, res.Code)
			require.Equal(t, "text/html; charset=utf-8", res.Header().Get(content.TypeKey))
			test.RequireResponseBodyContains(t, res, "404 Not Found")
		})
	}
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
	test.RequireEmptyResponseBody(t, res)
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

	require.True(t, mvc.NotFound(func(_ context.Context) (*mvc.View, *mvc.Error) {
		return mvc.NewFullView("views/error.tmpl"), &mvc.Error{Code: http.StatusNotFound}
	}))
	require.True(t, mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return mvc.NewFullView("views/hello.tmpl"), &test.Model, nil
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", http.NoBody)
	res := httptest.NewRecorder()

	mvc.NewHandler(mux).ServeHTTP(res, req)

	require.Equal(t, http.StatusMethodNotAllowed, res.Code)
	test.RequireResponseBodyNotContains(t, res, "custom")
}

func TestStaticFileSetsContentLength(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  test.ErrFileSystem{},
		Pool:        test.Pool,
		Layout:      test.Layout,
	})

	require.True(t, mvc.StaticFile("/asset", "asset.txt"))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/asset", http.NoBody)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "6", res.Header().Get("Content-Length"))
	test.RequireResponseBody(t, res, "hello")
}

type requestModel struct {
	Method string
	Path   string
}

type permissionFileSystem struct{}

func (permissionFileSystem) Open(string) (fs.File, error) {
	return nil, fs.ErrPermission
}
