package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	prom "github.com/prometheus/client_golang/prometheus/promhttp"
)

// RegisterMetrics for HTTP.
func RegisterMetrics(cfg *metrics.Config, mux *runtime.ServeMux) error {
	if !metrics.IsEnabled(cfg) || !cfg.IsPrometheus() {
		return nil
	}

	handler := prom.Handler()

	return mux.HandlePath("GET", "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		handler.ServeHTTP(w, r)
	})
}
