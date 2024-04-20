package http

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// RegisterMetrics for HTTP.
func RegisterMetrics(mux *runtime.ServeMux) error {
	handler := promhttp.Handler()

	return mux.HandlePath("GET", "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		handler.ServeHTTP(w, r)
	})
}
