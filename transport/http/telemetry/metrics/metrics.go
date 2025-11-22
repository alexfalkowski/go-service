package metrics

import (
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	prometheus "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Register prometheus.
func Register(name env.Name, cfg *metrics.Config, mux *http.ServeMux) {
	if cfg.IsEnabled() && cfg.IsPrometheus() {
		mux.Handle("GET "+http.Pattern(name, "/metrics"), prometheus.Handler())
	}
}
