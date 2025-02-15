package cmd

import (
	"context"
	"strings"

	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/time"
	"github.com/leaanthony/clir"
	"go.uber.org/dig"
	"go.uber.org/fx"
)

// New command.
func New(version string) *Command {
	name := os.ExecutableName()
	cli := clir.NewCli(name, name, version)

	cli.SetErrorFunction(func(_ string, err error) error {
		// Ignore errors related to testing.
		if strings.Contains(err.Error(), "test") {
			return nil
		}

		return err
	})

	return &Command{cli: cli}
}

// Command for application.
type Command struct {
	cli *clir.Cli
}

// AddServer to root.
func (c *Command) AddServer(name, description string, opts ...fx.Option) *clir.Command {
	cmd := c.cli.NewSubCommand(name, description)

	cmd.Action(func() error {
		return RunServer(context.Background(), name, opts...)
	})

	return cmd
}

// AddClient to root.
func (c *Command) AddClient(name, description string, opts ...fx.Option) *clir.Command {
	cmd := c.cli.NewSubCommand(name, description)

	cmd.Action(func() error {
		return RunClient(context.Background(), name, opts...)
	})

	return cmd
}

// Run the command.
func (c *Command) Run(args ...string) error {
	return c.cli.Run(args...)
}

// ExitOnError will run the command and exit on error.
func (c *Command) ExitOnError(args ...string) {
	if err := c.Run(args...); err != nil {
		os.Exit(1)
	}
}

// RunServer is a long running process.
func RunServer(ctx context.Context, name string, opts ...fx.Option) error {
	app := fx.New(options(opts)...)
	done := app.Done()

	if err := app.Start(ctx); err != nil {
		return prefix(name, err)
	}

	<-done

	return prefix(name, app.Stop(ctx))
}

// RunClient is a short lived process.
func RunClient(ctx context.Context, name string, opts ...fx.Option) error {
	app := fx.New(options(opts)...)

	if err := app.Start(ctx); err != nil {
		return prefix(name, err)
	}

	return prefix(name, app.Stop(ctx))
}

func options(opts []fx.Option) []fx.Option {
	return append(opts, fx.StartTimeout(time.Minute), fx.StopTimeout(time.Minute), fx.NopLogger)
}

func prefix(p string, err error) error {
	return errors.Prefix(p, dig.RootCause(err))
}
