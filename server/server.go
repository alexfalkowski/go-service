package server

import (
	"context"

	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/transport/meta"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Server is a generic server.
type Server struct {
	srv    Serverer
	sh     fx.Shutdowner
	logger *logger.Logger
	svc    string
}

// NewServer that can start and stop an underlying server.
func NewServer(svc string, srv Serverer, logger *logger.Logger, sh fx.Shutdowner) *Server {
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

	s.log(func(l *logger.Logger) {
		l.Info("starting server", addr, zap.String(meta.ServiceKey, s.svc))
	})

	if err := s.srv.Serve(); err != nil {
		serr := s.sh.Shutdown()

		s.log(func(l *logger.Logger) {
			l.Error("could not start server", zap.String(meta.ServiceKey, s.svc), addr, zap.Error(err), zap.NamedError("shutdown_error", serr))
		})
	}
}

// Stop the server.
func (s *Server) Stop(ctx context.Context) {
	if !s.srv.IsEnabled() {
		return
	}

	err := s.srv.Shutdown(ctx)

	s.log(func(l *logger.Logger) {
		l.Info("stopping server", zap.String(meta.ServiceKey, s.svc), zap.Error(err))
	})
}

func (s *Server) log(fn func(l *logger.Logger)) {
	if s.logger == nil {
		return
	}

	fn(s.logger)
}
