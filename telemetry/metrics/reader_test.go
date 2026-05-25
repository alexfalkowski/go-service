package metrics_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/header"
	"github.com/alexfalkowski/go-service/v2/telemetry/internal/otlp"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestInvalidReader(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	cfg := &metrics.Config{Kind: "wrong"}

	_, err := metrics.NewReader(lc, test.Name, cfg)
	require.Error(t, err)
}

func TestInvalidOTLPEndpoint(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	cfg := &metrics.Config{
		Kind: "otlp",
		URL:  "http://collector.example.com/v1/metrics",
		Headers: header.Map{
			"Authorization": "Bearer token",
		},
	}

	_, err := metrics.NewReader(lc, test.Name, cfg)
	require.ErrorIs(t, err, otlp.ErrInsecureEndpoint)
}

func TestMissingOTLPEndpointIgnoresEnv(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT", "http://collector.example.com/v1/metrics")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "https://collector.example.com")

	lc := fxtest.NewLifecycle(t)
	cfg := &metrics.Config{
		Kind: "otlp",
		Headers: header.Map{
			"Authorization": "Bearer token",
		},
	}

	_, err := metrics.NewReader(lc, test.Name, cfg)
	require.ErrorIs(t, err, otlp.ErrMissingEndpoint)
}
