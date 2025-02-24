package metrics

import (
	"context"

	"github.com/alexfalkowski/go-service/runtime"
	otlp "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
)

func newOtlpExporter(cfg *Config) *otlp.Exporter {
	exporter, err := otlp.New(context.Background(), otlp.WithEndpointURL(cfg.URL), otlp.WithHeaders(cfg.Headers))
	runtime.Must(prefix(err))

	return exporter
}
