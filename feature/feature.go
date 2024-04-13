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
func NewClient(params ClientParams) (*openfeature.Client, error) {
	p, err := provider(params)
	if err != nil {
		return nil, err
	}

	openfeature.SetProvider(p)

	return openfeature.NewClient(os.ExecutableName()), nil
}

func provider(params ClientParams) (openfeature.FeatureProvider, error) {
	c := params.Config

	if c != nil && c.Kind == "flipt" {
		opts := []grpc.ClientOption{
			grpc.WithClientLogger(params.Logger), grpc.WithClientTracer(params.Tracer),
			grpc.WithClientMetrics(params.Meter), grpc.WithClientRetry(c.Retry),
			grpc.WithClientUserAgent(c.UserAgent),
		}

		is, err := grpc.UnaryClientInterceptors(opts...)
		if err != nil {
			return nil, err
		}

		svc := transport.New(transport.WithAddress(c.Host), transport.WithUnaryClientInterceptor(is...))

		return flipt.NewProvider(flipt.WithAddress(c.Host), flipt.WithService(svc)), nil
	}

	return openfeature.NoopProvider{}, nil
}
