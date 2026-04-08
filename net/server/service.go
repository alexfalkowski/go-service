package server

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
)

// NewService constructs a [Service] that manages the lifetime of server.
//
// name is used for attribution in log attributes under `meta.SystemKey` and as
// the prefix when [Service.Stop] returns a shutdown error. logger may be nil;
// the resulting Service relies on the logger package's nil-safe methods and
// simply skips log emission in that case.
//
// server is expected to be a non-nil concrete implementation such as an HTTP
// or gRPC server adapter. sh is expected to be non-nil and is used to trigger
// application shutdown when server.Serve returns a non-nil error.
//
// NewService does not validate its arguments and does not start serving. Call
// [Service.Start] to launch server.Serve in the background.
func NewService(name string, server Server, logger *logger.Logger, sh di.Shutdowner) *Service {
	return &Service{name: name, server: server, logger: logger, sh: sh}
}

// Service manages starting and stopping a [Server] with optional logging and
// Fx shutdown integration.
//
// Service is a thin lifecycle adapter. It does not own listener construction or
// startup readiness checks; it only calls [Server.Serve], [Server.Shutdown], and
// [Server.String] at the appropriate lifecycle points.
//
// Concurrency and lifecycle semantics:
//   - [Service.Start] launches [Server.Serve] in a new goroutine and returns immediately.
//   - [Service.Start] does not report Serve errors to its caller. A non-nil Serve error is logged and triggers `di.Shutdowner`.
//   - A nil Serve return is treated as normal termination and does not trigger shutdown.
//   - [Service.Stop] calls [Server.Shutdown] synchronously with the provided context.
//   - Service does not guard against repeated or concurrent Start/Stop calls; callers should coordinate lifecycle transitions.
type Service struct {
	server Server
	sh     di.Shutdowner
	logger *logger.Logger
	name   string
}

// Start launches the underlying server asynchronously.
//
// Start returns immediately after spawning a goroutine that logs the server
// address and then calls [Server.Serve].
//
// If Serve returns a non-nil error, Start requests application shutdown via the
// configured `di.Shutdowner` and logs the failure. Because Start is
// asynchronous, it never returns that Serve error directly.
//
// Start performs no deduplication or synchronization. Calling it more than once
// for the same Service will start multiple goroutines and invoke Serve multiple
// times.
func (s *Service) Start() {
	go s.start()
}

func (s *Service) start() {
	addr := logger.String("addr", s.server.String())
	s.logger.Info("starting server", addr, logger.String(meta.SystemKey, s.name))

	if err := s.server.Serve(); err != nil {
		// Trigger application shutdown when serving terminates with an error.
		_ = s.sh.Shutdown()
		s.logger.Error("could not start server", logger.String(meta.SystemKey, s.name), addr, logger.Error(err))
	}
}

// Stop requests a graceful shutdown of the underlying server.
//
// The provided context controls shutdown deadlines and cancellation for
// [Server.Shutdown]. Stop logs the shutdown attempt regardless of whether the
// Service was previously started.
//
// If Shutdown returns an error, Stop logs the failure and returns the same
// error wrapped with [errors.Prefix] using the service name for attribution. A
// successful shutdown returns nil.
//
// Stop does not trigger `di.Shutdowner`; that mechanism is reserved for
// unexpected Serve failures observed by [Service.Start].
func (s *Service) Stop(ctx context.Context) error {
	s.logger.Info("stopping server", logger.String(meta.SystemKey, s.name))

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("could not stop server", logger.String(meta.SystemKey, s.name), logger.Error(err))

		return errors.Prefix(s.name, err)
	}
	return nil
}

// String returns [Server.String] unchanged.
//
// This is typically the human-readable address or identifier used by the
// underlying server and matches the value logged under the "addr" attribute
// during [Service.Start].
func (s *Service) String() string {
	return s.server.String()
}
