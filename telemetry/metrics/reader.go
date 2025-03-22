package metrics

import (
	"errors"

	"github.com/alexfalkowski/go-service/os"
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

		return metric.NewPeriodicReader(newOtlpExporter(cfg)), nil
	case cfg.IsPrometheus():
		return newPrometheusExporter(), nil
	default:
		return nil, ErrNotFound
	}
}
