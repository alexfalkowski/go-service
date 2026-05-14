package token_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/token"
	"github.com/stretchr/testify/require"
)

func TestRoundTripperDoesNotMutateRequest(t *testing.T) {
	roundTripper := token.NewRoundTripper(
		env.UserID("service-user"),
		staticGenerator("fresh-token"),
		roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			require.Equal(t, "Bearer fresh-token", req.Header.Get("Authorization"))

			return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: http.Header{}}, nil
		}),
	)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/hello", http.NoBody)
	require.NoError(t, err)

	res, err := roundTripper.RoundTrip(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Empty(t, req.Header.Values("Authorization"))
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type staticGenerator string

func (g staticGenerator) Generate(_, _ string) ([]byte, error) {
	return []byte(g), nil
}
