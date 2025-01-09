package grpc

import (
	"context"
	"net"

	sn "github.com/alexfalkowski/go-service/net"
	"google.golang.org/grpc"
)

// Server for gRPC.
type Server struct {
	server   *grpc.Server
	listener net.Listener
}

// Config for HTTP.
type Config struct {
	Address string
}

// NewServer for gRPC.
func NewServer(server *grpc.Server, cfg *Config) (*Server, error) {
	srv := &Server{server: server}

	if cfg == nil {
		return srv, nil
	}

	l, err := sn.Listener(cfg.Address)
	if err != nil {
		return srv, err
	}

	srv.listener = l

	return srv, nil
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

// IsEnabled for server.
func (s *Server) IsEnabled() bool {
	return s.listener != nil
}
