package server

import (
	"context"
	"log/slog"

	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/transport/meta"
	"go.uber.org/fx"
)

// NewServer that can start and stop an underlying server.
func NewServer(name string, serverer Serverer, logger *logger.Logger, sh fx.Shutdowner) *Server {
	return &Server{name: name, serverer: serverer, logger: logger, sh: sh}
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
	addr := slog.String("addr", s.serverer.String())

	s.log(func(l *logger.Logger) {
		l.Info("starting server", addr, slog.String(meta.ServiceKey, s.name))
	})

	if err := s.serverer.Serve(); err != nil {
		_ = s.sh.Shutdown()

		s.log(func(l *logger.Logger) {
			l.Error("could not start server", slog.String(meta.ServiceKey, s.name), addr, slog.String("error", err.Error()))
		})
	}
}

// Stop the server.
func (s *Server) Stop(ctx context.Context) {
	if !s.serverer.IsEnabled() {
		return
	}

	_ = s.serverer.Shutdown(ctx)

	s.log(func(l *logger.Logger) {
		l.Info("stopping server", slog.String(meta.ServiceKey, s.name))
	})
}

func (s *Server) log(fn func(l *logger.Logger)) {
	if s.logger == nil {
		return
	}

	fn(s.logger)
}
