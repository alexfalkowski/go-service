package cli_test

import (
	"log/slog"
	"testing"

	"github.com/alexfalkowski/go-service/v2/cli"
	"github.com/alexfalkowski/go-service/v2/config"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/module"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/client"
	"github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestApplicationRunsServerCommand(t *testing.T) {
	config := test.FilePath("configs/config.yaml")
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
	config := test.FilePath("configs/invalid_http.config.yaml")
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

func TestApplicationServerRunWithMissingEnvConfig(t *testing.T) {
	t.Setenv("MISSING_CONFIG", "")
	test.SetupCLI("server", "-config", "env:MISSING_CONFIG")

	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddServer("server", "Start the server.", test.Options()...)
			cmd.AddConfig(strings.Empty)
		},
	)

	err := app.Run(t.Context())
	require.ErrorIs(t, err, config.ErrEnvMissing)
	require.ErrorContains(t, err, "env MISSING_CONFIG")
}

func TestApplicationServerRunWithConfigFlag(t *testing.T) {
	config := test.FilePath("configs/config.yaml")
	test.SetupCLI("server", "-c", config)

	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddServer("server", "Start the server.", test.Options()...)
			cmd.AddConfig(strings.Empty)
		},
	)
	require.NoError(t, app.Run(t.Context()))
}

func TestApplicationServerDrainsStreamingRoute(t *testing.T) {
	cancel, code := startDrainingApplication(t)

	restClient := rest.NewClient(test.NewContentClient(client.WithRoundTripper(http.DefaultTransport)))
	received := make(chan string, 1)
	streamErr := make(chan error, 1)
	go func() {
		streamErr <- restClient.StreamGet(t.Context(), "http://127.0.0.1:11000/drain", rest.Options{Accept: media.NDJSON}, func(_ context.Context, stream *client.ResponseStream) error {
			var response test.Response
			if err := stream.Recv(&response); err != nil {
				return err
			}

			received <- response.Greeting

			var next test.Response
			if err := stream.Recv(&next); err != nil && !stream.IsFinished(err) {
				return err
			}

			return nil
		})
	}()

	requireStreamGreeting(t, received, streamErr)

	cancel()

	requireDrainedStream(t, streamErr)
	requireServerExit(t, code)
}

func TestApplicationServerRejectsInvalidConfiguration(t *testing.T) {
	tests := []struct {
		name   string
		config string
	}{
		{name: "rejects invalid HTTP configuration", config: test.FilePath("configs/invalid_http.config.yaml")},
		{name: "rejects invalid grpc configuration", config: test.FilePath("configs/invalid_grpc.config.yaml")},
		{name: "rejects invalid debug configuration", config: test.FilePath("configs/invalid_debug.config.yaml")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.SetupCLI("server", "-config", tt.config)

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

func TestApplicationRunsWhenServerIsDisabled(t *testing.T) {
	test.SetupCLI("server", "-config", test.FilePath("configs/disabled.config.yaml"))

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

func startDrainingApplication(t *testing.T) (context.CancelFunc, <-chan int) {
	t.Helper()

	test.SetupCLI("server", "-config", test.FilePath("configs/config.yaml"))

	started := make(chan struct{})
	opts := []di.Option{
		module.Server,
		di.Register(func(s *rest.Server) {
			s.StreamGet("/drain", func(ctx context.Context, stream *stream.Stream[test.Response]) error {
				if err := stream.Send(&test.Response{Greeting: "Hello Bob"}); err != nil {
					return err
				}

				<-ctx.Done()

				return ctx.Err()
			}, http.WithRouteUnauthenticated())
		}),
		di.Register(func(lc di.Lifecycle) {
			lc.Append(di.Hook{
				OnStart: func(context.Context) error {
					go func() {
						time.Sleep(10 * time.Millisecond)
						close(started)
					}()

					return nil
				},
			})
		}),
	}
	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddServer("server", "Start the server.", opts...)
			cmd.AddConfig(strings.Empty)
		},
	)

	ctx, cancel := context.WithCancel(t.Context())
	code := make(chan int, 1)
	go func() {
		code <- app.RunCode(ctx)
	}()

	requireServerStarted(t, started, code)

	return cancel, code
}

func requireServerStarted(t *testing.T, started <-chan struct{}, code <-chan int) {
	t.Helper()

	select {
	case <-started:
	case got := <-code:
		require.FailNowf(t, "server exited before startup completed", "exit code: %d", got)
	case <-time.After(5 * time.Second):
		require.FailNow(t, "timed out waiting for server startup")
	}
}

func requireStreamGreeting(t *testing.T, received <-chan string, streamErr <-chan error) {
	t.Helper()

	select {
	case greeting := <-received:
		require.Equal(t, "Hello Bob", greeting)
	case err := <-streamErr:
		require.FailNow(t, "stream ended before drain", err)
	case <-time.After(time.Second):
		require.FailNow(t, "timed out waiting for stream response")
	}
}

func requireDrainedStream(t *testing.T, streamErr <-chan error) {
	t.Helper()

	select {
	case err := <-streamErr:
		require.NoError(t, err)
	case <-time.After(time.Second):
		require.FailNow(t, "timed out waiting for drained stream to finish")
	}
}

func requireServerExit(t *testing.T, code <-chan int) {
	t.Helper()

	select {
	case got := <-code:
		require.Equal(t, os.ExitCodeSuccess, got)
	case <-time.After(time.Second):
		require.FailNow(t, "timed out waiting for server shutdown")
	}
}
