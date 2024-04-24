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
	if s.listener == nil {
		return nil
	}

	return s.server.Serve(s.listener)
}

// Shutdown the underlying server.
func (s *Server) Shutdown(_ context.Context) error {
	if s.listener == nil {
		return nil
	}

	s.server.GracefulStop()

	return nil
}

// String for server.
func (s *Server) String() string {
	l := s.listener
	if l != nil {
		return l.Addr().String()
	}

	return ""
}
