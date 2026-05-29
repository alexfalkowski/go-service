package metrics_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestIsEnabled(t *testing.T) {
	t.Cleanup(func() {
		metrics.NewMeterProvider(metrics.MeterProviderParams{Lifecycle: fxtest.NewLifecycle(t)})
	})

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

	rm := metrics.ResourceMetrics{}
	require.NoError(t, reader.Collect(t.Context(), &rm))
	attrs := resourceAttributes(rm)

	require.Equal(t, test.ID.String(), attrs[attributes.HostID("").Key])
	require.Equal(t, test.Name.String(), attrs[attributes.ServiceName("").Key])
	require.Equal(t, test.Version.String(), attrs[attributes.ServiceVersion("").Key])
	require.Equal(t, "development", attrs[attributes.DeploymentEnvironmentName("").Key])
}

func resourceAttributes(rm metrics.ResourceMetrics) map[attributes.Key]string {
	attrs := make(map[attributes.Key]string)
	for _, attr := range rm.Resource.Attributes() {
		attrs[attr.Key] = attr.Value.AsString()
	}
	return attrs
}
