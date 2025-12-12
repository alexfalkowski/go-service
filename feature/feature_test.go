package feature_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/feature"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestNoProvider(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	attrs := map[string]any{"favorite_color": "blue"}

	feature.Register(feature.ProviderParams{
		Lifecycle:      lc,
		Name:           test.Name,
		MetricProvider: test.NewPrometheusMeterProvider(lc),
	})

	client := feature.NewClient(test.Name)

	lc.RequireStart()

	v, err := client.BooleanValue(t.Context(), "v2_enabled", false, openfeature.NewEvaluationContext("tim@apple.com", attrs))
	require.NoError(t, err)
	require.False(t, v)

	lc.RequireStop()
}

func TestProvider(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	attrs := map[string]any{"favorite_color": "blue"}
	feature.Register(feature.ProviderParams{
		Lifecycle:       lc,
		Name:            test.Name,
		MetricProvider:  test.NewPrometheusMeterProvider(lc),
		FeatureProvider: openfeature.NoopProvider{},
	})

	client := feature.NewClient(test.Name)

	lc.RequireStart()

	v, err := client.BooleanValue(t.Context(), "v2_enabled", false, openfeature.NewEvaluationContext("tim@apple.com", attrs))
	require.NoError(t, err)
	require.False(t, v)

	lc.RequireStop()
}
