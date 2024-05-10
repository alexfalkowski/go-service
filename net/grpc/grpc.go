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
	Port    string
	Enabled bool
}

// NewServer for gRPC.
func NewServer(server *grpc.Server, cfg Config) (*Server, error) {
	s := &Server{server: server}

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
