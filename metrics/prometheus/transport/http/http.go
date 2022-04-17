package http

import (
	"net/http"

	shttp "github.com/alexfalkowski/go-service/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Register metrics for HTTP.
func Register(server *shttp.Server) error {
	handler := promhttp.Handler()

	return server.Mux.HandlePath("GET", "/metrics", func(w http.ResponseWriter, r *http.Request, p map[string]string) {
		handler.ServeHTTP(w, r)
	})
}
