package metrics

import (
	"context"

	"github.com/alexfalkowski/go-service/os"
	otlp "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

// NewReader for metrics. A nil reader means disabled.
func NewReader(fs os.FileSystem, cfg *Config) (metric.Reader, error) {
	switch {
	case !IsEnabled(cfg):
		return nil, nil
	case cfg.IsOTLP():
		if err := cfg.Headers.Secrets(fs); err != nil {
			return nil, prefix(err)
		}

		r, err := otlp.New(context.Background(), otlp.WithEndpointURL(cfg.URL), otlp.WithHeaders(cfg.Headers))

		return metric.NewPeriodicReader(r), prefix(err)
	case cfg.IsPrometheus():
		e, err := prometheus.New()

		return e, prefix(err)
	default:
		return nil, nil
	}
}
