package http

import (
	"context"
	"errors"
	"net"
	"net/http"

	sn "github.com/alexfalkowski/go-service/net"
)

// Server for HTTP.
type Server struct {
	server   *http.Server
	sec      Security
	listener net.Listener
}

// Config for HTTP.
type Config struct {
	Enabled  bool
	Port     string
	Security Security
}

// Security for HTTP.
type Security struct {
	Enabled           bool
	CertFile, KeyFile string
}

// NewServer for HTTP.
func NewServer(server *http.Server, cfg Config) (*Server, error) {
	s := &Server{server: server, sec: cfg.Security}

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
	if !s.sec.Enabled {
		return s.server.Serve(s.listener)
	}

	return s.server.ServeTLS(s.listener, s.sec.CertFile, s.sec.KeyFile)
}
