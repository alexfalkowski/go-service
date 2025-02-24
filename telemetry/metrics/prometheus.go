package metrics

import (
	"github.com/alexfalkowski/go-service/runtime"
	"go.opentelemetry.io/otel/exporters/prometheus"
)

func newPrometheusExporter() *prometheus.Exporter {
	exporter, err := prometheus.New()
	runtime.Must(prefix(err))

	return exporter
}
