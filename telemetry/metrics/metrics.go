package metrics

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	m "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/fx"
)

// Register metrics.
func Register(server *shttp.Server) error {
	handler := promhttp.Handler()

	return server.Mux.HandlePath("GET", "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		handler.ServeHTTP(w, r)
	})
}

// NewMeter with otel.
func NewMeter(fc fx.Lifecycle, env env.Environment, ver version.Version) (m.Meter, error) {
	exporter, err := prometheus.New(prometheus.WithoutTargetInfo())
	if err != nil {
		return nil, err
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	name := os.ExecutableName()
	attrs := []attribute.KeyValue{
		semconv.ServiceName(name),
		semconv.ServiceVersion(string(ver)),
		semconv.DeploymentEnvironment(string(env)),
	}
	meter := provider.Meter(os.ExecutableName(), m.WithInstrumentationVersion(string(ver)), m.WithInstrumentationAttributes(attrs...))

	fc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return provider.ForceFlush(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return provider.Shutdown(ctx)
		},
	})

	return meter, nil
}
