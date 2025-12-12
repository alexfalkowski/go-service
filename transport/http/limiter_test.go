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
	world := test.NewWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldClientLimiter(cfg),
		test.WithWorldServerLimiter(cfg),
		test.WithWorldHTTP(),
	)
	world.Register()
	world.HandleHello()
	world.RequireStart()

	url := world.PathServerURL("http", "hello")

	_, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
	require.NoError(t, err)

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "hello!", body)

	world.RequireStop()
}

func TestServerLimiter(t *testing.T) {
	for _, f := range []string{"user-agent", "ip"} {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig(f, "1s", 0)), test.WithWorldHTTP())
		world.Register()
		world.HandleHello()
		world.RequireStart()

		url := world.PathServerURL("http", "hello")

		_, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
		require.NoError(t, err)

		res, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
		require.NotEmpty(t, res.Header.Get("Ratelimit"))

		world.RequireStop()
	}
}

func TestClientLimiter(t *testing.T) {
	for _, f := range []string{"user-agent", "ip"} {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldClientLimiter(test.NewLimiterConfig(f, "1s", 0)), test.WithWorldHTTP())
		world.Register()
		world.HandleHello()
		world.RequireStart()

		url := world.PathServerURL("http", "hello")

		_, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
		require.NoError(t, err)

		_, _, err = world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
		require.Error(t, err)
		require.Equal(t, http.StatusTooManyRequests, status.Code(err))

		world.RequireStop()
	}
}

func TestServerClosedLimiter(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
	world.Register()
	world.HandleHello()
	world.RequireStart()

	err := world.Server.HTTPLimiter.Close(t.Context())
	require.NoError(t, err)

	url := world.PathServerURL("http", "hello")

	res, _, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)

	world.RequireStop()
}

func TestClientClosedLimiter(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldClientLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
	world.Register()
	world.HandleHello()
	world.RequireStart()

	url := world.PathServerURL("http", "hello")

	err := world.Client.HTTPLimiter.Close(t.Context())
	require.NoError(t, err)

	_, _, err = world.ResponseWithBody(t.Context(), url, http.MethodGet, http.Header{}, http.NoBody)
	require.Error(t, err)
	require.Equal(t, http.StatusInternalServerError, status.Code(err))

	world.RequireStop()
}
