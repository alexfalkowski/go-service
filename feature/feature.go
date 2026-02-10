package feature

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	hooks "github.com/open-feature/go-sdk-contrib/hooks/open-telemetry/pkg"
	"github.com/open-feature/go-sdk/openfeature"
)

// ProviderParams defines dependencies used to register an OpenFeature provider.
type ProviderParams struct {
	di.In
	Lifecycle       di.Lifecycle
	MetricProvider  metrics.MeterProvider
	FeatureProvider openfeature.FeatureProvider `optional:"true"`
	Name            env.Name
}

// Register registers the optional OpenFeature provider with the application lifecycle.
//
// If FeatureProvider is nil, Register is a no-op.
// When MetricProvider is available, Register installs OpenTelemetry hooks for metrics and traces.
// The provider is set during application start and shutdown is called during application stop.
func Register(params ProviderParams) {
	if params.FeatureProvider == nil {
		return
	}

	if params.MetricProvider != nil {
		h, err := hooks.NewMetricsHookForProvider(params.MetricProvider)
		runtime.Must(err)

		openfeature.AddHooks(h, hooks.NewTracesHook(hooks.WithErrorStatusEnabled()))
	}

	params.Lifecycle.Append(di.Hook{
		OnStart: func(_ context.Context) error {
			return openfeature.SetProviderAndWait(params.FeatureProvider)
		},
		OnStop: func(_ context.Context) error {
			openfeature.Shutdown()

			return nil
		},
	})
}

// NewClient returns an OpenFeature client named after the service.
func NewClient(name env.Name) *openfeature.Client {
	return openfeature.NewClient(name.String())
}
