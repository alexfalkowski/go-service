package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/v2/time"
)

const (
	// MethodDelete is an alias of http.MethodDelete.
	MethodDelete = http.MethodDelete

	// MethodGet is an alias of http.MethodGet.
	MethodGet = http.MethodGet

	// MethodPatch is an alias of http.MethodPatch.
	MethodPatch = http.MethodPatch

	// MethodPost is an alias of http.MethodPost.
	MethodPost = http.MethodPost

	// MethodPut is an alias of http.MethodPut.
	MethodPut = http.MethodPut
)

type (
	// Client is an alias for http.Client.
	Client = http.Client

	// Handler is an alias for http.Handler.
	Handler = http.Handler

	// Header is an alias for http.Header.
	Header = http.Header

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
	// NewRequestWithContext is an alias for http.NewRequestWithContext.
	NewRequestWithContext = http.NewRequestWithContext

	// NewServeMux is an alias for http.NewServeMux.
	NewServeMux = http.NewServeMux

	// NoBody is an alias for http.NoBody.
	NoBody = http.NoBody

	// StatusText is an alias for http.StatusText.
	StatusText = http.StatusText

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
