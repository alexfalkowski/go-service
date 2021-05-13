package http

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Register metrics for HTTP
func Register(server *http.Server) error {
	mux := server.Handler.(*runtime.ServeMux)
	handler := promhttp.Handler()

	return mux.HandlePath("GET", "/metrics", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		handler.ServeHTTP(w, r)
	})
}
