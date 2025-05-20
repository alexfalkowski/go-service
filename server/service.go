package server

import (
	"context"
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/transport/meta"
	"go.uber.org/fx"
)

// NewService that can start and stop an underlying server.
func NewService(name string, server Server, logger *logger.Logger, sh fx.Shutdowner) *Service {
	return &Service{name: name, server: server, logger: logger, sh: sh}
}

// Service handles the starting and stopping of a server.
type Service struct {
	server Server
	sh     fx.Shutdowner
	logger *logger.Logger
	name   string
}

// Start the server.
func (s *Service) Start() {
	go s.start()
}

func (s *Service) start() {
	addr := slog.String("addr", s.server.String())

	s.log(func(l *logger.Logger) {
		l.Info("starting server", addr, slog.String(meta.ServiceKey, s.name))
	})

	if err := s.server.Serve(); err != nil {
		_ = s.sh.Shutdown()

		s.log(func(l *logger.Logger) {
			l.Error("could not start server", slog.String(meta.ServiceKey, s.name), addr, logger.Error(err))
		})
	}
}

// Stop the server.
func (s *Service) Stop(ctx context.Context) {
	_ = s.server.Shutdown(ctx)

	s.log(func(l *logger.Logger) {
		l.Info("stopping server", slog.String(meta.ServiceKey, s.name))
	})
}

func (s *Service) log(fn func(l *logger.Logger)) {
	if s.logger == nil {
		return
	}

	fn(s.logger)
}
