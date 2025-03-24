package metrics

import (
	"context"
	"errors"

	"github.com/alexfalkowski/go-service/os"
	otlp "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

// ErrNotFound for metrics.
var ErrNotFound = errors.New("metrics: reader not found")

// NewReader for metrics. A nil reader means disabled.
func NewReader(fs os.FileSystem, cfg *Config) (metric.Reader, error) {
	switch {
	case !IsEnabled(cfg):
		return nil, nil
	case cfg.IsOTLP():
		if err := cfg.Headers.Secrets(fs); err != nil {
			return nil, prefix(err)
		}

		exporter, err := otlp.New(context.Background(), otlp.WithEndpointURL(cfg.URL), otlp.WithHeaders(cfg.Headers))
		if err != nil {
			return nil, prefix(err)
		}

		return metric.NewPeriodicReader(exporter), nil
	case cfg.IsPrometheus():
		exporter, err := prometheus.New()
		if err != nil {
			return nil, prefix(err)
		}

		return exporter, nil
	default:
		return nil, ErrNotFound
	}
}
