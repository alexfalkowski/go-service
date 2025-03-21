package feature

import (
	"context"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/runtime"
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
func Register(params ProviderParams) {
	if params.FeatureProvider == nil {
		return
	}

	if params.MetricProvider != nil {
		h, err := hooks.NewMetricsHookForProvider(params.MetricProvider)
		runtime.Must(err)

		openfeature.AddHooks(h, hooks.NewTracesHook(hooks.WithErrorStatusEnabled()))
	}

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			return openfeature.SetProviderAndWait(params.FeatureProvider)
		},
		OnStop: func(_ context.Context) error {
			openfeature.Shutdown()

			return nil
		},
	})
}

// NewClient for feature.
func NewClient(name env.Name) *openfeature.Client {
	return openfeature.NewClient(name.String())
}
