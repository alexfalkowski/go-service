package server

import (
	"context"

	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/transport/meta"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewServer that can start and stop an underlying server.
func NewServer(svc string, srv Serverer, logger *logger.Logger, sh fx.Shutdowner) *Server {
	return &Server{name: svc, serverer: srv, logger: logger, sh: sh}
}

// Server is a generic server.
type Server struct {
	serverer Serverer
	sh       fx.Shutdowner
	logger   *logger.Logger
	name     string
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

	s.log(func(l *logger.Logger) {
		l.Info("starting server", addr, zap.String(meta.ServiceKey, s.name))
	})

	if err := s.serverer.Serve(); err != nil {
		serr := s.sh.Shutdown()

		s.log(func(l *logger.Logger) {
			l.Error("could not start server", zap.String(meta.ServiceKey, s.name), addr, zap.Error(err), zap.NamedError("shutdown_error", serr))
		})
	}
}

// Stop the server.
func (s *Server) Stop(ctx context.Context) {
	if !s.serverer.IsEnabled() {
		return
	}

	err := s.serverer.Shutdown(ctx)

	s.log(func(l *logger.Logger) {
		l.Info("stopping server", zap.String(meta.ServiceKey, s.name), zap.Error(err))
	})
}

func (s *Server) log(fn func(l *logger.Logger)) {
	if s.logger == nil {
		return
	}

	fn(s.logger)
}
