package server

import (
	"context"

	"github.com/alexfalkowski/go-service/transport/meta"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Server is a generic server.
type Server struct {
	srv    Serverer
	sh     fx.Shutdowner
	logger *zap.Logger
	svc    string
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

	s.log(func(l *zap.Logger) {
		l.Info("starting server", addr, zap.String(meta.ServiceKey.Value(), s.svc))
	})

	if err := s.srv.Serve(); err != nil {
		serr := s.sh.Shutdown()

		s.log(func(l *zap.Logger) {
			l.Error("could not start server", zap.String(meta.ServiceKey.Value(), s.svc), addr, zap.Error(err), zap.NamedError("shutdown_error", serr))
		})
	}
}

// Stop the server.
func (s *Server) Stop(ctx context.Context) {
	if !s.srv.IsEnabled() {
		return
	}

	err := s.srv.Shutdown(ctx)

	s.log(func(l *zap.Logger) {
		l.Info("stopping server", zap.String(meta.ServiceKey.Value(), s.svc), zap.Error(err))
	})
}

func (s *Server) log(fn func(l *zap.Logger)) {
	if s.logger == nil {
		return
	}

	fn(s.logger)
}
