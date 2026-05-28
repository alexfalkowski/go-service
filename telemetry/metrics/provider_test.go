package metrics_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.uber.org/fx/fxtest"
)

func TestIsEnabled(t *testing.T) {
	t.Cleanup(func() {
		metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: fxtest.NewLifecycle(t)})
	})

	otel.SetMeterProvider(noop.NewMeterProvider())
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
	require.Equal(t, provider, otel.GetMeterProvider())

	lc.RequireStart()
	require.NoError(t, lc.Stop(t.Context()))

	require.False(t, metrics.IsEnabled())
	require.NotEqual(t, provider, otel.GetMeterProvider())
}

func TestMeterProviderResourceAttributes(t *testing.T) {
	t.Cleanup(func() {
		metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: fxtest.NewLifecycle(t)})
	})

	reader := metrics.NewManualReader()
	provider := metrics.NewMeterProvider(metrics.MeterProviderParams{
		Lifecycle:   fxtest.NewLifecycle(t),
		Config:      &metrics.Config{},
		Reader:      reader,
		ID:          test.ID,
		Name:        test.Name,
		Version:     test.Version,
		Environment: test.Environment,
	})
	counter, err := provider.Meter(test.Name.String()).Int64Counter("requests")
	require.NoError(t, err)
	counter.Add(t.Context(), 1)

	rm := metricdata.ResourceMetrics{}
	require.NoError(t, reader.Collect(t.Context(), &rm))
	attrs := resourceAttributes(rm)

	require.Equal(t, test.ID.String(), attrs[attributes.HostID("").Key])
	require.Equal(t, test.Name.String(), attrs[attributes.ServiceName("").Key])
	require.Equal(t, test.Version.String(), attrs[attributes.ServiceVersion("").Key])
	require.Equal(t, "development", attrs[attributes.DeploymentEnvironmentName("").Key])
}

func resourceAttributes(rm metricdata.ResourceMetrics) map[attribute.Key]string {
	attrs := make(map[attribute.Key]string)
	for _, attr := range rm.Resource.Attributes() {
		attrs[attr.Key] = attr.Value.AsString()
	}
	return attrs
}
