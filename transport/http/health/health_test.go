package health_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/mime"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	checks := []string{"healthz", "livez", "readyz"}

	for _, check := range checks {
		t.Run(check, func(t *testing.T) {
			world := test.NewStartedWorld(t,
				test.WithWorldTelemetry("otlp"),
				test.WithWorldHTTPHealth(test.Name.String(), test.StatusURL("200"), test.HealthObserve(check, "http")),
			)

			ctx := t.Context()
			ctx = meta.WithRequestID(ctx, meta.String("test-id"))
			ctx = meta.WithUserAgent(ctx, meta.String("test-user-agent"))

			header := http.Header{}
			url := world.NamedServerURL("http", check)

			res, body, err := world.ResponseWithBody(ctx, url, http.MethodGet, header, http.NoBody)
			require.NoError(t, err)

			require.Equal(t, http.StatusOK, res.StatusCode)
			require.Equal(t, "SERVING", body)
		})
	}
}

func TestReadinessNoop(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldHTTPHealth(test.Name.String(), test.StatusURL("500"), test.HealthObserve("readyz", "noop")),
	)

	header := http.Header{}
	header.Add("Request-Id", "test-id")
	header.Add("User-Agent", "test-user-agent")

	url := world.NamedServerURL("http", "readyz")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "SERVING", body)
	require.Equal(t, mime.TextMediaType, res.Header.Get(content.TypeKey))
}

func TestInvalidHealth(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldHTTPHealth(test.Name.String(), test.StatusURL("500"), test.HealthObserve("healthz", "http")),
	)

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
			world := test.NewStartedWorld(t,
				test.WithWorldTelemetry("otlp"),
				test.WithWorldHTTPHealth(test.Name.String(), test.StatusURL("200")),
			)

			ctx := t.Context()
			ctx = meta.WithRequestID(ctx, meta.String("test-id"))
			ctx = meta.WithUserAgent(ctx, meta.String("test-user-agent"))

			header := http.Header{}
			header.Set(content.TypeKey, mime.JSONMediaType)

			url := world.NamedServerURL("http", check)

			res, err := world.ResponseWithNoBody(ctx, url, http.MethodGet, header)
			require.NoError(t, err)

			require.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
		})
	}
}
