package http

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"

	sn "github.com/alexfalkowski/go-service/net"
)

type (
	// Server for HTTP.
	Server struct {
		server   *http.Server
		tls      *tls.Config
		listener net.Listener
	}

	// Config for HTTP.
	Config struct {
		TLS     *tls.Config
		Address string
	}
)

// NewServer for HTTP.
func NewServer(server *http.Server, cfg *Config) (*Server, error) {
	srv := &Server{server: server}

	if cfg == nil {
		return srv, nil
	}

	srv.tls = cfg.TLS

	l, err := sn.Listener(cfg.Address)
	if err != nil {
		return srv, err
	}

	srv.listener = l

	return srv, nil
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
