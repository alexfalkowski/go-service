package server

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/grpc/config"
	"github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"google.golang.org/grpc"
)

// Service is an alias for server.Service.
type Service = server.Service

// NewService builds a service that starts and stops a gRPC server with logging and shutdown wiring.
func NewService(name string, grpc *grpc.Server, cfg *config.Config, logger *logger.Logger, sh di.Shutdowner) (*Service, error) {
	serv, err := NewServer(grpc, cfg)
	if err != nil {
		return nil, err
	}

	return server.NewService(name, serv, logger, sh), nil
}

// NewServer builds a Server that listens on cfg.Address.
func NewServer(server *grpc.Server, cfg *config.Config) (*Server, error) {
	srv := &Server{server: server}
	n, a, _ := net.SplitNetworkAddress(cfg.Address)

	l, err := net.Listen(context.Background(), n, a)
	if err != nil {
		return srv, err
	}

	srv.listener = l
	return srv, nil
}

// Server wraps a grpc.Server with its listener.
type Server struct {
	server   *grpc.Server
	listener net.Listener
}

// Serve starts serving requests on the configured listener.
func (s *Server) Serve() error {
	return s.server.Serve(s.listener)
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(_ context.Context) error {
	s.server.GracefulStop()
	return nil
}

// String returns the listener address.
func (s *Server) String() string {
	return s.listener.Addr().String()
}
