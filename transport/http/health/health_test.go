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
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	netserver "github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/http/health"
	"github.com/stretchr/testify/require"
)

func TestHealthRegistersHealthLivenessAndReadinessEndpoints(t *testing.T) {
	checks := []struct {
		name string
		path string
	}{
		{name: "serves health endpoint", path: "healthz"},
		{name: "serves liveness endpoint", path: "livez"},
		{name: "serves readiness endpoint", path: "readyz"},
	}

	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			world := test.NewStartedWorld(t,
				test.WithWorldTelemetry("otlp"),
				test.WithWorldHTTPHealth(test.Name.String(), test.StatusURL("200"), test.HealthObserve(check.path, "http")),
			)

			ctx := meta.WithAttributes(t.Context(),
				meta.WithRequestID(meta.String("test-id")),
				meta.WithUserAgent(meta.String("test-user-agent")),
			)

			header := http.Header{}
			url := world.NamedServerURL("http", check.path)

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
	require.Equal(t, "text/plain; charset=utf-8", res.Header.Get(http.ContentTypeKey))
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
	router := http.NewRouter(mux, http.NewRoutePolicy())
	health.Register(health.RegisterParams{
		Name:   test.Name,
		Server: srv,
		Router: router,
		Drain:  drain,
	})

	require.Eventually(t, func() bool {
		return healthStatus(mux, "readyz") == http.StatusOK
	}, time.Second.Duration(), time.Millisecond.Duration())

	drain.Start()

	require.Equal(t, http.StatusServiceUnavailable, healthStatus(mux, "readyz"))
	require.Equal(t, http.StatusOK, healthStatus(mux, "livez"))
}

func TestRejectsInvalidHealth(t *testing.T) {
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
	require.Equal(t, "text/error; charset=utf-8", res.Header.Get(http.ContentTypeKey))
}

func TestHealthEndpointReportsMissingCheck(t *testing.T) {
	checks := []struct {
		name string
		path string
	}{
		{name: "rejects unobserved health endpoint", path: "healthz"},
		{name: "rejects unobserved liveness endpoint", path: "livez"},
		{name: "rejects unobserved readiness endpoint", path: "readyz"},
	}

	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			world := test.NewStartedWorld(t,
				test.WithWorldTelemetry("otlp"),
				test.WithWorldHTTPHealth(test.Name.String(), test.StatusURL("200")),
			)

			ctx := meta.WithAttributes(t.Context(),
				meta.WithRequestID(meta.String("test-id")),
				meta.WithUserAgent(meta.String("test-user-agent")),
			)

			header := http.Header{}
			header.Set(http.ContentTypeKey, media.JSON)

			url := world.NamedServerURL("http", check.path)

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
