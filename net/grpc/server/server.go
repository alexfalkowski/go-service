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

// NewService for gRPC.
func NewService(name string, grpc *grpc.Server, cfg *config.Config, logger *logger.Logger, sh di.Shutdowner) (*Service, error) {
	serv, err := NewServer(grpc, cfg)
	if err != nil {
		return nil, err
	}

	return server.NewService(name, serv, logger, sh), nil
}

// NewServer for gRPC.
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

// Server for gRPC.
type Server struct {
	server   *grpc.Server
	listener net.Listener
}

// Serve the underlying server.
func (s *Server) Serve() error {
	return s.server.Serve(s.listener)
}

// Shutdown the underlying server.
func (s *Server) Shutdown(_ context.Context) error {
	s.server.GracefulStop()
	return nil
}

// String for server.
func (s *Server) String() string {
	return s.listener.Addr().String()
}
