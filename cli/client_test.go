package cli_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v3/cli"
	"github.com/alexfalkowski/go-service/v3/context"
	"github.com/alexfalkowski/go-service/v3/internal/test"
	"github.com/alexfalkowski/go-service/v3/os"
	"github.com/alexfalkowski/go-service/v3/strings"
	"github.com/alexfalkowski/go-service/v3/time"
	"github.com/alexfalkowski/go-signal/v2"
	"github.com/stretchr/testify/require"
)

func TestApplicationClientRunCodeOnSuccess(t *testing.T) {
	test.SetupCLI("client")

	app := cli.NewApplication(
		func(c cli.Commander) {
			c.AddClient("client", "Start the client.")
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

func TestApplicationRunsClientCommand(t *testing.T) {
	test.SetupCLI("client")

	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddClient("client", "Start the client.")
			cmd.AddConfig(strings.Empty)
		},
	)
	require.NoError(t, app.Run(t.Context()))
}

func TestApplicationClientCanRunTwice(t *testing.T) {
	test.SetupCLI("client")

	app := cli.NewApplication(
		func(c cli.Commander) {
			c.AddClient("client", "Start the client.")
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
				signal.Constructor(func() string {
					panic("bad client")
				}),
				signal.Register(func(string) {}),
			)
		},
	)

	err := app.Run(t.Context())
	require.Error(t, err)
	require.ErrorContains(t, err, `panic: "bad client"`)
}

func TestApplicationClientStartUsesConfiguredTimeout(t *testing.T) {
	test.SetupCLI("client")

	app := cli.NewApplication(
		func(c cli.Commander) {
			c.AddClient(
				"client",
				"Start the client.",
				signal.StartTimeout((10 * time.Millisecond).Duration()),
				signal.Register(func(lc signal.Lifecycle) {
					lc.Append(signal.Hook{
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
				signal.Register(func(lc signal.Lifecycle, sh signal.Shutdowner) {
					lc.Append(signal.Hook{
						OnStart: func(context.Context) error {
							sh.Shutdown(3)
							return nil
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

func TestApplicationClientRejectsInvalidConfiguration(t *testing.T) {
	tests := []struct {
		name   string
		config string
	}{
		{name: "rejects invalid HTTP configuration", config: test.FilePath("configs/invalid_http.config.yaml")},
		{name: "rejects invalid grpc configuration", config: test.FilePath("configs/invalid_grpc.config.yaml")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.SetupCLI("client", "-config", tt.config)

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
