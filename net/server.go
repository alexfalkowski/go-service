package net

import (
	"context"

	tm "github.com/alexfalkowski/go-service/transport/meta"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Server for net.
type Server struct {
	svc      string
	serverer Serverer
	logger   *zap.Logger
	sh       fx.Shutdowner
}

// NewServer for net.
func NewServer(svc string, serverer Serverer, logger *zap.Logger, sh fx.Shutdowner) *Server {
	return &Server{svc: svc, serverer: serverer, logger: logger, sh: sh}
}

// Start the server.
func (s *Server) Start() {
	if !s.serverer.IsEnabled() {
		return
	}

	go s.start()
}

func (s *Server) start() {
	addr := zap.Stringer("addr", s.serverer)
	s.logger.Info("starting server", addr, zap.String(tm.ServiceKey, s.svc))

	if err := s.serverer.Serve(); err != nil {
		s.logger.Error("could not start server",
			zap.String(tm.ServiceKey, s.svc), addr,
			zap.Error(err), zap.NamedError("shutdown_error", s.sh.Shutdown()))
	}
}

// Stop the server.
func (s *Server) Stop(ctx context.Context) {
	if !s.serverer.IsEnabled() {
		return
	}

	s.logger.Info("stopping server",
		zap.String(tm.ServiceKey, s.svc), zap.Error(s.serverer.Shutdown(ctx)))
}
