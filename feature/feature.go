package feature

import (
	"context"
	"strings"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/google/uuid"
	flipt "github.com/open-feature/go-sdk-contrib/providers/flipt/pkg/provider"
	"github.com/open-feature/go-sdk-contrib/providers/flipt/pkg/service/transport"
	"github.com/open-feature/go-sdk/openfeature"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// RegisterParams for feature.
type ClientParams struct {
	fx.In

	Config    *Config
	Logger    *zap.Logger
	Tracer    trace.Tracer
	Meter     metric.Meter
	Name      env.Name
	UserAgent env.UserAgent
}

// NewFeatureProvider for feature.
func NewFeatureProvider(params ClientParams) openfeature.FeatureProvider {
	if !IsEnabled(params.Config) {
		return openfeature.NoopProvider{}
	}

	if params.Config.IsFlipt() {
		is := grpc.UnaryClientInterceptors(grpc.WithClientLogger(params.Logger),
			grpc.WithClientTracer(params.Tracer),
			grpc.WithClientMetrics(params.Meter),
			grpc.WithClientRetry(params.Config.Retry),
			grpc.WithClientUserAgent(string(params.UserAgent)),
			grpc.WithClientTimeout(params.Config.Timeout))
		svc := transport.New(transport.WithAddress(params.Config.Host), transport.WithUnaryClientInterceptor(is...))

		return flipt.NewProvider(flipt.WithAddress(params.Config.Host), flipt.WithService(svc))
	}

	return openfeature.NoopProvider{}
}

// NewClient for feature.
func NewClient(lc fx.Lifecycle, name env.Name, provider openfeature.FeatureProvider) *openfeature.Client {
	openfeature.SetProvider(provider)

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			openfeature.Shutdown()

			return nil
		},
	})

	return openfeature.NewClient(string(name))
}

// IsNotFoundError for feature.
func IsNotFoundError(err error) bool {
	return strings.Contains(err.Error(), string(openfeature.FlagNotFoundCode))
}

// Ping for feature.
func Ping(ctx context.Context, client *openfeature.Client) error {
	id := uuid.New().String()
	e := openfeature.NewEvaluationContext(id, nil)

	_, err := client.BooleanValue(ctx, id, false, e)
	if IsNotFoundError(err) {
		return nil
	}

	return errors.Prefix("ping", err)
}
