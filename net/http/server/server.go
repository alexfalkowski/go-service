package server

import (
	"crypto/tls"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/config"
	"github.com/alexfalkowski/go-service/v2/net/http/errors"
	"github.com/alexfalkowski/go-service/v2/server"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
)

// Service is an alias for server.Service.
type Service = server.Service

// NewService for http.
func NewService(name string, http *http.Server, cfg *config.Config, logger *logger.Logger, sh di.Shutdowner) (*Service, error) {
	serv, err := NewServer(http, cfg)
	if err != nil {
		return nil, err
	}

	return server.NewService(name, serv, logger, sh), nil
}

// NewServer for http.
func NewServer(server *http.Server, cfg *config.Config) (*Server, error) {
	srv := &Server{server: server, tls: cfg.TLS}
	n, a, _ := net.SplitNetworkAddress(cfg.Address)

	l, err := net.Listen(context.Background(), n, a)
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
	return errors.ServerError(s.serve())
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

		return s.server.ServeTLS(s.listener, strings.Empty, strings.Empty)
	}

	return s.server.Serve(s.listener)
}
