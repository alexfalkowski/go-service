package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/mime"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/stretchr/testify/require"
)

func TestRestNoContent(t *testing.T) {
	for _, method := range []string{http.MethodDelete, http.MethodGet} {
		t.Run(method, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP())

			test.RegisterHandlers("/hello", test.RestNoContent)

			url := world.NamedServerURL("http", "hello")
			err := world.Rest.Do(t.Context(), method, url, rest.NoOptions)
			require.NoError(t, err)
		})
	}
}

func TestRestRequestNoContent(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP())

			test.RegisterRequestHandlers("/hello", test.RestRequestNoContent)

			url := world.NamedServerURL("http", "hello")
			req := &test.Request{Name: "test"}
			opts := &rest.Options{
				ContentType: mime.JSONMediaType,
				Request:     req,
			}
			err := world.Rest.Do(t.Context(), method, url, opts)
			require.NoError(t, err)
		})
	}
}

func TestRestError(t *testing.T) {
	for _, method := range []string{http.MethodDelete, http.MethodGet} {
		t.Run(method, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP(), test.WithWorldLoggerConfig("tint"))

			test.RegisterHandlers("/hello", test.RestError)

			url := world.NamedServerURL("http", "hello")
			err := world.Rest.Do(t.Context(), method, url, rest.NoOptions)
			require.Error(t, err)
		})
	}
}

func TestRestRequestError(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP())

			test.RegisterRequestHandlers("/hello", test.RestRequestError)

			url := world.NamedServerURL("http", "hello")
			req := &test.Request{Name: "test"}
			opts := &rest.Options{
				ContentType: mime.JSONMediaType,
				Request:     req,
			}
			err := world.Rest.Do(t.Context(), method, url, opts)
			require.Error(t, err)
		})
	}
}

func TestRestWithContent(t *testing.T) {
	for _, method := range []string{http.MethodDelete, http.MethodGet} {
		t.Run(method, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

			test.RegisterHandlers("/hello", test.RestContent)

			url := world.NamedServerURL("http", "hello")
			resp := &test.Response{}
			opts := &rest.Options{
				Response: resp,
			}
			err := world.Rest.Do(t.Context(), method, url, opts)
			require.NoError(t, err)
			require.Equal(t, "Hello Bob", resp.Greeting)
		})
	}
}

func TestRestRequestWithContent(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

			test.RegisterRequestHandlers("/hello", test.RestRequestContent)

			url := world.NamedServerURL("http", "hello")
			req := &test.Request{Name: "test"}
			resp := &test.Response{}
			opts := &rest.Options{
				ContentType: mime.JSONMediaType,
				Request:     req,
				Response:    resp,
			}
			err := world.Rest.Do(t.Context(), method, url, opts)
			require.NoError(t, err)
			require.Equal(t, "Hello test", resp.Greeting)
		})
	}
}

func TestRestInvalidStatusCode(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	test.RegisterHandlers("/hello", test.RestInvalidStatusCode)

	url := world.NamedServerURL("http", "hello")
	err := world.Rest.Get(t.Context(), url, rest.NoOptions)
	require.Error(t, err)

	err = world.Rest.Delete(t.Context(), url, rest.NoOptions)
	require.Error(t, err)

	test.RegisterRequestHandlers("/hello", test.RestRequestInvalidStatusCode)

	url = world.NamedServerURL("http", "hello")
	req := &test.Request{}
	opts := &rest.Options{Request: req}

	err = world.Rest.Post(t.Context(), url, opts)
	require.Error(t, err)

	err = world.Rest.Put(t.Context(), url, opts)
	require.Error(t, err)

	err = world.Rest.Patch(t.Context(), url, opts)
	require.Error(t, err)
}
