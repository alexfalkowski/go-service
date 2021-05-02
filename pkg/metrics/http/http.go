package http

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Register metrics for HTTP
func Register(mux *runtime.ServeMux) error {
	handler := promhttp.Handler()

	return mux.HandlePath("GET", "/metrics", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		handler.ServeHTTP(w, r)
	})
}
