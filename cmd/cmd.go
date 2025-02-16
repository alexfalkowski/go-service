package cmd

import (
	"context"

	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/flags"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/time"
	"github.com/cristalhq/acmd"
	"go.uber.org/dig"
	"go.uber.org/fx"
)

// New command.
func New(version string) *Command {
	name := os.ExecutableName()

	return &Command{name: name, cmds: []acmd.Command{}, version: version}
}

// Command for application.
type Command struct {
	name    string
	version string
	cmds    []acmd.Command
}

// AddServer to the command.
func (c *Command) AddServer(name, description string, flags *flags.FlagSet, opts ...fx.Option) {
	cmd := acmd.Command{
		Name:        name,
		Description: description,
		ExecFunc: func(ctx context.Context, args []string) error {
			if err := flags.Parse(args); err != nil {
				return err
			}

			return RunServer(ctx, name, opts...)
		},
	}

	c.cmds = append(c.cmds, cmd)
}

// AddClient to the command.
func (c *Command) AddClient(name, description string, flags *flags.FlagSet, opts ...fx.Option) {
	cmd := acmd.Command{
		Name:        name,
		Description: description,
		ExecFunc: func(ctx context.Context, args []string) error {
			if err := flags.Parse(args); err != nil {
				return err
			}

			return RunClient(ctx, name, opts...)
		},
	}

	c.cmds = append(c.cmds, cmd)
}

// Run the command.
func (c *Command) Run(args ...string) error {
	if len(args) == 0 {
		args = flags.Sanitize(os.Args)
	}

	runner := acmd.RunnerOf(c.cmds, acmd.Config{
		AppName:        c.name,
		AppDescription: c.name,
		Version:        c.version,
		Args:           args,
	})

	return runner.Run()
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
