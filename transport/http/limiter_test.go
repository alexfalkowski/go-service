package http_test

import (
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/stretchr/testify/require"
)

func TestUnlimited(t *testing.T) {
	cfg := test.NewLimiterConfig("user-agent", "1s", 100)
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldClientLimiter(cfg),
		test.WithWorldServerLimiter(cfg),
		test.WithWorldHTTP(),
		test.WithWorldHello(),
	)

	url := world.PathServerURL("http", "hello")

	_, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
	require.NoError(t, err)

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "hello!", body)
}

func TestServerLimiter(t *testing.T) {
	for _, f := range []string{"user-agent", "ip", "service-method", "transport-service-method"} {
		t.Run(f, func(t *testing.T) {
			world := test.NewStartedWorld(t,
				test.WithWorldTelemetry("otlp"),
				test.WithWorldServerLimiter(test.NewLimiterConfig(f, "1s", 0)),
				test.WithWorldHTTP(),
				test.WithWorldHello(),
			)

			url := world.PathServerURL("http", "hello")

			_, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
			require.NoError(t, err)

			res, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
			require.NoError(t, err)
			require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
			require.Equal(t, `"default";r=0;t=1`, res.Header.Get("Ratelimit"))
			require.Equal(t, `"default";q=1;w=1`, res.Header.Get("Ratelimit-Policy"))
		})
	}
}

func TestServerLimiterSetsRetryAfter(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1h", 1)),
		test.WithWorldHTTP(),
		test.WithWorldHello(),
	)

	url := world.PathServerURL("http", "hello")

	res, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, _, err = world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	require.Contains(t, res.Header.Get("Ratelimit"), `"default";r=0;t=`)
	require.Equal(t, `"default";q=1;w=3600`, res.Header.Get("Ratelimit-Policy"))
	retryAfter, err := strconv.ParseUint(res.Header.Get("Retry-After"), 10, 64)
	require.NoError(t, err)
	require.Positive(t, retryAfter)
	require.LessOrEqual(t, retryAfter, uint64(3600))
}

func TestServerLimiterDoesNotBypassApplicationMetricsPath(t *testing.T) {
	world := test.NewWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 0)),
		test.WithWorldHTTP(),
	)
	world.Handle("GET /admin/metrics", http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		_, _ = res.Write([]byte("secret"))
	}))
	world.Start()

	url := world.PathServerURL("http", "admin/metrics")
	_, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
	require.NoError(t, err)

	res, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	require.Equal(t, `"default";r=0;t=1`, res.Header.Get("Ratelimit"))
	require.Equal(t, `"default";q=1;w=1`, res.Header.Get("Ratelimit-Policy"))
}

func TestClientLimiter(t *testing.T) {
	for _, f := range []string{"user-agent", "ip", "service-method", "transport-service-method"} {
		t.Run(f, func(t *testing.T) {
			world := test.NewStartedWorld(t,
				test.WithWorldTelemetry("otlp"),
				test.WithWorldClientLimiter(test.NewLimiterConfig(f, "1s", 0)),
				test.WithWorldHTTP(),
				test.WithWorldHello(),
			)

			url := world.PathServerURL("http", "hello")

			_, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
			require.NoError(t, err)

			_, _, err = world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
			require.Error(t, err)
			require.Equal(t, http.StatusTooManyRequests, status.Code(err))
		})
	}
}

func TestServerLimiterUsesVerifiedUserID(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldServerLimiter(test.NewLimiterConfig("user-id", "1s", 0)),
		test.WithWorldToken(&test.SequenceGenerator{}, test.AcceptingVerifier{}),
		test.WithWorldHTTP(),
		test.WithWorldHello(),
	)

	url := world.PathServerURL("http", "hello")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "hello!", body)

	res, body, err = world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	require.Equal(t, http.StatusText(http.StatusTooManyRequests), body)
}

func TestServerClosedLimiter(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)),
		test.WithWorldHTTP(),
		test.WithWorldHello(),
	)

	err := world.Server.HTTPLimiter.Close(t.Context())
	require.NoError(t, err)

	url := world.PathServerURL("http", "hello")

	res, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestClientClosedLimiter(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldClientLimiter(test.NewLimiterConfig("user-agent", "1s", 100)),
		test.WithWorldHTTP(),
		test.WithWorldHello(),
	)

	url := world.PathServerURL("http", "hello")

	err := world.Client.HTTPLimiter.Close(t.Context())
	require.NoError(t, err)

	_, _, err = world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
	require.Error(t, err)
	require.Equal(t, http.StatusInternalServerError, status.Code(err))
}
