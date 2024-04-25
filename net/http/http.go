package http

import (
	"context"
	"errors"
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
	server *http.Server
	cfg    Config
}

// Config for HTTP.
type Config struct {
	Listener net.Listener
	Security Security
}

// Security for HTTP.
type Security struct {
	Enabled           bool
	CertFile, KeyFile string
}

// NewServer for HTTP.
func NewServer(server *http.Server, cfg Config) *Server {
	return &Server{server: server, cfg: cfg}
}

// Serve the underlying server.
func (s *Server) Serve() error {
	err := s.serve()

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

// Shutdown the underlying server.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.cfg.Listener == nil {
		return nil
	}

	return s.server.Shutdown(ctx)
}

func (s *Server) String() string {
	l := s.cfg.Listener
	if l != nil {
		return l.Addr().String()
	}

	return ""
}

func (s *Server) serve() error {
	l := s.cfg.Listener
	if l == nil {
		return nil
	}

	tls := s.cfg.Security
	if !tls.Enabled {
		return s.server.Serve(l)
	}

	return s.server.ServeTLS(l, tls.CertFile, tls.KeyFile)
}
