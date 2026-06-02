package test

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/time"
)

// LifecycleOption returns a DI option that reports lifecycle start and stop events.
func LifecycleOption(started chan<- struct{}, stopped chan<- error) di.Option {
	return di.Module(
		di.Constructor(slog.Default),
		di.Register(func(lc di.Lifecycle) {
			lc.Append(di.Hook{
				OnStart: func(context.Context) error {
					// Signal readiness after Start has had a chance to return so cancellation
					// exercises the post-start shutdown path instead of racing startup.
					go func() {
						time.Sleep(10 * time.Millisecond)
						close(started)
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					stopped <- ctx.Err()
					return nil
				},
			})
		}),
	)
}

// ShutdownExitCodeOption returns a DI option that requests shutdown with code.
func ShutdownExitCodeOption(code int) di.Option {
	return di.Module(
		di.NoLogger,
		di.Constructor(slog.Default),
		di.Register(func(lc di.Lifecycle, sh di.Shutdowner) {
			lc.Append(di.Hook{
				OnStart: func(context.Context) error {
					go func() {
						time.Sleep(10 * time.Millisecond)
						_ = sh.Shutdown(di.ExitCode(code))
					}()
					return nil
				},
			})
		}),
	)
}

// ServerFailureOption returns a DI option that registers a server whose Serve method fails.
func ServerFailureOption() di.Option {
	return di.Module(
		di.NoLogger,
		di.Constructor(slog.Default),
		di.Register(func(lc di.Lifecycle, sh di.Shutdowner) {
			server.Register(lc, []*server.Service{
				server.NewService("test", DelayedFailingServer{}, nil, sh),
			})
		}),
	)
}

// DelayedFailingServer is a [server.Server] test double whose Serve method fails after a short delay.
type DelayedFailingServer struct{}

// Serve returns ErrFailed after a short delay.
func (DelayedFailingServer) Serve() error {
	time.Sleep(10 * time.Millisecond)

	return ErrFailed
}

// Shutdown implements [server.Server] and always succeeds.
func (DelayedFailingServer) Shutdown(context.Context) error {
	return nil
}

// String returns a stable identifier for logs and assertions.
func (DelayedFailingServer) String() string {
	return "test"
}
