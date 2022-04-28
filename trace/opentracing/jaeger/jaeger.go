package jaeger

import (
	"context"
	"fmt"

	"github.com/alexfalkowski/go-service/os"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/uber/jaeger-client-go"
	jconfig "github.com/uber/jaeger-client-go/config"
	jprometheus "github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/fx"
)

const (
	eventsPerSecond = 100
)

// TracerParams for jaeger.
type TracerParams struct {
	Lifecycle fx.Lifecycle
	Name      string
	Config    *Config
}

// NewTracer for jaeger.
func NewTracer(params TracerParams) (opentracing.Tracer, error) {
	c := jconfig.Configuration{
		ServiceName: fmt.Sprintf("%s:%s", os.ExecutableName(), params.Name),
		Sampler: &jconfig.SamplerConfig{
			Type:  jaeger.SamplerTypeRateLimiting,
			Param: eventsPerSecond,
		},
		Reporter: &jconfig.ReporterConfig{
			LocalAgentHostPort: params.Config.Host,
			LogSpans:           false,
		},
	}

	options := []jconfig.Option{
		jconfig.Logger(jaeger.NullLogger),
		jconfig.Metrics(jprometheus.New(jprometheus.WithRegisterer(prometheus.NewRegistry()))),
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
