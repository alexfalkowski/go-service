package health_test

import (
	"encoding/json"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/mime"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	checks := []string{"healthz", "livez", "readyz"}

	for _, check := range checks {
		t.Run(check, func(t *testing.T) {
			world := test.NewStartedHTTPWorld(t, func(world *test.World) {
				server := world.HealthServer(test.Name.String(), test.StatusURL("200"))
				err := server.Observe(test.Name.String(), check, "http")
				require.NoError(t, err)
				test.RegisterHealth(server)
			}, test.WithWorldTelemetry("otlp"))

			header := http.Header{}
			header.Set(content.TypeKey, mime.JSONMediaType)

			url := world.NamedServerURL("http", check)

			res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
			require.NoError(t, err)

			require.Equal(t, http.StatusOK, res.StatusCode)
			require.Equal(t, mime.JSONMediaType, res.Header.Get(content.TypeKey))

			var resp healthResponse
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, "SERVING", resp.Status)
			require.Empty(t, resp.Meta)
		})
	}
}

func TestReadinessNoop(t *testing.T) {
	world := test.NewStartedHTTPWorld(t, func(world *test.World) {
		server := world.HealthServer(test.Name.String(), test.StatusURL("500"))
		err := server.Observe(test.Name.String(), "readyz", "noop")
		require.NoError(t, err)
		test.RegisterHealth(server)
	}, test.WithWorldTelemetry("otlp"))

	header := http.Header{}
	header.Add("Request-Id", "test-id")
	header.Add("User-Agent", "test-user-agent")
	header.Set(content.TypeKey, mime.JSONMediaType)

	url := world.NamedServerURL("http", "readyz")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, mime.JSONMediaType, res.Header.Get(content.TypeKey))

	var resp healthResponse
	require.NoError(t, json.Unmarshal([]byte(body), &resp))
	require.Equal(t, "SERVING", resp.Status)
	require.Empty(t, resp.Meta)
}

func TestInvalidHealth(t *testing.T) {
	world := test.NewStartedHTTPWorld(t, func(world *test.World) {
		server := world.HealthServer(test.Name.String(), test.StatusURL("500"))
		err := server.Observe(test.Name.String(), "healthz", "http")
		require.NoError(t, err)
		test.RegisterHealth(server)
	}, test.WithWorldTelemetry("otlp"))

	header := http.Header{}
	url := world.NamedServerURL("http", "healthz")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)

	require.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
	require.Equal(t, "http: http checker: invalid status code", body)
	require.Equal(t, mime.ErrorMediaType, res.Header.Get(content.TypeKey))
}

func TestMissingHealth(t *testing.T) {
	checks := []string{"healthz", "livez", "readyz"}

	for _, check := range checks {
		t.Run(check, func(t *testing.T) {
			world := test.NewStartedHTTPWorld(t, func(world *test.World) {
				server := world.HealthServer(test.Name.String(), test.StatusURL("200"))
				test.RegisterHealth(server)
			}, test.WithWorldTelemetry("otlp"))

			header := http.Header{}
			header.Set(content.TypeKey, mime.JSONMediaType)

			url := world.NamedServerURL("http", check)

			res, err := world.ResponseWithNoBody(t.Context(), url, http.MethodGet, header)
			require.NoError(t, err)

			require.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
		})
	}
}

type healthResponse struct {
	Meta   map[string]string `json:"meta"`
	Status string            `json:"status"`
}
