package http_test

import (
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
	for _, f := range []string{"user-agent", "ip"} {
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
			require.NotEmpty(t, res.Header.Get("Ratelimit"))
		})
	}
}

func TestServerLimiterDoesNotBypassApplicationMetricsPath(t *testing.T) {
	world := test.NewWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 0)),
		test.WithWorldHTTP(),
	)
	http.HandleFunc(world.ServeMux, "GET /admin/metrics", func(res http.ResponseWriter, _ *http.Request) {
		_, _ = res.Write([]byte("secret"))
	})
	world.Start()

	url := world.PathServerURL("http", "admin/metrics")
	_, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
	require.NoError(t, err)

	res, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	require.NotEmpty(t, res.Header.Get("Ratelimit"))
}

func TestClientLimiter(t *testing.T) {
	for _, f := range []string{"user-agent", "ip"} {
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

func TestClientLimiterUsesGeneratedToken(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldClientLimiter(test.NewLimiterConfig("token", "1s", 0)),
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
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "hello!", body)
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
