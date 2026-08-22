package http_test

import (
	"log/slog"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/mvc"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func TestRouteRendersSuccessfulView(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldCompression(), test.WithWorldHTTP())

	full, _ := world.MVCServer.NewViewPair("views/hello.tmpl")

	controller := func(_ context.Context) (*mvc.View, *test.Page, error) {
		return full, &test.Model, nil
	}

	world.MVCServer.Delete("/hello", controller)
	world.MVCServer.Get("/hello", controller)
	world.MVCServer.Post("/hello", controller)
	world.MVCServer.Put("/hello", controller)
	world.MVCServer.Patch("/hello", controller)

	header := http.Header{}
	header.Set(http.ContentTypeKey, media.HTML)

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "text/html; charset=utf-8", res.Header.Get(http.ContentTypeKey))

	_, err = html.Parse(strings.NewReader(body))
	require.NoError(t, err)
}

func TestRoutePartialViewSuccess(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldCompression(), test.WithWorldHTTP())

	_, partial := world.MVCServer.NewViewPair("views/hello.tmpl")

	world.MVCServer.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return partial, &test.Model, nil
	})

	header := http.Header{}
	header.Set(http.ContentTypeKey, media.HTML)

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "text/html; charset=utf-8", res.Header.Get(http.ContentTypeKey))

	_, err = html.Parse(strings.NewReader(body))
	require.NoError(t, err)
}

func TestRouteRendersErrorView(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRoundTripper(http.DefaultTransport), test.WithWorldHTTP())

	view := world.MVCServer.NewFullView("views/error.tmpl")
	world.MVCServer.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return view, &test.Model, status.ServiceUnavailableError(test.ErrInternal)
	})

	header := http.Header{}
	header.Set(http.ContentTypeKey, media.HTML)

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	require.Equal(t, "text/html; charset=utf-8", res.Header.Get(http.ContentTypeKey))

	_, err = html.Parse(strings.NewReader(body))
	require.NoError(t, err)
}

func TestRouteRendersNotFoundView(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	view := world.MVCServer.NewFullView("views/error.tmpl")
	require.True(t, world.MVCServer.NotFound(func(_ context.Context) (*mvc.View, *mvc.Error) {
		return view, &mvc.Error{Code: http.StatusNotFound, Message: http.StatusText(http.StatusNotFound)}
	}))

	header := http.Header{}
	header.Set("Accept", media.HTML)
	header.Set(http.ContentTypeKey, media.HTML)

	url := world.PathServerURL("http", "missing")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
	require.Equal(t, "text/html; charset=utf-8", res.Header.Get(http.ContentTypeKey))

	_, err = html.Parse(strings.NewReader(body))
	require.NoError(t, err)
}

func TestNotFoundUsesContentFallbackWithoutHTMLAccept(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	view := world.MVCServer.NewFullView("views/error.tmpl")
	require.True(t, world.MVCServer.NotFound(func(_ context.Context) (*mvc.View, *mvc.Error) {
		return view, &mvc.Error{Code: http.StatusNotFound, Message: http.StatusText(http.StatusNotFound)}
	}))

	header := http.Header{}
	header.Set(http.ContentTypeKey, media.JSON)

	url := world.PathServerURL("http", "missing")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
	require.Equal(t, "text/error; charset=utf-8", res.Header.Get(http.ContentTypeKey))
	require.Equal(t, "http: not found", body)
}

func TestNotFoundHandlesHTMXRequest(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	view := world.MVCServer.NewPartialView("views/error.tmpl")
	require.True(t, world.MVCServer.NotFound(func(_ context.Context) (*mvc.View, *mvc.Error) {
		return view, &mvc.Error{Code: http.StatusNotFound, Message: http.StatusText(http.StatusNotFound)}
	}))

	header := http.Header{}
	header.Set("Accept", "*/*")
	header.Set("Hx-Request", "true")
	header.Set(http.ContentTypeKey, media.HTML)

	url := world.PathServerURL("http", "missing")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
	require.Equal(t, "text/html; charset=utf-8", res.Header.Get(http.ContentTypeKey))

	_, err = html.Parse(strings.NewReader(body))
	require.NoError(t, err)
}

func TestStaticFileRouteServesExistingFile(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	world.MVCServer.StaticFile("/robots.txt", "static/robots.txt")

	header := http.Header{}
	header.Set(http.ContentTypeKey, media.Text)

	url := world.PathServerURL("http", "robots.txt")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "text/plain; charset=utf-8", res.Header.Get(http.ContentTypeKey))
}

func TestStaticFileRouteReportsMissingFile(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	world.MVCServer.StaticFile("/robots.txt", "static/bob.txt")

	header := http.Header{}
	url := world.PathServerURL("http", "robots.txt")

	res, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestStaticPathRouteServesExistingFile(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	world.MVCServer.StaticPathValue("/{file}", "file", "static")

	header := http.Header{}
	header.Set(http.ContentTypeKey, media.Text)

	url := world.PathServerURL("http", "robots.txt")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "text/plain; charset=utf-8", res.Header.Get(http.ContentTypeKey))
}

func TestStaticPathRouteReportsMissingFile(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	world.MVCServer.StaticPathValue("/{file}", "file", "static")

	header := http.Header{}
	url := world.PathServerURL("http", "bob.txt")

	res, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestNewServerRejectsMissingViews(t *testing.T) {
	noFileSystem := mvc.NewServer(mvc.ServerParams{
		Router:      newTestRouter(),
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		Pool:        test.Pool,
		Layout:      test.Layout,
	})

	require.False(t, noFileSystem.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return nil, &test.Model, nil
	}))
	require.False(t, noFileSystem.NotFound(func(_ context.Context) (*mvc.View, *test.Page) {
		return nil, nil
	}))
	require.False(t, noFileSystem.StaticFile("/robots.txt", "static/robots.txt"))
	require.False(t, noFileSystem.StaticPathValue("/{file}", "file", "static"))

	noLayout := mvc.NewServer(mvc.ServerParams{
		Router:      newTestRouter(),
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  test.FileSystem,
		Pool:        test.Pool,
	})

	require.False(t, noLayout.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return nil, &test.Model, nil
	}))
	require.False(t, noLayout.NotFound(func(_ context.Context) (*mvc.View, *test.Page) {
		return nil, nil
	}))
	require.False(t, noLayout.StaticFile("/robots.txt", "static/robots.txt"))
	require.False(t, noLayout.StaticPathValue("/{file}", "file", "static"))
}

func newTestRouter() *http.Router {
	return http.NewRouter(http.NewServeMux(), http.NewRoutePolicy())
}
