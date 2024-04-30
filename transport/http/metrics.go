package http

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	prom "github.com/prometheus/client_golang/prometheus/promhttp"
)

// RegisterMetrics for HTTP.
func RegisterMetrics(mux *runtime.ServeMux) error {
	handler := prom.Handler()

	return mux.HandlePath("GET", "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		handler.ServeHTTP(w, r)
	})
}
