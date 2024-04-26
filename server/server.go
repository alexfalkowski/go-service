package server

import (
	"context"

	"github.com/alexfalkowski/go-service/transport/meta"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Server is a generic server.
type Server struct {
	svc    string
	srv    Serverer
	logger *zap.Logger
	sh     fx.Shutdowner
}

// NewServer that can start and stop an underlying server.
func NewServer(svc string, srv Serverer, logger *zap.Logger, sh fx.Shutdowner) *Server {
	return &Server{svc: svc, srv: srv, logger: logger, sh: sh}
}

// Start the server.
func (s *Server) Start() {
	if !s.srv.IsEnabled() {
		return
	}

	go s.start()
}

func (s *Server) start() {
	addr := zap.Stringer("addr", s.srv)
	s.logger.Info("starting server", addr, zap.String(meta.ServiceKey, s.svc))

	if err := s.srv.Serve(); err != nil {
		s.logger.Error("could not start server",
			zap.String(meta.ServiceKey, s.svc), addr,
			zap.Error(err), zap.NamedError("shutdown_error", s.sh.Shutdown()))
	}
}

// Stop the server.
func (s *Server) Stop(ctx context.Context) {
	if !s.srv.IsEnabled() {
		return
	}

	s.logger.Info("stopping server",
		zap.String(meta.ServiceKey, s.svc), zap.Error(s.srv.Shutdown(ctx)))
}
