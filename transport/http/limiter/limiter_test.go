package limiter_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/strings"
	httplimiter "github.com/alexfalkowski/go-service/v2/transport/http/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/limiter"
	"github.com/stretchr/testify/require"
)

func TestRoundTripperClosesBodyOnLimiterError(t *testing.T) {
	client, err := httplimiter.NewClientLimiter(test.NoopLifecycle{}, limiter.NewKeyMap(), test.NewLimiterConfig("user-agent", "1s", 1))
	require.NoError(t, err)
	require.NoError(t, client.Close(t.Context()))

	rt := httplimiter.NewRoundTripper(client, test.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("unexpected round trip")
		return nil, nil
	}))
	body := &test.TrackedBody{Reader: strings.NewReader("body")}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com/hello", body)
	require.NoError(t, err)

	res, err := rt.RoundTrip(req)
	require.Nil(t, res)
	require.Error(t, err)
	require.True(t, body.Closed)
}

func TestRoundTripperClosesBodyOnLimiterDenial(t *testing.T) {
	client, err := httplimiter.NewClientLimiter(test.NoopLifecycle{}, limiter.NewKeyMap(), test.NewLimiterConfig("user-agent", "1s", 1))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close(t.Context())) })

	rt := httplimiter.NewRoundTripper(client, test.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: http.Header{}}, nil
	}))
	first, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com/hello", http.NoBody)
	require.NoError(t, err)
	res, err := rt.RoundTrip(first)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())

	body := &test.TrackedBody{Reader: strings.NewReader("body")}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.com/hello", body)
	require.NoError(t, err)

	res, err = rt.RoundTrip(req)
	require.Nil(t, res)
	require.Error(t, err)
	require.True(t, body.Closed)
}
