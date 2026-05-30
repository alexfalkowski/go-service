package cli_test

import (
	"log/slog"
	"testing"

	"github.com/alexfalkowski/go-service/v2/cli"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestApplicationRun(t *testing.T) {
	config := test.FilePath("configs/config.yml")

	os.Args = []string{test.Name.String(), "server", "-config", config}
	cli.Name = test.Name
	cli.Version = test.Version

	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddServer("server", "Start the server.", test.Options()...)
			cmd.AddConfig(strings.Empty)
		},
	)
	require.NoError(t, app.Run(t.Context()))
}

func TestApplicationRunCodeWithError(t *testing.T) {
	config := test.FilePath("configs/invalid_http.config.yml")

	os.Args = []string{test.Name.String(), "server", "-config", config}

	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddServer("server", "Start the server.", test.Options()...)
			cmd.AddConfig(strings.Empty)
		},
	)

	require.Equal(t, os.ExitCodeFailure, app.RunCode(t.Context()))
}

func TestApplicationRunCodeOnSuccess(t *testing.T) {
	os.Args = []string{test.Name.String(), "client"}
	cli.Name = test.Name
	cli.Version = test.Version

	app := cli.NewApplication(
		func(c cli.Commander) {
			c.AddClient("client", "Start the client.", di.NoLogger)
		},
	)

	require.Equal(t, os.ExitCodeSuccess, app.RunCode(t.Context()))
}

func TestApplicationRunWithInvalidFlag(t *testing.T) {
	os.Args = []string{test.Name.String(), "server", "--invalid-flag"}
	cli.Name = test.Name
	cli.Version = test.Version

	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddServer("server", "Start the server.", test.Options()...)
			cmd.AddConfig(strings.Empty)
		},
	)
	require.Error(t, app.Run(t.Context()))

	os.Args = []string{test.Name.String(), "client", "--invalid-flag"}
	cli.Name = test.Name
	cli.Version = test.Version

	app = cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddClient("client", "Start the client.", test.Options()...)
			cmd.AddConfig(strings.Empty)
		},
	)
	require.Error(t, app.Run(t.Context()))
}

func TestApplicationDuplicateCommand(t *testing.T) {
	var err error

	func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				err = runtime.ConvertRecover(recovered)
			}
		}()

		cli.NewApplication(func(c cli.Commander) {
			c.AddServer("server", "Start the server.", test.Options()...)
			c.AddClient("server", "Start the client.", test.Options()...)
		})
	}()

	require.ErrorIs(t, err, cli.ErrCommandRegistered)
	require.ErrorContains(t, err, "server")
}

func TestApplicationRunWithInvalidParams(t *testing.T) {
	config := test.FilePath("configs/config.yml")

	os.Args = []string{test.Name.String(), "server", "-config", config}
	cli.Name = test.Name
	cli.Version = test.Version

	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddServer("server", "Start the server.", test.Options()...)
			cmd.AddConfig(strings.Empty)
		},
	)
	require.NoError(t, app.Run(t.Context()))
}

func TestApplicationInvalid(t *testing.T) {
	configs := []string{
		test.FilePath("configs/invalid_http.config.yml"),
		test.FilePath("configs/invalid_grpc.config.yml"),
		test.FilePath("configs/invalid_debug.config.yml"),
	}

	for _, config := range configs {
		t.Run(config, func(t *testing.T) {
			os.Args = []string{test.Name.String(), "server", "-config", config}
			cli.Name = test.Name
			cli.Version = test.Version

			app := cli.NewApplication(
				func(c cli.Commander) {
					cmd := c.AddServer("server", "Start the server.", test.Options()...)
					cmd.AddConfig(strings.Empty)
				},
			)

			err := app.Run(t.Context())
			require.Error(t, err)
			require.Contains(t, err.Error(), "unknown port")
		})
	}
}

func TestApplicationDisabled(t *testing.T) {
	os.Args = []string{test.Name.String(), "server", "-config", test.FilePath("configs/disabled.config.yml")}
	cli.Name = test.Name
	cli.Version = test.Version

	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddServer("server", "Start the server.", test.Options()...)
			cmd.AddConfig(strings.Empty)
		},
	)
	require.NoError(t, app.Run(t.Context()))
}

func TestApplicationClient(t *testing.T) {
	os.Args = []string{test.Name.String(), "client"}
	cli.Name = test.Name
	cli.Version = test.Version

	opts := []di.Option{di.NoLogger}
	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddClient("client", "Start the client.", opts...)
			cmd.AddConfig(strings.Empty)
		},
	)
	require.NoError(t, app.Run(t.Context()))
}

func TestApplicationClientCanRunTwice(t *testing.T) {
	os.Args = []string{test.Name.String(), "client"}
	cli.Name = test.Name
	cli.Version = test.Version

	app := cli.NewApplication(
		func(c cli.Commander) {
			c.AddClient("client", "Start the client.", di.NoLogger)
		},
	)

	require.NoError(t, app.Run(t.Context()))
	require.NoError(t, app.Run(t.Context()))
}

