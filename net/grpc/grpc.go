package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

// Server for HTTP.
type Server struct {
	server   *grpc.Server
	listener net.Listener
}

// NewServer for HTTP.
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
