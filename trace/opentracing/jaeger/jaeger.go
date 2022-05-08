package jaeger

import (
	"context"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jconfig "github.com/uber/jaeger-client-go/config"
	"go.uber.org/fx"
)

const (
	eventsPerSecond = 100
)

// TracerParams for jaeger.
type TracerParams struct {
	Lifecycle fx.Lifecycle
	Name      string
	Version   version.Version
	Host      string
}

// NewTracer for jaeger.
func NewTracer(params TracerParams) (opentracing.Tracer, error) {
	c := jconfig.Configuration{
		ServiceName: params.Name,
		Sampler: &jconfig.SamplerConfig{
			Type:  jaeger.SamplerTypeRateLimiting,
			Param: eventsPerSecond,
		},
		Reporter: &jconfig.ReporterConfig{
			LocalAgentHostPort: params.Host,
			LogSpans:           false,
		},
	}

	options := []jconfig.Option{
		jconfig.Logger(jaeger.NullLogger),
		jconfig.Tag("name", os.ExecutableName()),
		jconfig.Tag("version", params.Version),
	}

	tracer, closer, err := c.NewTracer(options...)
	if err != nil {
		return nil, err
	}

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return closer.Close()
		},
	})

	return tracer, nil
}
