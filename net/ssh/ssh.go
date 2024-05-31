package ssh

import (
	"context"
	"errors"
	"net"

	sn "github.com/alexfalkowski/go-service/net"
	"github.com/gliderlabs/ssh"
)

// Server for SSH.
type Server struct {
	server   *ssh.Server
	listener net.Listener
}

// Config for SSH.
type Config struct {
	Port string
}

// NewServer for SSH.
func NewServer(server *ssh.Server, cfg *Config) (*Server, error) {
	s := &Server{server: server}

	if cfg == nil {
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
	err := s.serve()
	if errors.Is(err, ssh.ErrServerClosed) {
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
	return s.server.Serve(s.listener)
}
