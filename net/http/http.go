package http

import (
	"context"
	"net"
	"net/http"
)

// ResponseWriter with status for http.
type ResponseWriter struct {
	StatusCode int

	http.ResponseWriter
}

// WriteHeader sends an HTTP response header with the provided status code.
func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// Server for HTTP.
type Server struct {
	server            *http.Server
	listener          net.Listener
	certFile, keyFile string
}

// NewServer for HTTP.
func NewServer(server *http.Server, listener net.Listener, certFile, keyFile string) *Server {
	return &Server{server: server, listener: listener, certFile: certFile, keyFile: keyFile}
}

// Serve the underlying server.
func (s *Server) Serve() error {
	if s.certFile != "" && s.keyFile != "" {
		return s.server.ServeTLS(s.listener, s.certFile, s.keyFile)
	}

	return s.server.Serve(s.listener)
}

// Shutdown the underlying server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
