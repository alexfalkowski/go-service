package meta_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
	httpmeta "github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/stretchr/testify/require"
)

func TestWithContent(t *testing.T) {
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/test", http.NoBody)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	enc := test.NewEncoder(nil)

	ctx := httpmeta.WithContent(t.Context(), req, res, enc)

	require.Same(t, req, httpmeta.Request(ctx))
	require.Same(t, res, httpmeta.Response(ctx))
	require.Same(t, enc, httpmeta.Encoder(ctx))
}

func TestWithContentAllowsPartialContent(t *testing.T) {
	res := httptest.NewRecorder()

	ctx := httpmeta.WithContent(t.Context(), nil, res, nil)

	require.Same(t, res, httpmeta.Response(ctx))
}

func TestRoundTripperAppendDoesNotOverwriteRequestID(t *testing.T) {
	roundTripper := httpmeta.NewRoundTripper(
		env.UserAgent("agent"),
		test.StaticIDGenerator("request-id"),
		test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			req.Header.Add("User-Agent", "next-agent")

			require.Equal(t, []string{"agent", "next-agent"}, req.Header.Values("User-Agent"))
			require.Equal(t, []string{"request-id"}, req.Header.Values("Request-Id"))

			return &http.Response{StatusCode: http.StatusOK, Header: http.Header{}, Body: http.NoBody}, nil
		}),
	)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)
	require.NoError(t, err)

	res, err := roundTripper.RoundTrip(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Empty(t, req.Header.Values("User-Agent"))
	require.Empty(t, req.Header.Values("Request-Id"))
}

func TestRoundTripperHandlesNilRequestHeader(t *testing.T) {
	roundTripper := httpmeta.NewRoundTripper(
		env.UserAgent("agent"),
		test.StaticIDGenerator("request-id"),
		test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			require.Equal(t, "agent", req.Header.Get("User-Agent"))
			require.Equal(t, "request-id", req.Header.Get("Request-Id"))

			return &http.Response{StatusCode: http.StatusOK, Header: http.Header{}, Body: http.NoBody}, nil
		}),
	)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)
	require.NoError(t, err)
	req.Header = nil

	res, err := roundTripper.RoundTrip(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Nil(t, req.Header)
}

func TestRoundTripperStoresServiceMethod(t *testing.T) {
	roundTripper := httpmeta.NewRoundTripper(
		env.UserAgent("agent"),
		test.StaticIDGenerator("request-id"),
		test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			require.Equal(t, meta.Ignored("/users/123"), meta.ServiceMethod(req.Context()))
			require.NotContains(t, meta.CamelStrings(req.Context(), meta.NoPrefix), meta.ServiceMethodKey)

			return &http.Response{StatusCode: http.StatusOK, Header: http.Header{}, Body: http.NoBody}, nil
		}),
	)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/users/123", http.NoBody)
	require.NoError(t, err)

	res, err := roundTripper.RoundTrip(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
}

func TestHandlerAppendDoesNotOverwriteRequestID(t *testing.T) {
	handler := httpmeta.NewHandler(env.Name("service"), env.UserAgent("agent"), env.Version("v1"), test.StaticIDGenerator("request-id"))
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/test", http.NoBody)
	require.NoError(t, err)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req, func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Add("Service-Version", "v2")

		require.Equal(t, []string{"1", "v2"}, res.Header().Values("Service-Version"))
		require.Equal(t, []string{"request-id"}, res.Header().Values("Request-Id"))
	})
}

func TestHandlerStoresServiceMethodFromPath(t *testing.T) {
	handler := httpmeta.NewHandler(env.Name("service"), env.UserAgent("agent"), env.Version("v1"), test.StaticIDGenerator("request-id"))
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/users/123", http.NoBody)
	require.NoError(t, err)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req, func(_ http.ResponseWriter, req *http.Request) {
		require.Equal(t, meta.Ignored("/users/123"), meta.ServiceMethod(req.Context()))
		require.NotContains(t, meta.CamelStrings(req.Context(), meta.NoPrefix), meta.ServiceMethodKey)
	})
}

func TestHandlerStoresServiceMethodFromPattern(t *testing.T) {
	handler := httpmeta.NewHandler(env.Name("service"), env.UserAgent("agent"), env.Version("v1"), test.StaticIDGenerator("request-id"))
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/users/123", http.NoBody)
	require.NoError(t, err)
	req.Pattern = "GET /users/{id}"
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req, func(_ http.ResponseWriter, req *http.Request) {
		require.Equal(t, meta.Ignored("GET /users/{id}"), meta.ServiceMethod(req.Context()))
		require.NotContains(t, meta.CamelStrings(req.Context(), meta.NoPrefix), meta.ServiceMethodKey)
	})
}

func TestHandlerStoresGeolocationAsIgnored(t *testing.T) {
	handler := httpmeta.NewHandler(env.Name("service"), env.UserAgent("agent"), env.Version("v1"), test.StaticIDGenerator("request-id"))
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/test", http.NoBody)
	require.NoError(t, err)
	req.Header.Set("Geolocation", "geo:47,11")
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req, func(_ http.ResponseWriter, req *http.Request) {
		geolocation := meta.Geolocation(req.Context())

		require.Equal(t, "geo:47,11", geolocation.Value())
		require.Empty(t, geolocation.String())
		require.NotContains(t, meta.CamelStrings(req.Context(), meta.NoPrefix), meta.GeolocationKey)
	})
}
