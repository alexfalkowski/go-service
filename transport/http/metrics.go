package http

import (
	"net/http"

	hm "github.com/alexfalkowski/go-service/net/http/mux"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	prom "github.com/prometheus/client_golang/prometheus/promhttp"
)

// RegisterMetrics for HTTP.
func RegisterMetrics(cfg *metrics.Config, mux hm.ServeMux) error {
	if !metrics.IsEnabled(cfg) || !cfg.IsPrometheus() {
		return nil
	}

	handler := prom.Handler()

	return mux.Handle("GET", "/metrics", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}
