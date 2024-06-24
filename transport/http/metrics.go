package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/telemetry/metrics"
	prom "github.com/prometheus/client_golang/prometheus/promhttp"
)

// RegisterMetrics for HTTP.
func RegisterMetrics(cfg *metrics.Config, mux *http.ServeMux) {
	if !metrics.IsEnabled(cfg) || !cfg.IsPrometheus() {
		return
	}

	handler := prom.Handler()

	mux.HandleFunc("GET /metrics", func(res http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(res, req)
	})
}
