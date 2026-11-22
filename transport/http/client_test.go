package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v3/config/client"
	tls "github.com/alexfalkowski/go-service/v3/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v3/env"
	"github.com/alexfalkowski/go-service/v3/internal/test"
	"github.com/alexfalkowski/go-service/v3/net/http"
	"github.com/alexfalkowski/go-service/v3/net/http/status"
	"github.com/alexfalkowski/go-service/v3/net/url"
	"github.com/alexfalkowski/go-service/v3/time"
	transporthttp "github.com/alexfalkowski/go-service/v3/transport/http"
	"github.com/alexfalkowski/go-service/v3/transport/http/breaker"
	httplimiter "github.com/alexfalkowski/go-service/v3/transport/http/limiter"
	httpretry "github.com/alexfalkowski/go-service/v3/transport/http/retry"
	"github.com/alexfalkowski/go-service/v3/transport/limiter"
	"github.com/alexfalkowski/go-signal/v2/signaltest"
	"github.com/stretchr/testify/require"
)

func TestNewClientFromConfigReturnsEndpointAndAppliesOptions(t *testing.T) {
	cfg := &client.Config{Endpoint: "https://example.com/api/", Timeout: time.Second}

	conn, endpoint, err := transporthttp.NewClientFromConfig(cfg)
	require.NoError(t, err)
	require.Equal(t, cfg.Endpoint, endpoint.String())
	require.Equal(t, cfg.Timeout.Duration(), conn.Timeout)

	conn, _, err = transporthttp.NewClientFromConfig(cfg, transporthttp.WithClientTimeout(time.Minute))
	require.NoError(t, err)
	require.Equal(t, time.Minute.Duration(), conn.Timeout)
}

func TestNewClientFromConfigRejectsInvalidEndpoint(t *testing.T) {
	for _, tt := range []struct {
		cfg  *client.Config
		err  error
		name string
	}{
		{name: "missing config", err: url.ErrInvalidSecureURL},
		{name: "missing endpoint", cfg: &client.Config{}, err: url.ErrInvalidSecureURL},
		{name: "non-loopback HTTP", cfg: &client.Config{Endpoint: "http://example.com"}, err: url.ErrInsecureURL},
	} {
		t.Run(tt.name, func(t *testing.T) {
			conn, endpoint, err := transporthttp.NewClientFromConfig(tt.cfg)
			require.ErrorIs(t, err, tt.err)
			require.Nil(t, conn)
			require.Nil(t, endpoint)
		})
	}
}

func TestNewClientFromConfigUsesBreakerThreshold(t *testing.T) {
	cfg := &client.Config{Endpoint: "https://example.com", Breaker: test.NewBreaker(1)}
	base := &test.StatusSequenceRoundTripper{Codes: []int{http.StatusServiceUnavailable, http.StatusOK}}
	conn, endpoint, err := transporthttp.NewClientFromConfig(cfg,
		transporthttp.WithClientRoundTripper(base),
	)
	require.NoError(t, err)

	res, err := conn.Get(endpoint.String())
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusServiceUnavailable, res.StatusCode)

	res, err = conn.Get(endpoint.String())
	require.ErrorIs(t, err, breaker.ErrOpenState)
	require.Nil(t, res)
	require.Equal(t, 1, base.Calls, "the open breaker rejects the second request locally")
}

func TestNewClientWithLifecycle(t *testing.T) {
	for _, tt := range []struct {
		name       string
		closesIdle bool
	}{
		{name: "with idle cleanup", closesIdle: true},
		{name: "without idle cleanup"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			lc := signaltest.NewLifecycle(t)
			t.Cleanup(lc.RequireStop)
			base := test.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: http.Header{}}, nil
			})
			pool := &idleClosingTransport{RoundTripperFunc: base}
			var rt http.RoundTripper = base
			if tt.closesIdle {
				rt = pool
			}

			conn, err := transporthttp.NewClient(
				transporthttp.WithClientRoundTripper(rt),
				transporthttp.WithClientLifecycle(lc),
			)
			require.NoError(t, err)
			lc.RequireStart()

			res, err := conn.Get("http://localhost")
			require.NoError(t, err)
			require.NoError(t, res.Body.Close())
			require.Equal(t, http.StatusOK, res.StatusCode)
			require.False(t, pool.closed, "the pool stays open while the application is running")

			lc.RequireStop()
			require.Equal(t, tt.closesIdle, pool.closed, "shutdown closes the pool when the transport supports idle cleanup")
		})
	}
}

