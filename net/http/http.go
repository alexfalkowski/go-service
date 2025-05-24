package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/v2/time"
)

type (
	// ServeMux is an alias for http.ServeMux.
	ServeMux = http.ServeMux

	// Server is an alias for http.Server.
	Server = http.Server
)

// NewServeMux is an alias for http.NewServeMux.
var NewServeMux = http.NewServeMux

// NewServer for http.
func NewServer(timeout time.Duration, handler http.Handler) *Server {
	return &http.Server{
		Handler:     handler,
		ReadTimeout: timeout, WriteTimeout: timeout,
		IdleTimeout: timeout, ReadHeaderTimeout: timeout,
		Protocols: Protocols(),
	}
}
