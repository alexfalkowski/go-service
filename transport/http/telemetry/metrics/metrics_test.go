package metrics_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	transportmetrics "github.com/alexfalkowski/go-service/v2/transport/http/telemetry/metrics"
	"github.com/stretchr/testify/require"
)

func TestRegisterDisabled(t *testing.T) {
	mux := http.NewServeMux()

	transportmetrics.Register(test.Name, nil, mux)

	res := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/test/metrics", http.NoBody)
	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusNotFound, res.Code)
}

func TestRegisterNonPrometheus(t *testing.T) {
	mux := http.NewServeMux()

	transportmetrics.Register(test.Name, &metrics.Config{Kind: "otlp"}, mux)

	res := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/test/metrics", http.NoBody)
	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusNotFound, res.Code)
}
