package feature

import (
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	flipt "github.com/open-feature/go-sdk-contrib/providers/flipt/pkg/provider"
	"github.com/open-feature/go-sdk-contrib/providers/flipt/pkg/service/transport"
	"github.com/open-feature/go-sdk/openfeature"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// RegisterParams for feature.
type ClientParams struct {
	fx.In

	Config *Config
	Logger *zap.Logger
	Tracer tracer.Tracer
	Meter  metric.Meter
}

// NewClient for feature.
func NewClient(params ClientParams) (*openfeature.Client, error) {
	p, err := provider(params)
	if err != nil {
		return nil, err
	}

	openfeature.SetProvider(p)

	return openfeature.NewClient(os.ExecutableName()), nil
}

func provider(params ClientParams) (openfeature.FeatureProvider, error) {
	if params.Config.Kind == "flipt" {
		opts := []grpc.ClientOption{
			grpc.WithClientLogger(params.Logger), grpc.WithClientTracer(params.Tracer),
			grpc.WithClientMetrics(params.Meter), grpc.WithClientRetry(&params.Config.Retry),
			grpc.WithClientUserAgent(params.Config.UserAgent),
		}

		is, err := grpc.UnaryClientInterceptors(opts...)
		if err != nil {
			return nil, err
		}

		svc := transport.New(transport.WithAddress(params.Config.Host), transport.WithUnaryClientInterceptor(is...))

		return flipt.NewProvider(flipt.WithAddress(params.Config.Host), flipt.WithService(svc)), nil
	}

	return openfeature.NoopProvider{}, nil
}
