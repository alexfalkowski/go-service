package http

import (
	"fmt"
	"net/http"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/strings"
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

	// StatusBadRequest is an alias of http.StatusBadRequest.
	StatusBadRequest = http.StatusBadRequest

	// StatusConflict is an alias of http.StatusConflict.
	StatusConflict = http.StatusConflict

	// StatusForbidden is an alias of http.StatusForbidden.
	StatusForbidden = http.StatusForbidden

	// StatusGatewayTimeout is an alias of http.StatusGatewayTimeout.
	StatusGatewayTimeout = http.StatusGatewayTimeout

	// StatusOK is an alias of http.StatusOK.
	StatusOK = http.StatusOK

	// StatusInternalServerError is an alias of http.StatusInternalServerError.
	StatusInternalServerError = http.StatusInternalServerError

	// StatusNotFound is an alias of http.StatusNotFound.
	StatusNotFound = http.StatusNotFound

	// StatusNotImplemented is an alias of http.StatusNotImplemented.
	StatusNotImplemented = http.StatusNotImplemented

	// StatusServiceUnavailable is an alias of http.StatusServiceUnavailable.
	StatusServiceUnavailable = http.StatusServiceUnavailable

	// StatusTooManyRequests is an alias of http.StatusTooManyRequests.
	StatusTooManyRequests = http.StatusTooManyRequests

	// StatusUnauthorized is an alias of http.StatusUnauthorized.
	StatusUnauthorized = http.StatusUnauthorized
)

type (
	// Client is an alias for http.Client.
	Client = http.Client

	// Handler is an alias for http.Handler.
	Handler = http.Handler

	// HandlerFunc is an alias for http.HandlerFunc.
	HandlerFunc = http.HandlerFunc

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

	// ResponseWriter is an alias for http.ResponseWriter.
	ResponseWriter = http.ResponseWriter

	// RoundTripper is an alias for http.RoundTripper.
	RoundTripper = http.RoundTripper
)

var (
	// DefaultTransport is an alias for http.DefaultTransport.
	DefaultTransport = http.DefaultTransport

	// ErrUseLastResponse is an alias for http.ErrUseLastResponse.
	ErrUseLastResponse = http.ErrUseLastResponse

	// ErrServerClosed is an alias for http.ErrServerClosed.
	ErrServerClosed = http.ErrServerClosed

	// NewRequestWithContext is an alias for http.NewRequestWithContext.
	NewRequestWithContext = http.NewRequestWithContext

	// NewServeMux is an alias for http.NewServeMux.
	NewServeMux = http.NewServeMux

	// NoBody is an alias for http.NoBody.
	NoBody = http.NoBody

	// StatusText is an alias for http.StatusText.
	StatusText = http.StatusText
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

// Path will strip / from the start.
func Path(req *http.Request) string {
	path := req.URL.Path
	if strings.IsEmpty(path) {
		return path
	}

	return path[1:]
}

// Pattern will create a pattern with the format /name/pattern.
func Pattern(pattern string, name env.Name) string {
	return fmt.Sprintf("/%s%s", name.String(), pattern)
}
