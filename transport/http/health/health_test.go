package health_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-health/v2/checker"
	healthserver "github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	netserver "github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/http/health"
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

			ctx := meta.WithAttributes(t.Context(),
				meta.WithRequestID(meta.String("test-id")),
				meta.WithUserAgent(meta.String("test-user-agent")),
			)

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
	require.Equal(t, "text/plain; charset=utf-8", res.Header.Get(content.TypeKey))
}

func TestReadinessCache(t *testing.T) {
	world := test.NewStartedWorld(t,
		test.WithWorldTelemetry("otlp"),
		test.WithWorldHTTPHealth(test.Name.String(), test.StatusURL("500"), test.HealthObserve("readyz", "cache")),
	)

	header := http.Header{}
	url := world.NamedServerURL("http", "readyz")

	res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "SERVING", body)
}

func TestReadinessDrains(t *testing.T) {
	srv := healthserver.NewServer()
	reg := healthserver.NewRegistration("noop", time.Millisecond.Duration(), checker.NewNoopChecker())
	srv.Register(test.Name.String(), reg)
	require.NoError(t, srv.Observe(test.Name.String(), "readyz", "noop"))
	require.NoError(t, srv.Observe(test.Name.String(), "livez", "noop"))
	require.NoError(t, srv.Start(t.Context()))
	t.Cleanup(func() {
		require.NoError(t, srv.Stop(context.Background()))
	})

	drain := netserver.NewDrain()
	mux := http.NewServeMux()
	health.Register(health.RegisterParams{
		Name:   test.Name,
		Server: srv,
		Mux:    mux,
		Drain:  drain,
	})

	require.Eventually(t, func() bool {
		return healthStatus(mux, "readyz") == http.StatusOK
	}, time.Second.Duration(), time.Millisecond.Duration())

	drain.Start()

	require.Equal(t, http.StatusServiceUnavailable, healthStatus(mux, "readyz"))
	require.Equal(t, http.StatusOK, healthStatus(mux, "livez"))
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
	require.Equal(t, "http: service unavailable", body)
	require.Equal(t, "text/error; charset=utf-8", res.Header.Get(content.TypeKey))
}

func TestMissingHealth(t *testing.T) {
	checks := []string{"healthz", "livez", "readyz"}

	for _, check := range checks {
		t.Run(check, func(t *testing.T) {
			world := test.NewStartedWorld(t,
				test.WithWorldTelemetry("otlp"),
				test.WithWorldHTTPHealth(test.Name.String(), test.StatusURL("200")),
			)

			ctx := meta.WithAttributes(t.Context(),
				meta.WithRequestID(meta.String("test-id")),
				meta.WithUserAgent(meta.String("test-user-agent")),
			)

			header := http.Header{}
			header.Set(content.TypeKey, media.JSON)

			url := world.NamedServerURL("http", check)

			res, err := world.ResponseWithNoBody(ctx, url, http.MethodGet, header)
			require.NoError(t, err)

			require.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
		})
	}
}

func healthStatus(mux *http.ServeMux, check string) int {
	res := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, http.Pattern(test.Name, "/"+check), http.NoBody)

	mux.ServeHTTP(res, req)

	return res.Code
}