func TestApplicationClientRecoversFromPanic(t *testing.T) {
	os.Args = []string{test.Name.String(), "client"}
	cli.Name = test.Name
	cli.Version = test.Version

	app := cli.NewApplication(
		func(c cli.Commander) {
			c.AddClient(
				"client",
				"Start the client.",
				di.Constructor(func() string {
					panic("bad client")
				}),
				di.Register(func(string) {}),
			)
		},
	)

	err := app.Run(t.Context())
	require.Error(t, err)
	require.Contains(t, err.Error(), `panic: "bad client"`)
}

func TestApplicationServerHonorsContextCancellation(t *testing.T) {
	os.Args = []string{test.Name.String(), "server"}
	cli.Name = test.Name
	cli.Version = test.Version

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

func TestApplicationServerShutdownExitCodeIsReturned(t *testing.T) {
	os.Args = []string{test.Name.String(), "server"}
	cli.Name = test.Name
	cli.Version = test.Version

	app := cli.NewApplication(
		func(c cli.Commander) {
			c.AddServer("server", "Start the server.", test.ShutdownExitCodeOption(3))
		},
	)

	require.Equal(t, 3, app.RunCode(t.Context()))
}

func TestApplicationServerShutdownExitCodeIsReturnedWhenStopFails(t *testing.T) {
	os.Args = []string{test.Name.String(), "server"}
	cli.Name = test.Name
	cli.Version = test.Version

	app := cli.NewApplication(
		func(c cli.Commander) {
			c.AddServer("server", "Start the server.", shutdownExitCodeAndStopErrorOption(3))
		},
	)

	err := app.Run(t.Context())
	require.ErrorIs(t, err, test.ErrFailed)
	require.Equal(t, 3, app.RunCode(t.Context()))
}

func TestApplicationClientShutdownExitCodeIsReturned(t *testing.T) {
	os.Args = []string{test.Name.String(), "client"}
	cli.Name = test.Name
	cli.Version = test.Version

	app := cli.NewApplication(
		func(c cli.Commander) {
			c.AddClient(
				"client",
				"Start the client.",
				di.Register(func(lc di.Lifecycle, sh di.Shutdowner) {
					lc.Append(di.Hook{
						OnStart: func(context.Context) error {
							return sh.Shutdown(di.ExitCode(3))
						},
					})
				}),
			)
		},
	)

	require.Equal(t, 3, app.RunCode(t.Context()))
}

func TestApplicationClientShutdownExitCodeIsReturnedWhenStopFails(t *testing.T) {
	os.Args = []string{test.Name.String(), "client"}
	cli.Name = test.Name
	cli.Version = test.Version

	app := cli.NewApplication(
		func(c cli.Commander) {
			c.AddClient("client", "Start the client.", shutdownExitCodeAndStopErrorOption(3))
		},
	)

	err := app.Run(t.Context())
	require.ErrorIs(t, err, test.ErrFailed)
	require.Equal(t, 3, app.RunCode(t.Context()))
}

func TestApplicationServerServeFailureReturnsServeFailureExitCode(t *testing.T) {
	os.Args = []string{test.Name.String(), "server"}
	cli.Name = test.Name
	cli.Version = test.Version

	app := cli.NewApplication(
		func(c cli.Commander) {
			c.AddServer("server", "Start the server.", test.ServerFailureOption())
		},
	)

	require.Equal(t, os.ExitCodeServeFailure, app.RunCode(t.Context()))
}

func TestApplicationInvalidClient(t *testing.T) {
	configs := []string{
		test.FilePath("configs/invalid_http.config.yml"),
		test.FilePath("configs/invalid_grpc.config.yml"),
	}

	for _, config := range configs {
		t.Run(config, func(t *testing.T) {
			os.Args = []string{test.Name.String(), "client", "-config", config}
			cli.Name = test.Name
			cli.Version = test.Version

			app := cli.NewApplication(
				func(c cli.Commander) {
					cmd := c.AddClient("client", "Start the client.", test.Options()...)
					cmd.AddConfig(strings.Empty)
				},
			)

			err := app.Run(t.Context())
			require.Error(t, err)
			require.Contains(t, err.Error(), "unknown port")
		})
	}
}

func shutdownExitCodeAndStopErrorOption(code int) di.Option {
	return di.Module(
		di.NoLogger,
		di.Constructor(slog.Default),
		di.Register(func(lc di.Lifecycle, sh di.Shutdowner) {
			lc.Append(di.Hook{
				OnStart: func(context.Context) error {
					return sh.Shutdown(di.ExitCode(code))
				},
				OnStop: func(context.Context) error {
					return test.ErrFailed
				},
			})
		}),
	)
}
