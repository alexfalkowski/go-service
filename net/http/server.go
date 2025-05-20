package http

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/alexfalkowski/go-service/v2/net"
)

// NewServer for HTTP.
func NewServer(server *http.Server, cfg *Config) (*Server, error) {
	srv := &Server{server: server, tls: cfg.TLS}

	l, err := net.Listen(cfg.Address)
	if err != nil {
		return srv, err
	}

	srv.listener = l

	return srv, nil
}

// Server for HTTP.
type Server struct {
	server   *http.Server
	tls      *tls.Config
	listener net.Listener
}

// Serve the underlying server.
func (s *Server) Serve() error {
	return ServerError(s.serve())
}

// Shutdown the underlying server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) String() string {
	return s.listener.Addr().String()
}

func (s *Server) serve() error {
	if s.tls != nil {
		s.server.TLSConfig = s.tls

		return s.server.ServeTLS(s.listener, "", "")
	}

	return s.server.Serve(s.listener)
}
