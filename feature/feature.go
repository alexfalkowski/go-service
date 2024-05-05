package feature

import (
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/transport/grpc"
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

	Config *Config
	Logger *zap.Logger
	Tracer trace.Tracer
	Meter  metric.Meter
}

// NewClient for feature.
func NewClient(params ClientParams) *openfeature.Client {
	openfeature.SetProvider(provider(params))

	return openfeature.NewClient(os.ExecutableName())
}

func provider(params ClientParams) openfeature.FeatureProvider {
	if !IsEnabled(params.Config) {
		return openfeature.NoopProvider{}
	}

	if params.Config.IsFlipt() {
		is := grpc.UnaryClientInterceptors(grpc.WithClientLogger(params.Logger),
			grpc.WithClientTracer(params.Tracer),
			grpc.WithClientMetrics(params.Meter),
			grpc.WithClientRetry(params.Config.Retry),
			grpc.WithClientUserAgent(params.Config.UserAgent))
		svc := transport.New(transport.WithAddress(params.Config.Host), transport.WithUnaryClientInterceptor(is...))

		return flipt.NewProvider(flipt.WithAddress(params.Config.Host), flipt.WithService(svc))
	}

	return openfeature.NoopProvider{}
}
