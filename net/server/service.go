package server

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
)

// NewService constructs a Service that manages the lifetime of an underlying Server.
//
// name is used for attribution in logs (meta.SystemKey). server is the concrete server implementation
// (e.g. HTTP or gRPC). logger may be nil to disable logging. sh is used to trigger application shutdown
// when the underlying Server terminates unexpectedly.
//
// NewService does not start the server; call (*Service).Start to begin serving.
func NewService(name string, server Server, logger *logger.Logger, sh di.Shutdowner) *Service {
	return &Service{name: name, server: server, logger: logger, sh: sh}
}

// Service manages starting and stopping a Server with optional logging and shutdown integration.
//
// Concurrency/lifecycle expectations:
//   - Start launches the server in a new goroutine.
//   - Stop requests a graceful shutdown and should be called during application shutdown.
//   - If Server.Serve returns a non-nil error, Service triggers application shutdown via di.Shutdowner.
type Service struct {
	server Server
	sh     di.Shutdowner
	logger *logger.Logger
	name   string
}

// Start launches the underlying server asynchronously.
//
// Start returns immediately. The underlying Server.Serve runs in a separate goroutine.
func (s *Service) Start() {
	go s.start()
}

func (s *Service) start() {
	addr := logger.String("addr", s.server.String())

	s.log(func(l *logger.Logger) {
		l.Info("starting server", addr, logger.String(meta.SystemKey, s.name))
	})

	if err := s.server.Serve(); err != nil {
		// Trigger application shutdown when serving terminates with an error.
		_ = s.sh.Shutdown()

		s.log(func(l *logger.Logger) {
			l.Error("could not start server", logger.String(meta.SystemKey, s.name), addr, logger.Error(err))
		})
	}
}

// Stop requests a graceful shutdown of the underlying server and logs the stop event.
//
// Stop calls Server.Shutdown with ctx and ignores the returned error. The provided context controls
// shutdown deadlines/cancellation.
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
