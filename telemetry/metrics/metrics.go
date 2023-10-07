package metrics

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/os"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	m "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/fx"
)

// Register metrics.
func Register(server *shttp.Server) error {
	handler := promhttp.Handler()

	return server.Mux.HandlePath("GET", "/metrics", func(w http.ResponseWriter, r *http.Request, p map[string]string) {
		handler.ServeHTTP(w, r)
	})
}

// NewMeter with otel.
func NewMeter(fc fx.Lifecycle) (m.Meter, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter(os.ExecutableName())

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
