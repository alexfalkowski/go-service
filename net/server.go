package net

import (
	"context"
	"net"

	tm "github.com/alexfalkowski/go-service/transport/meta"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Server for net.
type Server struct {
	svc      string
	serverer Serverer
	listener net.Listener
	logger   *zap.Logger
	sh       fx.Shutdowner
}

// NewServer for net.
func NewServer(svc string, serverer Serverer, listener net.Listener, logger *zap.Logger, sh fx.Shutdowner) *Server {
	return &Server{svc: svc, serverer: serverer, listener: listener, logger: logger, sh: sh}
}

// Start the server.
func (s *Server) Start() {
	if s.listener == nil {
		return
	}

	go s.start()
}

func (s *Server) start() {
	s.logger.Info("starting server", zap.Stringer("addr", s.listener.Addr()), zap.String(tm.ServiceKey, s.svc))

	if err := s.serverer.Serve(); err != nil {
		fields := []zapcore.Field{zap.Stringer("addr", s.listener.Addr()), zap.Error(err), zap.String(tm.ServiceKey, s.svc)}

		if err := s.sh.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		s.logger.Error("could not start server", fields...)
	}
}

// Stop the server.
func (s *Server) Stop(ctx context.Context) {
	if s.listener == nil {
		return
	}

	err := s.serverer.Shutdown(ctx)
	s.logger.Info("stopping server", zap.String(tm.ServiceKey, s.svc), zap.Error(err))
}
