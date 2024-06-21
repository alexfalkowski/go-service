package http

import (
	"net/http"

	sh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	prom "github.com/prometheus/client_golang/prometheus/promhttp"
)

// RegisterMetrics for HTTP.
func RegisterMetrics(cfg *metrics.Config, mux sh.ServeMux) error {
	if !metrics.IsEnabled(cfg) || !cfg.IsPrometheus() {
		return nil
	}

	handler := prom.Handler()

	return mux.Handle("GET", "/metrics", func(res http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(res, req)
	})
}
