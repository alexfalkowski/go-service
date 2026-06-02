package feature_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/feature"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestRegister(t *testing.T) {
	for _, tt := range []struct {
		provider openfeature.FeatureProvider
		name     string
	}{
		{name: "no provider"},
		{provider: openfeature.NoopProvider{}, name: "provider"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			lc := fxtest.NewLifecycle(t)
			meterProvider, err := test.NewPrometheusMeterProvider(lc)
			require.NoError(t, err)

			feature.Register(feature.ProviderParams{
				Lifecycle:       lc,
				Name:            test.Name,
				MetricProvider:  meterProvider,
				FeatureProvider: tt.provider,
			})

			client := feature.NewClient(test.Name)

			lc.RequireStart()

			attrs := map[string]any{"favorite_color": "blue"}
			evalCtx := openfeature.NewEvaluationContext("tim@apple.com", attrs)
			v, err := client.BooleanValue(t.Context(), "v2_enabled", false, evalCtx)
			require.NoError(t, err)
			require.False(t, v)

			lc.RequireStop()
		})
	}
}
