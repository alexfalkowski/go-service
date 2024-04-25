package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

// Server for gRPC.
type Server struct {
	server   *grpc.Server
	listener net.Listener
}

// NewServer for gRPC.
func NewServer(server *grpc.Server, listener net.Listener) *Server {
	return &Server{server: server, listener: listener}
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
