package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/v2/time"
)

type (
	// Client is an alias for http.Client.
	Client = http.Client

	// Handler is an alias for http.Handler.
	Handler = http.Handler

	// Request is an alias for http.Request.
	Request = http.Request

	// Response is an alias for http.Response.
	Response = http.Response

	// ServeMux is an alias for http.ServeMux.
	ServeMux = http.ServeMux

	// Server is an alias for http.Server.
	Server = http.Server

	// RoundTripper is an alias for http.RoundTripper.
	RoundTripper = http.RoundTripper
)

var (
	// NewServeMux is an alias for http.NewServeMux.
	NewServeMux = http.NewServeMux

	// ErrUseLastResponse is an alias for http.ErrUseLastResponse.
	ErrUseLastResponse = http.ErrUseLastResponse
)

// NewServer for http.
func NewServer(timeout time.Duration, handler Handler) *Server {
	return &http.Server{
		Handler:     handler,
		ReadTimeout: timeout, WriteTimeout: timeout,
		IdleTimeout: timeout, ReadHeaderTimeout: timeout,
		Protocols: Protocols(),
	}
}
