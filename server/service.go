package server

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
)

// NewService that can start and stop an underlying server.
func NewService(name string, server Server, logger *logger.Logger, sh di.Shutdowner) *Service {
	return &Service{name: name, server: server, logger: logger, sh: sh}
}

// Service handles the starting and stopping of a server.
type Service struct {
	server Server
	sh     di.Shutdowner
	logger *logger.Logger
	name   string
}

// Start the server.
func (s *Service) Start() {
	go s.start()
}

func (s *Service) start() {
	addr := logger.String("addr", s.server.String())

	s.log(func(l *logger.Logger) {
		l.Info("starting server", addr, logger.String(meta.SystemKey, s.name))
	})

	if err := s.server.Serve(); err != nil {
		_ = s.sh.Shutdown()

		s.log(func(l *logger.Logger) {
			l.Error("could not start server", logger.String(meta.SystemKey, s.name), addr, logger.Error(err))
		})
	}
}

// Stop the server.
func (s *Service) Stop(ctx context.Context) {
	_ = s.server.Shutdown(ctx)

	s.log(func(l *logger.Logger) {
		l.Info("stopping server", logger.String(meta.SystemKey, s.name))
	})
}

func (s *Service) log(fn func(l *logger.Logger)) {
	if s.logger == nil {
		return
	}

	fn(s.logger)
}
