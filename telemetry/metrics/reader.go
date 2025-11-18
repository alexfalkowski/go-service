package metrics

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	otlp "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

// ErrNotFound for metrics.
var ErrNotFound = errors.New("metrics: reader not found")

// NewReader for metrics. A nil reader means disabled.
func NewReader(lc di.Lifecycle, name env.Name, cfg *Config) (metric.Reader, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	switch cfg.Kind {
	case "otlp":
		exporter, err := otlp.New(context.Background(), otlp.WithEndpointURL(cfg.URL), otlp.WithHeaders(cfg.Headers))
		if err != nil {
			return nil, prefix(err)
		}

		reader := metric.NewPeriodicReader(exporter)
		lc.Append(di.Hook{
			OnStop: func(ctx context.Context) error {
				return shutdown(reader, ctx)
			},
		})
		return reader, nil
	case "prometheus":
		exporter, err := prometheus.New(prometheus.WithNamespace(name.String()))
		if err != nil {
			return nil, prefix(err)
		}

		lc.Append(di.Hook{
			OnStop: func(ctx context.Context) error {
				return shutdown(exporter, ctx)
			},
		})
		return exporter, nil
	default:
		return nil, ErrNotFound
	}
}

func shutdown(reader metric.Reader, ctx context.Context) error {
	err := reader.Shutdown(ctx)
	if err != nil {
		if errors.Is(err, metric.ErrReaderShutdown) {
			return nil
		}
		return err
	}

	return nil
}
