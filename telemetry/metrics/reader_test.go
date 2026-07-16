package metrics_test

import (
	"testing"

	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/header"
	"github.com/alexfalkowski/go-service/v2/telemetry/internal/otlp"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx/fxtest"
)

func TestInvalidReader(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	cfg := &metrics.Config{Kind: "wrong"}

	_, err := metrics.NewReader(metrics.ReaderParams{Lifecycle: lc, Config: cfg, FS: test.FS, Name: test.Name})
	require.Error(t, err)
}

func TestConfigGetProtocol(t *testing.T) {
	require.Equal(t, otlp.ProtocolHTTP, (*metrics.Config)(nil).GetProtocol())
	require.Equal(t, otlp.ProtocolHTTP, (&metrics.Config{}).GetProtocol())
	require.Equal(t, otlp.ProtocolGRPC, (&metrics.Config{Protocol: otlp.ProtocolGRPC}).GetProtocol())
}

func TestReaderShutdownIgnoresAlreadyShutdownReader(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	reader, err := metrics.NewReader(metrics.ReaderParams{
		Lifecycle: lc,
		Config:    &metrics.Config{Kind: "prometheus"},
		FS:        test.FS,
		Name:      test.Name,
	})
	require.NoError(t, err)

	lc.RequireStart()
	require.NoError(t, reader.Shutdown(t.Context()))

	require.NoError(t, lc.Stop(t.Context()))
}

func TestPrometheusReaderShapesOutput(t *testing.T) {
	t.Cleanup(func() {
		metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: fxtest.NewLifecycle(t)})
	})

	lc := fxtest.NewLifecycle(t)
	cfg := &metrics.Config{
		Kind: "prometheus",
		Prometheus: &metrics.PrometheusConfig{
			WithoutTargetInfo: true,
			WithoutSuffixes:   true,
		},
	}

	reader, err := metrics.NewReader(metrics.ReaderParams{Lifecycle: lc, Config: cfg, FS: test.FS, Name: test.Name})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = reader.Shutdown(t.Context())
	})

	provider := metrics.NewMeterProvider(metrics.MeterProviderParams{
		Lifecycle:   lc,
		Config:      cfg,
		Reader:      reader,
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	})

	histogram, err := provider.Meter(test.Name.String()).Float64Histogram("prometheus_reader_shapes_output_duration", metric.WithUnit("s"))
	require.NoError(t, err)
	histogram.Record(t.Context(), 0.3)

	counter, err := provider.Meter(test.Name.String()).Int64Counter("prometheus_reader_shapes_output_count")
	require.NoError(t, err)
	counter.Add(t.Context(), 1)

	names := gatherFamilyNames(t)

	require.NotContains(t, names, "target_info")
	require.Contains(t, names, "test_prometheus_reader_shapes_output_duration")
	require.NotContains(t, names, "test_prometheus_reader_shapes_output_duration_seconds")
	require.Contains(t, names, "test_prometheus_reader_shapes_output_count")
	require.NotContains(t, names, "test_prometheus_reader_shapes_output_count_total")
}

func TestPrometheusReaderKeepsDefaultOutput(t *testing.T) {
	t.Cleanup(func() {
		metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: fxtest.NewLifecycle(t)})
	})

	lc := fxtest.NewLifecycle(t)
	cfg := &metrics.Config{Kind: "prometheus"}

	reader, err := metrics.NewReader(metrics.ReaderParams{Lifecycle: lc, Config: cfg, FS: test.FS, Name: test.Name})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = reader.Shutdown(t.Context())
	})

	provider := metrics.NewMeterProvider(metrics.MeterProviderParams{
		Lifecycle:   lc,
		Config:      cfg,
		Reader:      reader,
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	})

	histogram, err := provider.Meter(test.Name.String()).Float64Histogram("prometheus_reader_keeps_default_output_duration", metric.WithUnit("s"))
	require.NoError(t, err)
	histogram.Record(t.Context(), 0.3)

	names := gatherFamilyNames(t)

	require.Contains(t, names, "target_info")
	require.Contains(t, names, "test_prometheus_reader_keeps_default_output_duration_seconds")
}

func gatherFamilyNames(t *testing.T) []string {
	t.Helper()

	families, err := prometheus.DefaultGatherer.Gather()
	require.NoError(t, err)

	names := make([]string, 0, len(families))
	for _, family := range families {
		names = append(names, family.GetName())
	}

	return names
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

	_, err := metrics.NewReader(metrics.ReaderParams{Lifecycle: lc, Config: cfg, FS: test.FS, Name: test.Name})
	require.ErrorIs(t, err, otlp.ErrInsecureEndpoint)
}

func TestOTLPGRPCReader(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	cfg := &metrics.Config{
		Kind:     "otlp",
		Protocol: "grpc",
		URL:      "localhost:4317",
	}

	reader, err := metrics.NewReader(metrics.ReaderParams{Lifecycle: lc, Config: cfg, FS: test.FS, Name: test.Name})
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

	_, err := metrics.NewReader(metrics.ReaderParams{Lifecycle: lc, Config: cfg, FS: test.FS, Name: test.Name})
	require.ErrorIs(t, err, otlp.ErrInsecureEndpoint)
}

func TestOTLPGRPCReaderWithTLSHeaders(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	cfg := &metrics.Config{
		Kind:     "otlp",
		Protocol: "grpc",
		URL:      "collector.example.com:4317",
		TLS:      &tls.Config{ServerName: "collector.example.com"},
		Headers: header.Map{
			"Authorization": "Bearer token",
		},
	}

	reader, err := metrics.NewReader(metrics.ReaderParams{Lifecycle: lc, Config: cfg, FS: test.FS, Name: test.Name})
	require.NoError(t, err)
	require.NotNil(t, reader)
	require.NoError(t, reader.Shutdown(t.Context()))
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

	_, err := metrics.NewReader(metrics.ReaderParams{Lifecycle: lc, Config: cfg, FS: test.FS, Name: test.Name})
	require.ErrorIs(t, err, otlp.ErrMissingEndpoint)
}
