package cli_test

import (
	"log/slog"
	"testing"

	"github.com/alexfalkowski/go-service/v2/cli"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestApplicationServerRun(t *testing.T) {
	config := test.FilePath("configs/config.yml")
	test.SetupCLI("server", "-config", config)

	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddServer("server", "Start the server.", test.Options()...)
			cmd.AddConfig(strings.Empty)
		},
	)
	require.NoError(t, app.Run(t.Context()))
}

func TestApplicationServerRunCodeWithError(t *testing.T) {
	config := test.FilePath("configs/invalid_http.config.yml")
	test.SetupCLI("server", "-config", config)

	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddServer("server", "Start the server.", test.Options()...)
			cmd.AddConfig(strings.Empty)
		},
	)

	require.Equal(t, os.ExitCodeFailure, app.RunCode(t.Context()))
}

func TestApplicationServerRunWithInvalidFlag(t *testing.T) {
	test.SetupCLI("server", "--invalid-flag")

	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddServer("server", "Start the server.", test.Options()...)
			cmd.AddConfig(strings.Empty)
		},
	)

	require.Error(t, app.Run(t.Context()))
}

func TestApplicationServerRunWithConfigFlag(t *testing.T) {
	config := test.FilePath("configs/config.yml")
	test.SetupCLI("server", "-c", config)

	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddServer("server", "Start the server.", test.Options()...)
			cmd.AddConfig(strings.Empty)
		},
	)
	require.NoError(t, app.Run(t.Context()))
}

func TestApplicationServerInvalidConfig(t *testing.T) {
	configs := []string{
		test.FilePath("configs/invalid_http.config.yml"),
		test.FilePath("configs/invalid_grpc.config.yml"),
		test.FilePath("configs/invalid_debug.config.yml"),
	}

	for _, config := range configs {
		t.Run(config, func(t *testing.T) {
			test.SetupCLI("server", "-config", config)

			app := cli.NewApplication(
				func(c cli.Commander) {
					cmd := c.AddServer("server", "Start the server.", test.Options()...)
					cmd.AddConfig(strings.Empty)
				},
			)

			err := app.Run(t.Context())
			require.Error(t, err)
			require.ErrorContains(t, err, "unknown port")
		})
	}
}

func TestApplicationServerDisabled(t *testing.T) {
	test.SetupCLI("server", "-config", test.FilePath("configs/disabled.config.yml"))

	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddServer("server", "Start the server.", test.Options()...)
			cmd.AddConfig(strings.Empty)
		},
	)
	require.NoError(t, app.Run(t.Context()))
}

func TestApplicationServerHonorsContextCancellation(t *testing.T) {
	test.SetupCLI("server")

	started := make(chan struct{})
	stopped := make(chan error, 1)
	app := cli.NewApplication(
		func(c cli.Commander) {
			c.AddServer("server", "Start the server.", test.LifecycleOption(started, stopped))
		},
	)

	ctx, cancel := context.WithCancel(t.Context())
	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run(ctx)
	}()

	select {
	case <-started:
	case err := <-errCh:
		require.FailNow(t, "server exited before startup completed", err.Error())
	case <-time.After(time.Second):
		require.FailNow(t, "timed out waiting for server startup")
	}

	cancel()

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(time.Second):
		require.FailNow(t, "timed out waiting for server shutdown after cancellation")
	}

	select {
	case err := <-stopped:
		require.NoError(t, err)
	case <-time.After(time.Second):
		require.FailNow(t, "timed out waiting for server stop hook")
	}
}

func TestApplicationServerStopAfterCancellationUsesFxTimeout(t *testing.T) {
	test.SetupCLI("server")

	started := make(chan struct{})
	app := cli.NewApplication(
		func(c cli.Commander) {
			c.AddServer(
				"server",
				"Start the server.",
				di.NoLogger,
				di.Constructor(slog.Default),
				fx.StopTimeout((10 * time.Millisecond).Duration()),
				di.Register(func(lc di.Lifecycle) {
					lc.Append(di.Hook{
						OnStart: func(context.Context) error {
							go func() {
								time.Sleep(10 * time.Millisecond)
								close(started)
							}()

							return nil
						},
						OnStop: func(ctx context.Context) error {
							<-ctx.Done()

							return ctx.Err()
						},
					})
				}),
			)
		},
	)

	ctx, cancel := context.WithCancel(t.Context())
	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run(ctx)
	}()

	select {
	case <-started:
	case err := <-errCh:
		require.FailNow(t, "server exited before startup completed", err.Error())
	case <-time.After(time.Second):
		require.FailNow(t, "timed out waiting for server startup")
	}

	cancel()

	select {
	case err := <-errCh:
		require.ErrorIs(t, err, context.DeadlineExceeded)
	case <-time.After(time.Second):
		require.FailNow(t, "timed out waiting for server shutdown after cancellation")
	}
}

func TestApplicationServerShutdownExitCodeIsReturned(t *testing.T) {
	test.SetupCLI("server")

	app := cli.NewApplication(
		func(c cli.Commander) {
			c.AddServer("server", "Start the server.", test.ShutdownExitCodeOption(3))
		},
	)

	require.Equal(t, 3, app.RunCode(t.Context()))
}

func TestApplicationServerShutdownExitCodeIsReturnedWhenStopFails(t *testing.T) {
	for _, tt := range []struct {
		name string
		code int
	}{
		{name: "positive", code: 3},
		{name: "negative", code: -1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			test.SetupCLI("server")
			app := cli.NewApplication(
				func(c cli.Commander) {
					c.AddServer("server", "Start the server.", test.ShutdownExitCodeAndStopErrorOption(tt.code))
				},
			)

			err := app.Run(t.Context())
			require.ErrorIs(t, err, test.ErrFailed)
			require.Equal(t, tt.code, app.RunCode(t.Context()))
		})
	}
}

func TestApplicationServerServeFailureReturnsServeFailureExitCode(t *testing.T) {
	test.SetupCLI("server")

	app := cli.NewApplication(
		func(c cli.Commander) {
			c.AddServer("server", "Start the server.", test.ServerFailureOption())
		},
	)

	require.Equal(t, os.ExitCodeServeFailure, app.RunCode(t.Context()))
}
