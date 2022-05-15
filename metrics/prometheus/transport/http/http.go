package http

import (
	"net/http"

	shttp "github.com/alexfalkowski/go-service/transport/http"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Register metrics for HTTP.
func Register(server *shttp.Server) error {
	return server.Register(register)
}

func register(mux *runtime.ServeMux) error {
	handler := promhttp.Handler()

	return mux.HandlePath("GET", "/metrics", func(w http.ResponseWriter, r *http.Request, p map[string]string) {
		handler.ServeHTTP(w, r)
	})
}
