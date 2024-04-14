package metrics

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	m "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/fx"
)

// Register metrics.
func Register(server *shttp.Server) error {
	handler := promhttp.Handler()

	return server.ServeMux().HandlePath("GET", "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		handler.ServeHTTP(w, r)
	})
}

// NewNoopMeter for metrics.
func NewNoopMeter() m.Meter {
	return noop.Meter{}
}

// NewMeter for metrics.
func NewMeter(lc fx.Lifecycle, env env.Environment, ver version.Version, cfg *Config) (m.Meter, error) {
	if !IsEnabled(cfg) {
		return NewNoopMeter(), nil
	}

	r, err := reader(cfg)
	if err != nil {
		return nil, err
	}

	name := os.ExecutableName()
	attrs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(name),
		semconv.ServiceVersion(string(ver)),
		semconv.DeploymentEnvironment(string(env)),
	)

	provider := metric.NewMeterProvider(metric.WithReader(r), metric.WithResource(attrs))
	meter := provider.Meter(name)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return provider.Shutdown(ctx)
		},
	})

	return meter, nil
}

func reader(cfg *Config) (metric.Reader, error) {
	if cfg.IsOTLP() {
		r, err := otlpmetrichttp.New(context.Background(), otlpmetrichttp.WithEndpointURL(cfg.Host))
		if err != nil {
			return nil, err
		}

		return metric.NewPeriodicReader(r), nil
	}

	return prometheus.New(prometheus.WithoutTargetInfo())
}
