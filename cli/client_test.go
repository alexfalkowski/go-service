package cli_test

import (
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

func TestApplicationClientRunCodeOnSuccess(t *testing.T) {
	test.SetupCLI("client")

	app := cli.NewApplication(
		func(c cli.Commander) {
			c.AddClient("client", "Start the client.", di.NoLogger)
		},
	)

	require.Equal(t, os.ExitCodeSuccess, app.RunCode(t.Context()))
}

func TestApplicationClientRunWithInvalidFlag(t *testing.T) {
	test.SetupCLI("client", "--invalid-flag")

	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddClient("client", "Start the client.", test.Options()...)
			cmd.AddConfig(strings.Empty)
		},
	)

	require.Error(t, app.Run(t.Context()))
}

func TestApplicationClient(t *testing.T) {
	test.SetupCLI("client")

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
	test.SetupCLI("client")

	app := cli.NewApplication(
		func(c cli.Commander) {
			c.AddClient("client", "Start the client.", di.NoLogger)
		},
	)

	require.NoError(t, app.Run(t.Context()))
	require.NoError(t, app.Run(t.Context()))
}

func TestApplicationClientRecoversFromPanic(t *testing.T) {
	test.SetupCLI("client")

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
	require.ErrorContains(t, err, `panic: "bad client"`)
}

func TestApplicationClientStartUsesFxTimeout(t *testing.T) {
	test.SetupCLI("client")

	app := cli.NewApplication(
		func(c cli.Commander) {
			c.AddClient(
				"client",
				"Start the client.",
				di.NoLogger,
				fx.StartTimeout((10 * time.Millisecond).Duration()),
				di.Register(func(lc di.Lifecycle) {
					lc.Append(di.Hook{
						OnStart: func(ctx context.Context) error {
							<-ctx.Done()

							return ctx.Err()
						},
					})
				}),
			)
		},
	)

	require.ErrorIs(t, app.Run(t.Context()), context.DeadlineExceeded)
}

func TestApplicationClientShutdownExitCodeIsReturned(t *testing.T) {
	test.SetupCLI("client")

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
	for _, tt := range []struct {
		name string
		code int
	}{
		{name: "positive", code: 3},
		{name: "negative", code: -1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			test.SetupCLI("client")
			app := cli.NewApplication(
				func(c cli.Commander) {
					c.AddClient("client", "Start the client.", test.ShutdownExitCodeAndStopErrorOption(tt.code))
				},
			)

			err := app.Run(t.Context())
			require.ErrorIs(t, err, test.ErrFailed)
			require.Equal(t, tt.code, app.RunCode(t.Context()))
		})
	}
}

func TestApplicationClientInvalidConfig(t *testing.T) {
	configs := []string{
		test.FilePath("configs/invalid_http.config.yml"),
		test.FilePath("configs/invalid_grpc.config.yml"),
	}

	for _, config := range configs {
		t.Run(config, func(t *testing.T) {
			test.SetupCLI("client", "-config", config)

			app := cli.NewApplication(
				func(c cli.Commander) {
					cmd := c.AddClient("client", "Start the client.", test.Options()...)
					cmd.AddConfig(strings.Empty)
				},
			)

			err := app.Run(t.Context())
			require.Error(t, err)
			require.ErrorContains(t, err, "unknown port")
		})
	}
}
