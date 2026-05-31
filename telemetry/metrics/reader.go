package metrics

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/internal/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

// ErrNotFound is returned when the configured metrics reader kind is unknown.
var ErrNotFound = errors.New("metrics: reader not found")

// NewReader constructs an OpenTelemetry SDK metric reader based on cfg.
//
// If metrics are disabled (`cfg == nil`), it returns (nil, nil).
//
// The constructed reader is registered with the provided lifecycle and is shut down on stop.
// If the reader was already shut down, the shutdown error is ignored.
func NewReader(lc di.Lifecycle, name env.Name, cfg *Config) (metric.Reader, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	reader, err := newReader(name, cfg)
	if err != nil {
		return nil, err
	}

	lc.Append(di.Hook{
		OnStop: func(ctx context.Context) error {
			if err := reader.Shutdown(ctx); err != nil {
				if errors.Is(err, metric.ErrReaderShutdown) {
					return nil
				}
				return err
			}

			return nil
		},
	})
	return reader, nil
}

func newReader(name env.Name, cfg *Config) (metric.Reader, error) {
	switch cfg.Kind {
	case "otlp":
		if err := otlp.ValidateEndpoint(cfg.URL, cfg.Headers); err != nil {
			return nil, prefix(err)
		}

		opts := []otlpmetrichttp.Option{otlpmetrichttp.WithHeaders(cfg.Headers)}
		if !strings.IsEmpty(cfg.URL) {
			opts = append(opts, otlpmetrichttp.WithEndpointURL(cfg.URL))
		}

		exporter, err := otlpmetrichttp.New(context.Background(), opts...)
		if err != nil {
			return nil, prefix(err)
		}
		return metric.NewPeriodicReader(exporter), nil
	case "prometheus":
		exporter, err := prometheus.New(prometheus.WithNamespace(name.String()))
		if err != nil {
			return nil, prefix(err)
		}
		return exporter, nil
	default:
		return nil, ErrNotFound
	}
}
