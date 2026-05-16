package http_test

import (
	"log/slog"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/mvc"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func TestRouteSuccess(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldCompression(), test.WithWorldHTTP())

	full, _ := mvc.NewViewPair("views/hello.tmpl")

	controller := func(_ context.Context) (*mvc.View, *test.Page, error) {
		return full, &test.Model, nil
	}

	mvc.Delete("/hello", controller)
	mvc.Get("/hello", controller)
	mvc.Post("/hello", controller)
	mvc.Put("/hello", controller)
	mvc.Patch("/hello", controller)

	header := http.Header{}
	header.Set(content.TypeKey, media.HTML)

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, media.WithUTF8(media.HTML), res.Header.Get(content.TypeKey))

	_, err = html.Parse(strings.NewReader(body))
	require.NoError(t, err)
}

func TestRoutePartialViewSuccess(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldCompression(), test.WithWorldHTTP())

	_, partial := mvc.NewViewPair("views/hello.tmpl")

	mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return partial, &test.Model, nil
	})

	header := http.Header{}
	header.Set(content.TypeKey, media.HTML)

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, media.WithUTF8(media.HTML), res.Header.Get(content.TypeKey))

	_, err = html.Parse(strings.NewReader(body))
	require.NoError(t, err)
}

func TestRouteError(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRoundTripper(http.DefaultTransport), test.WithWorldHTTP())

	mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return mvc.NewFullView("views/error.tmpl"), &test.Model, status.ServiceUnavailableError(test.ErrInternal)
	})

	header := http.Header{}
	header.Set(content.TypeKey, media.HTML)

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	require.Equal(t, media.WithUTF8(media.HTML), res.Header.Get(content.TypeKey))

	_, err = html.Parse(strings.NewReader(body))
	require.NoError(t, err)
}

func TestNotFound(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	require.True(t, mvc.NotFound(func(_ context.Context) (*mvc.View, *notFoundModel) {
		return mvc.NewFullView("views/error.tmpl"), &notFoundModel{Error: http.StatusText(http.StatusNotFound)}
	}))

	header := http.Header{}
	header.Set("Accept", media.HTML)
	header.Set(content.TypeKey, media.HTML)

	url := world.PathServerURL("http", "missing")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
	require.Equal(t, media.WithUTF8(media.HTML), res.Header.Get(content.TypeKey))

	_, err = html.Parse(strings.NewReader(body))
	require.NoError(t, err)
}

func TestNotFoundUsesContentFallbackWithoutHTMLAccept(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	require.True(t, mvc.NotFound(func(_ context.Context) (*mvc.View, *notFoundModel) {
		return mvc.NewFullView("views/error.tmpl"), &notFoundModel{Error: http.StatusText(http.StatusNotFound)}
	}))

	header := http.Header{}
	header.Set(content.TypeKey, media.JSON)

	url := world.PathServerURL("http", "missing")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
	require.Equal(t, media.WithUTF8(media.Error), res.Header.Get(content.TypeKey))
	require.Equal(t, http.StatusText(http.StatusNotFound), body)
}

func TestNotFoundHandlesHTMXRequest(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	require.True(t, mvc.NotFound(func(_ context.Context) (*mvc.View, *notFoundModel) {
		return mvc.NewPartialView("views/error.tmpl"), &notFoundModel{Error: http.StatusText(http.StatusNotFound)}
	}))

	header := http.Header{}
	header.Set("Accept", "*/*")
	header.Set("Hx-Request", "true")
	header.Set(content.TypeKey, media.HTML)

	url := world.PathServerURL("http", "missing")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
	require.Equal(t, media.WithUTF8(media.HTML), res.Header.Get(content.TypeKey))

	_, err = html.Parse(strings.NewReader(body))
	require.NoError(t, err)
}

func TestStaticFileSuccess(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	mvc.StaticFile("/robots.txt", "static/robots.txt")

	header := http.Header{}
	header.Set(content.TypeKey, media.Text)

	url := world.PathServerURL("http", "robots.txt")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, media.WithUTF8(media.Text), res.Header.Get(content.TypeKey))
}

func TestStaticFileError(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	mvc.StaticFile("/robots.txt", "static/bob.txt")

	header := http.Header{}
	url := world.PathServerURL("http", "robots.txt")

	res, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestStaticPathValueSuccess(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	mvc.StaticPathValue("/{file}", "file", "static")

	header := http.Header{}
	header.Set(content.TypeKey, media.Text)

	url := world.PathServerURL("http", "robots.txt")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, media.WithUTF8(media.Text), res.Header.Get(content.TypeKey))
}

func TestStaticPathValueError(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	mvc.StaticPathValue("/{file}", "file", "static")

	header := http.Header{}
	url := world.PathServerURL("http", "bob.txt")

	res, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestMissingViews(t *testing.T) {
	mvc.Register(mvc.RegisterParams{
		Mux:         http.NewServeMux(),
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		Pool:        test.Pool,
		Layout:      test.Layout,
	})

	require.False(t, mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return mvc.NewFullView("views/hello.tmpl"), &test.Model, nil
	}))
	require.False(t, mvc.NotFound(func(_ context.Context) (*mvc.View, *test.Page) {
		return mvc.NewFullView("views/error.tmpl"), nil
	}))
	require.False(t, mvc.StaticFile("/robots.txt", "static/robots.txt"))
	require.False(t, mvc.StaticPathValue("/{file}", "file", "static"))

	mvc.Register(mvc.RegisterParams{
		Mux:         http.NewServeMux(),
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  test.FileSystem,
		Pool:        test.Pool,
	})

	require.False(t, mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return mvc.NewFullView("views/hello.tmpl"), &test.Model, nil
	}))
	require.False(t, mvc.NotFound(func(_ context.Context) (*mvc.View, *test.Page) {
		return mvc.NewFullView("views/error.tmpl"), nil
	}))
	require.False(t, mvc.StaticFile("/robots.txt", "static/robots.txt"))
	require.False(t, mvc.StaticPathValue("/{file}", "file", "static"))
}

type notFoundModel struct {
	Error string
}
