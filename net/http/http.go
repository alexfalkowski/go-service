package http

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"

	sn "github.com/alexfalkowski/go-service/net"
)

// Server for HTTP.
type Server struct {
	server   *http.Server
	tls      *tls.Config
	listener net.Listener
}

// Config for HTTP.
type Config struct {
	Enabled bool
	Port    string
	TLS     *tls.Config
}

// NewServer for HTTP.
func NewServer(server *http.Server, cfg Config) (*Server, error) {
	s := &Server{server: server, tls: cfg.TLS}

	if !cfg.Enabled {
		return s, nil
	}

	l, err := sn.Listener(cfg.Port)
	if err != nil {
		return s, err
	}

	s.listener = l

	return s, nil
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
	return s.server.Shutdown(ctx)
}

func (s *Server) String() string {
	return s.listener.Addr().String()
}

// IsEnabled for server.
func (s *Server) IsEnabled() bool {
	return s.listener != nil
}

func (s *Server) serve() error {
	if s.tls != nil {
		s.server.TLSConfig = s.tls

		return s.server.ServeTLS(s.listener, "", "")
	}

	return s.server.Serve(s.listener)
}
