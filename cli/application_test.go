package cli_test

import (
	"log/slog"
	"testing"

	"github.com/alexfalkowski/go-service/v3/cli"
	"github.com/alexfalkowski/go-service/v3/context"
	"github.com/alexfalkowski/go-service/v3/errors"
	"github.com/alexfalkowski/go-service/v3/feature"
	"github.com/alexfalkowski/go-service/v3/internal/test"
	"github.com/alexfalkowski/go-service/v3/module"
	"github.com/alexfalkowski/go-service/v3/os"
	"github.com/alexfalkowski/go-service/v3/runtime"
	"github.com/alexfalkowski/go-service/v3/strings"
	"github.com/alexfalkowski/go-service/v3/transport/http"
	"github.com/alexfalkowski/go-signal/v2"
	"github.com/avfs/avfs/vfs/memfs"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/stretchr/testify/require"
)

func TestNewApplicationRejectsDuplicateCommand(t *testing.T) {
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

func TestApplicationConstructsFeatureProviderWithFileBackedTLS(t *testing.T) {
	config, err := test.FS.ReadFile(test.Path("configs/valid_feature_provider_tls.yaml"))
	require.NoError(t, err)
	ca, err := test.FS.ReadFile(test.Path("certs/rootCA.pem"))
	require.NoError(t, err)

	for _, tt := range []struct {
		name       string
		module     signal.Option
		addCommand func(cli.Commander, string, string, ...cli.Option) *cli.Command
	}{
		{name: "server", module: module.Server, addCommand: cli.Commander.AddServer},
		{name: "client", module: module.Client, addCommand: cli.Commander.AddClient},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fs := &os.FS{VFS: memfs.New()}
			// Restore the CLI filesystem so the next graph must register its own CA source.
			t.Cleanup(func() { http.Register(cli.FS) })
			require.NoError(t, fs.WriteFile("/feature-provider-ca.pem", ca, 0o600))
			require.NoError(t, fs.WriteFile("/config.yaml", config, 0o600))
			test.SetupCLI(tt.name, "-config", "file:/config.yaml")

			app := cli.NewApplication(func(c cli.Commander) {
				cmd := tt.addCommand(c, tt.name, "Start the application.", tt.module,
					signal.Decorate(func(*os.FS) *os.FS { return fs }),
					signal.Constructor(func(cfg *feature.Config, lc signal.Lifecycle) (openfeature.FeatureProvider, error) {
						_, _, err := http.NewClientFromConfig(cfg.Config, http.WithClientLifecycle(lc))
						return openfeature.NoopProvider{}, err
					}),
					signal.Register(func(sh signal.Shutdowner) { sh.Shutdown(os.ExitCodeSuccess) }),
				)
				cmd.AddConfig(strings.Empty)
			})
			require.NoError(t, app.Run(t.Context()))
		})
	}
}

func TestApplicationRollsBackStartupFailure(t *testing.T) {
	for _, name := range []string{"client", "server"} {
		t.Run(name, func(t *testing.T) {
			test.SetupCLI(name)
			stopped := false
			cleanupErr := errors.New("rollback failed")
			option := signal.Register(func(lc signal.Lifecycle, sh signal.Shutdowner) {
				lc.Append(signal.Hook{
					OnStop: func(context.Context) error {
						stopped = true
						return cleanupErr
					},
				})
				lc.Append(signal.Hook{
					OnStart: func(context.Context) error {
						sh.Shutdown(3)
						return test.ErrFailed
					},
				})
			})
			app := cli.NewApplication(func(c cli.Commander) {
				if name == "server" {
					c.AddServer(name, "Start the server.", signal.Constructor(slog.Default), option)
					return
				}

				c.AddClient(name, "Start the client.", option)
			})

			err := app.Run(t.Context())
			require.ErrorIs(t, err, test.ErrFailed)
			require.ErrorIs(t, err, cleanupErr)
			require.ErrorContains(t, err, name+": failed to run")
			require.True(t, stopped, "startup failure must roll back completed hooks")

			var exitErr *signal.ExitError
			require.ErrorAs(t, err, &exitErr)
			require.Equal(t, 3, exitErr.Code, "startup failure must preserve the shutdown code")
			require.Equal(t, 3, app.RunCode(t.Context()), "RunCode must preserve the shutdown code alongside startup and rollback errors")
		})
	}
}
