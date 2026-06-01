package http_test

import (
	"net/http/httptest"
	"testing"

	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/time"
	transporthttp "github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/breaker"
	httplimiter "github.com/alexfalkowski/go-service/v2/transport/http/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/http/retry"
	"github.com/alexfalkowski/go-service/v2/transport/limiter"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	transporthttp.Register(test.FS)

	_, err := transporthttp.NewClient(transporthttp.WithClientTLS(&tls.Config{Cert: "bob", Key: "bob"}))
	require.Error(t, err)

	_, err = transporthttp.NewClient(transporthttp.WithClientTLS(&tls.Config{}))
	require.NoError(t, err)
}

func TestClientRoundTripperBypassesTLSConfig(t *testing.T) {
	transporthttp.Register(test.FS)

	called := false
	client, err := transporthttp.NewClient(
		transporthttp.WithClientRoundTripper(test.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
			called = true
			return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: http.Header{}}, nil
		})),
		transporthttp.WithClientTLS(&tls.Config{Cert: "bob", Key: "bob"}),
	)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", http.NoBody)
	require.NoError(t, err)

	res, err := client.Do(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.True(t, called)
}

func TestClientNegativeTimeoutUsesDefault(t *testing.T) {
	client, err := transporthttp.NewClient(transporthttp.WithClientTimeout(-time.Second))
	require.NoError(t, err)
	require.Equal(t, time.DefaultTimeout.Duration(), client.Timeout)
}

func TestClientWithTokenDoesNotFollowCrossOriginRedirect(t *testing.T) {
	client, err := transporthttp.NewClient(
		transporthttp.WithClientTokenGenerator(env.UserID("service-user"), test.NewGenerator("secret", nil)),
	)
	require.NoError(t, err)

	prev, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com/start", http.NoBody)
	require.NoError(t, err)

	next, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://other.example.com/target", http.NoBody)
	require.NoError(t, err)

	require.ErrorIs(t, client.CheckRedirect(next, []*http.Request{prev}), http.ErrUseLastResponse)
}

func TestRoundTripperWithTokenDoesNotSendAuthorizationToCrossOriginRedirect(t *testing.T) {
	var attackerAuthorization string
	attacker := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		attackerAuthorization = req.Header.Get("Authorization")
	}))
	t.Cleanup(attacker.Close)

	var trustedAuthorization string
	trusted := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		trustedAuthorization = req.Header.Get("Authorization")
		res.Header().Set("Location", attacker.URL+"/target")
		res.WriteHeader(http.StatusTemporaryRedirect)
	}))
	t.Cleanup(trusted.Close)

	rt, err := transporthttp.NewRoundTripper(
		transporthttp.WithClientTokenGenerator(env.UserID("service-user"), test.NewGenerator("secret", nil)),
	)
	require.NoError(t, err)

	client := &http.Client{Transport: rt}
	res, err := client.Get(trusted.URL + "/start")
	require.ErrorIs(t, err, http.ErrUseLastResponse)
	require.Nil(t, res)
	require.Equal(t, "Bearer secret", trustedAuthorization)
	require.Empty(t, attackerAuthorization)
}

func TestRoundTripperGeneratesTokenPerRetryAttempt(t *testing.T) {
	base := &test.AuthRoundTripper{Codes: []int{http.StatusTooManyRequests, http.StatusOK}}
	rt, err := transporthttp.NewRoundTripper(
		transporthttp.WithClientRoundTripper(base),
		transporthttp.WithClientRetry(&retry.Config{
			Attempts: 2,
			Timeout:  time.Second,
			Backoff:  time.Millisecond,
		}),
		transporthttp.WithClientTokenGenerator(env.UserID("service-user"), &test.SequenceGenerator{}),
	)
	require.NoError(t, err)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/hello", http.NoBody)
	require.NoError(t, err)

	res, err := rt.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, []string{"Bearer token-1", "Bearer token-2"}, base.AuthValues)
	require.Equal(t, []int{1, 1}, base.AuthCounts)
}

func TestRoundTripperLimiterDenialDoesNotOpenBreaker(t *testing.T) {
	clientLimiter, err := httplimiter.NewClientLimiter(test.NoopLifecycle{}, limiter.NewKeyMap(), test.NewLimiterConfig("user-agent", "1s", 1))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, clientLimiter.Close(t.Context())) })

	calls := 0
	rt, err := transporthttp.NewRoundTripper(
		transporthttp.WithClientRoundTripper(test.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
			calls++
			return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: http.Header{}}, nil
		})),
		transporthttp.WithClientLimiter(clientLimiter),
		transporthttp.WithClientBreaker(breaker.WithSettings(breaker.Settings{
			ReadyToTrip: func(counts breaker.Counts) bool {
				return counts.ConsecutiveFailures >= 1
			},
		})),
	)
	require.NoError(t, err)

	t.Run("allows first request", func(t *testing.T) {
		res, err := rt.RoundTrip(newUserAgentRequest(t, "first-agent"))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.NoError(t, res.Body.Close())
		require.Equal(t, 1, calls)
	})

	t.Run("denies exhausted key", func(t *testing.T) {
		res, err := rt.RoundTrip(newUserAgentRequest(t, "first-agent"))
		require.Nil(t, res)
		require.Error(t, err)
		require.Equal(t, http.StatusTooManyRequests, status.Code(err))
		require.Equal(t, 1, calls)
	})

	t.Run("denial does not open breaker", func(t *testing.T) {
		res, err := rt.RoundTrip(newUserAgentRequest(t, "first-agent"))
		require.Nil(t, res)
		require.Error(t, err)
		require.Equal(t, http.StatusTooManyRequests, status.Code(err))
		require.Equal(t, 1, calls)
	})

	t.Run("allows different key", func(t *testing.T) {
		res, err := rt.RoundTrip(newUserAgentRequest(t, "second-agent"))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.NoError(t, res.Body.Close())
		require.Equal(t, 2, calls)
	})
}

func newUserAgentRequest(t *testing.T, userAgent string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/hello", http.NoBody)
	require.NoError(t, err)
	req.Header.Set("User-Agent", userAgent)

	return req
}
