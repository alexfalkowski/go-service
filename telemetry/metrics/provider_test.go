package metrics_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/time"
	sync "github.com/alexfalkowski/go-sync"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestIsEnabled(t *testing.T) {
	t.Cleanup(func() {
		metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: fxtest.NewLifecycle(t)})
	})

	metrics.SetMeterProvider(metrics.NewNoopMeterProvider())
	require.False(t, metrics.IsEnabled())

	metrics.SetMeterProvider(nil)
	require.False(t, metrics.IsEnabled())

	manualProvider := metrics.NewMeterProvider(metrics.MeterProviderParams{
		Lifecycle:   fxtest.NewLifecycle(t),
		Config:      &metrics.Config{},
		Reader:      metrics.NewManualReader(),
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	})
	require.True(t, metrics.IsEnabled())

	metrics.SetMeterProvider(metrics.NewNoopMeterProvider())
	require.False(t, metrics.IsEnabled())

	metrics.SetMeterProvider(manualProvider)
	require.True(t, metrics.IsEnabled())

	metrics.SetMeterProvider(metrics.NewNoopMeterProvider())
	require.False(t, metrics.IsEnabled())

	metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: fxtest.NewLifecycle(t)})
	require.False(t, metrics.IsEnabled())

	provider := metrics.NewMeterProvider(metrics.MeterProviderParams{
		Lifecycle:   fxtest.NewLifecycle(t),
		Config:      &metrics.Config{},
		Reader:      metrics.NewManualReader(),
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	})
	require.NotNil(t, provider)
	require.True(t, metrics.IsEnabled())

	metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: fxtest.NewLifecycle(t)})
	require.False(t, metrics.IsEnabled())
}

func TestMeterProviderStopResetsGlobalState(t *testing.T) {
	t.Cleanup(func() {
		metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: fxtest.NewLifecycle(t)})
	})

	lc := fxtest.NewLifecycle(t)
	provider := metrics.NewMeterProvider(metrics.MeterProviderParams{
		Lifecycle:   lc,
		Config:      &metrics.Config{},
		Reader:      metrics.NewManualReader(),
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	})
	require.True(t, metrics.IsEnabled())
	require.Equal(t, provider, metrics.GetMeterProvider())

	lc.RequireStart()
	require.NoError(t, lc.Stop(t.Context()))

	require.False(t, metrics.IsEnabled())
	require.NotEqual(t, provider, metrics.GetMeterProvider())
}

func TestMeterProviderAttributes(t *testing.T) {
	t.Cleanup(func() {
		metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: fxtest.NewLifecycle(t)})
	})

	reader := metrics.NewManualReader()
	provider := metrics.NewMeterProvider(metrics.MeterProviderParams{
		Lifecycle: fxtest.NewLifecycle(t),
		Config:    &metrics.Config{},
		Reader:    reader,
		Attributes: attributes.Map{
			"k8s.namespace.name":                      "payments",
			string(attributes.ServiceName("").Key):    "configured",
			string(attributes.ServiceVersion("").Key): "configured",
		},
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	})
	counter, err := provider.Meter(test.Name.String()).Int64Counter("requests")
	require.NoError(t, err)
	counter.Add(t.Context(), 1)

	rm := metrics.ResourceMetrics{}
	require.NoError(t, reader.Collect(t.Context(), &rm))
	attrs := make(map[attributes.Key]string)
	for _, attr := range rm.Resource.Attributes() {
		attrs[attr.Key] = attr.Value.AsString()
	}

	require.Equal(t, test.ID.String(), attrs[attributes.HostID("").Key])
	require.Equal(t, test.Name.String(), attrs[attributes.ServiceName("").Key])
	require.Equal(t, test.Version.String(), attrs[attributes.ServiceVersion("").Key])
	require.Equal(t, "development", attrs[attributes.DeploymentEnvironmentName("").Key])
	require.Equal(t, "payments", attrs[attributes.Key("k8s.namespace.name")])
}

func TestMeterProviderViews(t *testing.T) {
	t.Cleanup(func() {
		metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: fxtest.NewLifecycle(t)})
	})

	boundaries := []float64{0.1, 0.5, 1}
	reader := metrics.NewManualReader()
	provider := metrics.NewMeterProvider(metrics.MeterProviderParams{
		Lifecycle:   fxtest.NewLifecycle(t),
		Config:      &metrics.Config{Views: map[string][]float64{"test.duration": boundaries}},
		Reader:      reader,
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	})

	histogram, err := provider.Meter(test.Name.String()).Float64Histogram("test.duration")
	require.NoError(t, err)
	histogram.Record(t.Context(), 0.3)

	rm := metrics.ResourceMetrics{}
	require.NoError(t, reader.Collect(t.Context(), &rm))

	require.Equal(t, boundaries, collectHistogramBounds(t, rm, "test.duration"))
}

func TestOTLPReaderUsesConfiguredInterval(t *testing.T) {
	t.Cleanup(func() {
		metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: fxtest.NewLifecycle(t)})
	})

	var exports sync.Int64
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		exports.Add(1)
		res.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	lc := fxtest.NewLifecycle(t)
	cfg := &metrics.Config{
		Kind:     "otlp",
		URL:      server.URL,
		Interval: 20 * time.Millisecond,
		Timeout:  time.Second,
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
	counter, err := provider.Meter(test.Name.String()).Int64Counter("requests")
	require.NoError(t, err)
	counter.Add(t.Context(), 1)

	require.Eventually(t, func() bool {
		return exports.Load() > 0
	}, (2 * time.Second).Duration(), (20 * time.Millisecond).Duration())
}

func collectHistogramBounds(t *testing.T, rm metrics.ResourceMetrics, name string) []float64 {
	t.Helper()

	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name != name {
				continue
			}

			histogram, ok := m.Data.(metrics.Histogram[float64])
			require.True(t, ok)
			require.NotEmpty(t, histogram.DataPoints)

			return histogram.DataPoints[0].Bounds
		}
	}

	t.Fatalf("histogram %q not collected", name)
	return nil
}
