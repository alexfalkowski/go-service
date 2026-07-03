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

func TestConfigGetProtocol(t *testing.T) {
	require.Equal(t, otlp.ProtocolHTTP, (*metrics.Config)(nil).GetProtocol())
	require.Equal(t, otlp.ProtocolHTTP, (&metrics.Config{}).GetProtocol())
	require.Equal(t, otlp.ProtocolGRPC, (&metrics.Config{Protocol: otlp.ProtocolGRPC}).GetProtocol())
}

func TestReaderShutdownIgnoresAlreadyShutdownReader(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	reader, err := metrics.NewReader(lc, test.Name, &metrics.Config{Kind: "prometheus"})
	require.NoError(t, err)

	lc.RequireStart()
	require.NoError(t, reader.Shutdown(t.Context()))

	require.NoError(t, lc.Stop(t.Context()))
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

func TestOTLPGRPCReader(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	cfg := &metrics.Config{
		Kind:     "otlp",
		Protocol: "grpc",
		URL:      "localhost:4317",
	}

	reader, err := metrics.NewReader(lc, test.Name, cfg)
	require.NoError(t, err)
	require.NotNil(t, reader)
	require.NoError(t, reader.Shutdown(t.Context()))
}

func TestInvalidOTLPGRPCEndpoint(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	cfg := &metrics.Config{
		Kind:     "otlp",
		Protocol: "grpc",
		URL:      "collector.example.com:4317",
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
