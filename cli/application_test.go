package cli_test

import (
	"log/slog"
	"testing"

	"github.com/alexfalkowski/go-service/v3/cli"
	"github.com/alexfalkowski/go-service/v3/context"
	"github.com/alexfalkowski/go-service/v3/errors"
	"github.com/alexfalkowski/go-service/v3/internal/test"
	"github.com/alexfalkowski/go-service/v3/runtime"
	"github.com/alexfalkowski/go-signal/v2"
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
