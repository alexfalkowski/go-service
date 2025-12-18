package cli_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/cli"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestApplicationRun(t *testing.T) {
	config := test.FilePath("configs/config.yml")

	os.Args = []string{test.Name.String(), "server", "-i", config}
	cli.Name = test.Name
	cli.Version = test.Version

	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddServer("server", "Start the server.", test.Options()...)
			cmd.AddInput(strings.Empty)
		},
	)
	require.NoError(t, app.Run(t.Context()))
}

func TestApplicationExitOnRun(t *testing.T) {
	config := test.FilePath("configs/invalid_http.config.yml")

	os.Args = []string{test.Name.String(), "server", "-i", config}

	var exitCode int
	exit := func(code int) {
		exitCode = code
	}

	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddServer("server", "Start the server.", test.Options()...)
			cmd.AddInput(strings.Empty)
		},
		cli.WithApplicationExit(exit),
	)

	app.ExitOnError(t.Context())
	require.Equal(t, 1, exitCode)
}

func TestApplicationRunWithInvalidFlag(t *testing.T) {
	os.Args = []string{test.Name.String(), "server", "--invalid-flag"}
	cli.Name = test.Name
	cli.Version = test.Version

	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddServer("server", "Start the server.", test.Options()...)
			cmd.AddInput(strings.Empty)
		},
	)
	require.Error(t, app.Run(t.Context()))

	os.Args = []string{test.Name.String(), "client", "--invalid-flag"}
	cli.Name = test.Name
	cli.Version = test.Version

	app = cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddClient("client", "Start the client.", test.Options()...)
			cmd.AddInput(strings.Empty)
		},
	)
	require.Error(t, app.Run(t.Context()))
}

func TestApplicationRunWithInvalidParams(t *testing.T) {
	config := test.FilePath("configs/config.yml")

	os.Args = []string{test.Name.String(), "server", "-i", config}
	cli.Name = test.Name
	cli.Version = test.Version

	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddServer("server", "Start the server.", test.Options()...)
			cmd.AddInput(strings.Empty)
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
		os.Args = []string{test.Name.String(), "server", "-i", config}
		cli.Name = test.Name
		cli.Version = test.Version

		app := cli.NewApplication(
			func(c cli.Commander) {
				cmd := c.AddServer("server", "Start the server.", test.Options()...)
				cmd.AddInput(strings.Empty)
			},
		)

		err := app.Run(t.Context())
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown port")
	}
}

func TestApplicationDisabled(t *testing.T) {
	os.Args = []string{test.Name.String(), "server", "-i", test.FilePath("configs/disabled.config.yml")}
	cli.Name = test.Name
	cli.Version = test.Version

	app := cli.NewApplication(
		func(c cli.Commander) {
			cmd := c.AddServer("server", "Start the server.", test.Options()...)
			cmd.AddInput(strings.Empty)
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
			cmd.AddInput(strings.Empty)
		},
	)
	require.NoError(t, app.Run(t.Context()))
}

func TestApplicationInvalidClient(t *testing.T) {
	configs := []string{
		test.FilePath("configs/invalid_http.config.yml"),
		test.FilePath("configs/invalid_grpc.config.yml"),
	}

	for _, config := range configs {
		os.Args = []string{test.Name.String(), "client", "-i", config}
		cli.Name = test.Name
		cli.Version = test.Version

		app := cli.NewApplication(
			func(c cli.Commander) {
				cmd := c.AddClient("client", "Start the client.", test.Options()...)
				cmd.AddInput(strings.Empty)
			},
		)

		err := app.Run(t.Context())
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown port")
	}
}
