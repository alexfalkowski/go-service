package limiter_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
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

func TestHandlerStoresLimiterOnContextWhenAllowed(t *testing.T) {
	server, err := httplimiter.NewServerLimiter(test.NoopLifecycle{}, limiter.NewKeyMap(), test.NewLimiterConfig("user-agent", "1s", 1))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, server.Close(t.Context())) })

	handler := httplimiter.NewHandler(http.NewRoutePolicy(), server)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/hello", http.NoBody)
	require.NoError(t, err)
	res := httptest.NewRecorder()

	var called bool
	handler.ServeHTTP(res, req, func(_ http.ResponseWriter, req *http.Request) {
		called = true
		require.Same(t, server.Limiter, meta.Limiter(req.Context()))
	})

	require.True(t, called)
	require.Equal(t, http.StatusOK, res.Code)
}

func TestHandlerDoesNotStoreLimiterOnContextWhenDenied(t *testing.T) {
	server, err := httplimiter.NewServerLimiter(test.NoopLifecycle{}, limiter.NewKeyMap(), test.NewLimiterConfig("user-agent", "1s", 0))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, server.Close(t.Context())) })

	routePolicy := http.NewRoutePolicy()
	first, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/hello", http.NoBody)
	require.NoError(t, err)
	httplimiter.NewHandler(routePolicy, server).ServeHTTP(httptest.NewRecorder(), first, func(http.ResponseWriter, *http.Request) {})

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/hello", http.NoBody)
	require.NoError(t, err)
	res := httptest.NewRecorder()

	httplimiter.NewHandler(routePolicy, server).ServeHTTP(res, req, func(http.ResponseWriter, *http.Request) {
		t.Fatal("unexpected next call")
	})

	require.Equal(t, http.StatusTooManyRequests, res.Code)
	require.Nil(t, meta.Limiter(req.Context()))
}

func TestHandlerBypassesOperationPathWithNoLimiterOnContext(t *testing.T) {
	server, err := httplimiter.NewServerLimiter(test.NoopLifecycle{}, limiter.NewKeyMap(), test.NewLimiterConfig("user-agent", "1s", 0))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, server.Close(t.Context())) })

	routePolicy := http.NewRoutePolicy()
	routePolicy.Operation("GET /healthz")

	handler := httplimiter.NewHandler(routePolicy, server)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/healthz", http.NoBody)
	require.NoError(t, err)
	res := httptest.NewRecorder()

	var called bool
	handler.ServeHTTP(res, req, func(_ http.ResponseWriter, req *http.Request) {
		called = true
		require.Nil(t, meta.Limiter(req.Context()))
	})

	require.True(t, called)
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
