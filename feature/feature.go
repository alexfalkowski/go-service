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
//
// It is intended for dependency injection (Fx/Dig). The FeatureProvider dependency is marked as optional
// so services may include feature wiring without necessarily providing a concrete provider.
type ProviderParams struct {
	di.In

	// Lifecycle is used to register start/stop hooks that set the OpenFeature provider and shut down the SDK.
	Lifecycle di.Lifecycle

	// MetricProvider is an optional OpenTelemetry meter provider used to install OpenFeature telemetry hooks.
	// When nil, telemetry hooks are not installed.
	MetricProvider metrics.MeterProvider

	// FeatureProvider is the OpenFeature provider to register.
	//
	// It is optional: if not present in the DI graph, Register is a no-op and OpenFeature uses its default
	// provider semantics.
	FeatureProvider openfeature.FeatureProvider `optional:"true"`

	// Name is the service name. It is typically used when constructing an OpenFeature client (see NewClient).
	Name env.Name
}

// Register registers an optional OpenFeature FeatureProvider with the application lifecycle.
//
// Disabled behavior: if params.FeatureProvider is nil (not provided), Register is a no-op.
//
// Enabled behavior:
//   - If a MetricProvider is available, Register installs OpenTelemetry hooks so evaluations emit metrics
//     and traces.
//   - Register appends lifecycle hooks that:
//   - set the OpenFeature provider during application start (openfeature.SetProviderAndWait), and
//   - shut down the OpenFeature SDK during application stop (openfeature.Shutdown).
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
//
// The returned client is created via openfeature.NewClient using the service name string. Callers use the
// client to evaluate feature flags/values against the currently configured provider.
func NewClient(name env.Name) *openfeature.Client {
	return openfeature.NewClient(name.String())
}
