package meta_test

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/stretchr/testify/require"
)

func TestWithContent(t *testing.T) {
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/test", http.NoBody)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	enc := &encoder{}

	ctx := meta.WithContent(t.Context(), req, res, enc)

	require.Same(t, req, meta.Request(ctx))
	require.Same(t, res, meta.Response(ctx))
	require.Same(t, enc, meta.Encoder(ctx))
}

func TestWithContentAllowsPartialContent(t *testing.T) {
	res := httptest.NewRecorder()

	ctx := meta.WithContent(t.Context(), nil, res, nil)

	require.Same(t, res, meta.Response(ctx))
}

func TestRoundTripperAppendDoesNotOverwriteRequestID(t *testing.T) {
	roundTripper := meta.NewRoundTripper(
		env.UserAgent("agent"),
		staticGenerator("request-id"),
		roundTripperFunc(func(req *http.Request) (*http.Response, error) {
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

func TestHandlerAppendDoesNotOverwriteRequestID(t *testing.T) {
	handler := meta.NewHandler(env.UserAgent("agent"), env.Version("v1"), staticGenerator("request-id"))
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/test", http.NoBody)
	require.NoError(t, err)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req, func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Add("Service-Version", "v2")

		require.Equal(t, []string{"1", "v2"}, res.Header().Values("Service-Version"))
		require.Equal(t, []string{"request-id"}, res.Header().Values("Request-Id"))
	})
}

type encoder struct{}

func (e *encoder) Decode(_ io.Reader, _ any) error {
	return nil
}

func (e *encoder) Encode(_ io.Writer, _ any) error {
	return nil
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type staticGenerator string

func (g staticGenerator) Generate() string {
	return string(g)
}
