package test

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v3/cli"
	"github.com/alexfalkowski/go-service/v3/context"
	"github.com/alexfalkowski/go-service/v3/net/server"
	"github.com/alexfalkowski/go-service/v3/os"
	"github.com/alexfalkowski/go-service/v3/time"
	"github.com/alexfalkowski/go-signal/v2"
)

// SetupCLI configures process arguments and CLI identity for tests.
func SetupCLI(args ...string) {
	os.Args = append([]string{Name.String()}, args...)
	cli.Name = Name
	cli.Version = Version
}

// LifecycleOption returns a DI option that reports lifecycle start and stop events.
func LifecycleOption(started chan<- struct{}, stopped chan<- error) signal.Option {
	return signal.Module(
		signal.Constructor(slog.Default),
		signal.Register(func(lc signal.Lifecycle) {
			lc.Append(signal.Hook{
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
func ShutdownExitCodeOption(code int) signal.Option {
	return signal.Module(
		signal.Constructor(slog.Default),
		signal.Register(func(lc signal.Lifecycle, sh signal.Shutdowner) {
			lc.Append(signal.Hook{
				OnStart: func(context.Context) error {
					go func() {
						time.Sleep(10 * time.Millisecond)
						sh.Shutdown(code)
					}()
					return nil
				},
			})
		}),
	)
}

// ShutdownExitCodeAndStopErrorOption returns a DI option that requests shutdown with code and fails on stop.
func ShutdownExitCodeAndStopErrorOption(code int) signal.Option {
	return signal.Module(
		signal.Constructor(slog.Default),
		signal.Register(func(lc signal.Lifecycle, sh signal.Shutdowner) {
			lc.Append(signal.Hook{
				OnStart: func(context.Context) error {
					sh.Shutdown(code)
					return nil
				},
				OnStop: func(context.Context) error {
					return ErrFailed
				},
			})
		}),
	)
}

// ServerFailureOption returns a DI option that registers a server whose Serve method fails.
func ServerFailureOption() signal.Option {
	return signal.Module(
		signal.Constructor(slog.Default),
		signal.Register(func(lc signal.Lifecycle, sh signal.Shutdowner) {
			server.Register(server.RegisterParams{
				Lifecycle: lc,
				Drain:     server.NewDrain(),
				Services: []*server.Service{
					server.NewService("test", DelayedFailingServer{}, nil, sh),
				},
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