func TestNewClientValidatesTLSConfig(t *testing.T) {
	transporthttp.Register(test.FS)

	_, err := transporthttp.NewClient(transporthttp.WithClientTLS(&tls.Config{Cert: "bob", Key: "bob"}))
	require.Error(t, err)

	_, err = transporthttp.NewClient(transporthttp.WithClientTLS(&tls.Config{}))
	require.NoError(t, err)

	cfg := &client.Config{Endpoint: "https://example.com", TLS: &tls.Config{Cert: "bob", Key: "bob"}}
	conn, endpoint, err := transporthttp.NewClientFromConfig(cfg)
	require.Error(t, err)
	require.Nil(t, conn)
	require.Nil(t, endpoint)
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
		attackerAuthorization = req.Header.Get(http.AuthorizationKey)
	}))
	t.Cleanup(attacker.Close)

	var trustedAuthorization string
	trusted := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		trustedAuthorization = req.Header.Get(http.AuthorizationKey)
		res.Header().Set(http.LocationKey, attacker.URL+"/target")
		res.WriteHeader(http.StatusTemporaryRedirect)
	}))
	t.Cleanup(trusted.Close)

	rt, err := transporthttp.NewRoundTripper(
		transporthttp.WithClientTokenGenerator(env.UserID("service-user"), test.NewGenerator("secret", nil)),
	)
	require.NoError(t, err)

	client := &http.Client{Transport: rt}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, trusted.URL+"/start", http.NoBody)
	require.NoError(t, err)

	res, err := client.Do(req)
	require.ErrorIs(t, err, http.ErrUseLastResponse)
	require.Nil(t, res)
	require.Equal(t, "Bearer secret", trustedAuthorization)
	require.Empty(t, attackerAuthorization)
}

func TestRoundTripperRetriesThroughClientLimiter(t *testing.T) {
	clientLimiter, err := httplimiter.NewClientLimiter(
		test.NoopLifecycle{},
		limiter.NewKeyMap(),
		test.NewLimiterConfig("user-agent", "1m", 1),
	)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, clientLimiter.Close(t.Context())) })

	base := &test.StatusSequenceRoundTripper{Codes: []int{http.StatusServiceUnavailable, http.StatusOK}}
	rt, err := transporthttp.NewRoundTripper(
		transporthttp.WithClientRoundTripper(base),
		transporthttp.WithClientLimiter(clientLimiter),
		transporthttp.WithClientRetry(httpretry.NewConfig(test.FastRetryConfig)),
		transporthttp.WithClientUserAgent(test.UserAgent),
	)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/hello", http.NoBody)
	require.NoError(t, err)

	res, err := rt.RoundTrip(req)
	require.Nil(t, res)
	require.Error(t, err)
	require.True(t, status.IsLocalError(err))
	require.Equal(t, http.StatusTooManyRequests, status.Code(err))
	require.Equal(t, 1, base.Calls)
}

func TestRoundTripperRetriesThroughClientBreaker(t *testing.T) {
	base := &test.StatusSequenceRoundTripper{Codes: []int{http.StatusServiceUnavailable, http.StatusOK}}
	rt, err := transporthttp.NewRoundTripper(
		transporthttp.WithClientRoundTripper(base),
		transporthttp.WithClientRetry(httpretry.NewConfig(test.FastRetryConfig)),
		transporthttp.WithClientBreaker(breaker.NewConfig(
			test.NewBreaker(1),
			http.StatusServiceUnavailable,
		)),
	)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/hello", http.NoBody)
	require.NoError(t, err)

	res, err := rt.RoundTrip(req)
	require.Nil(t, res)
	require.ErrorIs(t, err, breaker.ErrOpenState)
	require.True(t, status.IsLocalError(err))
	require.Equal(t, http.StatusServiceUnavailable, status.Code(err))
	require.Equal(t, 1, base.Calls)
}

type idleClosingTransport struct {
	test.RoundTripperFunc
	closed bool
}

// CloseIdleConnections records shutdown of the pool.
func (r *idleClosingTransport) CloseIdleConnections() {
	r.closed = true
}
