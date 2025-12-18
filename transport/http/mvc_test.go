package http_test

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/mime"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/mvc"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func TestRouteSuccess(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldCompression(), test.WithWorldHTTP())
	world.Register()
	world.RequireStart()

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
	header.Set(content.TypeKey, mime.HTMLMediaType)

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, mime.HTMLMediaType, res.Header.Get(content.TypeKey))

	_, err = html.Parse(strings.NewReader(body))
	require.NoError(t, err)

	world.RequireStop()
}

func TestRoutePartialViewSuccess(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldCompression(), test.WithWorldHTTP())
	world.Register()
	world.RequireStart()

	_, partial := mvc.NewViewPair("views/hello.tmpl")

	mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return partial, &test.Model, nil
	})

	header := http.Header{}
	header.Set(content.TypeKey, mime.HTMLMediaType)

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, mime.HTMLMediaType, res.Header.Get(content.TypeKey))

	_, err = html.Parse(strings.NewReader(body))
	require.NoError(t, err)

	world.RequireStop()
}

func TestRouteError(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRoundTripper(http.DefaultTransport), test.WithWorldHTTP())
	world.Register()
	world.RequireStart()

	mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return mvc.NewFullView("views/error.tmpl"), &test.Model, status.ServiceUnavailableError(test.ErrInternal)
	})

	header := http.Header{}
	header.Set(content.TypeKey, mime.HTMLMediaType)

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	require.Equal(t, mime.HTMLMediaType, res.Header.Get(content.TypeKey))

	_, err = html.Parse(strings.NewReader(body))
	require.NoError(t, err)

	world.RequireStop()
}

func TestStaticFileSuccess(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
	world.Register()
	world.RequireStart()

	mvc.StaticFile("/robots.txt", "static/robots.txt")

	header := http.Header{}
	header.Set(content.TypeKey, mime.TextMediaType)

	url := world.PathServerURL("http", "robots.txt")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, mime.TextMediaType, res.Header.Get(content.TypeKey))

	world.RequireStop()
}

func TestStaticFileError(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
	world.Register()
	world.RequireStart()

	mvc.StaticFile("/robots.txt", "static/bob.txt")

	header := http.Header{}
	url := world.PathServerURL("http", "robots.txt")

	res, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)

	world.RequireStop()
}

func TestStaticPathValueSuccess(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
	world.Register()
	world.RequireStart()

	mvc.StaticPathValue("/{file}", "file", "static")

	header := http.Header{}
	header.Set(content.TypeKey, mime.TextMediaType)

	url := world.PathServerURL("http", "robots.txt")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, mime.TextMediaType, res.Header.Get(content.TypeKey))

	world.RequireStop()
}

func TestStaticPathValueError(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
	world.Register()
	world.RequireStart()

	mvc.StaticPathValue("/{file}", "file", "static")

	header := http.Header{}
	url := world.PathServerURL("http", "bob.txt")

	res, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)

	world.RequireStop()
}

func TestMissingViews(t *testing.T) {
	mvc.Register(mvc.RegisterParams{
		Mux:         http.NewServeMux(),
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		Layout:      test.Layout,
	})

	require.False(t, mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return mvc.NewFullView("views/hello.tmpl"), &test.Model, nil
	}))
	require.False(t, mvc.StaticFile("/robots.txt", "static/robots.txt"))
	require.False(t, mvc.StaticPathValue("/{file}", "file", "static"))

	mvc.Register(mvc.RegisterParams{
		Mux:         http.NewServeMux(),
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  test.FileSystem,
	})

	require.False(t, mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return mvc.NewFullView("views/hello.tmpl"), &test.Model, nil
	}))
	require.False(t, mvc.StaticFile("/robots.txt", "static/robots.txt"))
	require.False(t, mvc.StaticPathValue("/{file}", "file", "static"))
}
