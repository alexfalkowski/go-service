package metrics

import (
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	prometheus "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Register registers the Prometheus metrics endpoint on mux when metrics are enabled and Prometheus is selected.
//
// Routing:
// The handler is registered as a GET route using the pattern built from `http.Pattern(name, "/metrics")`.
// This results in a service-prefixed route (for example, `/<service>/metrics`).
//
// Enablement:
// Registration is a no-op unless cfg is both enabled (`cfg.IsEnabled()`) and configured for Prometheus
// export (`cfg.IsPrometheus()`).
//
// Handler:
// The handler is provided by `promhttp.Handler()` and serves the Prometheus text exposition format.
func Register(name env.Name, cfg *metrics.Config, mux *http.ServeMux) {
	if cfg.IsEnabled() && cfg.IsPrometheus() {
		mux.Handle("GET "+http.Pattern(name, "/metrics"), prometheus.Handler())
	}
}
