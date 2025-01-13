package feature

import (
	"context"

	"github.com/alexfalkowski/go-service/env"
	hooks "github.com/open-feature/go-sdk-contrib/hooks/open-telemetry/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
)

// ProviderParams for feature.
type ProviderParams struct {
	fx.In

	Lifecycle       fx.Lifecycle
	MetricProvider  metric.MeterProvider
	FeatureProvider openfeature.FeatureProvider `optional:"true"`
	Name            env.Name
}

// Register for feature.
func Register(params ProviderParams) error {
	provider := params.FeatureProvider
	if provider == nil {
		provider = openfeature.NoopProvider{}
	}

	h, err := hooks.NewMetricsHookForProvider(params.MetricProvider)
	if err != nil {
		return err
	}

	openfeature.AddHooks(h, hooks.NewTracesHook(hooks.WithErrorStatusEnabled()))

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			return openfeature.SetProviderAndWait(provider)
		},
		OnStop: func(_ context.Context) error {
			openfeature.Shutdown()

			return nil
		},
	})

	return nil
}

// NewClient for feature.
func NewClient(name env.Name) *openfeature.Client {
	return openfeature.NewClient(string(name))
}
