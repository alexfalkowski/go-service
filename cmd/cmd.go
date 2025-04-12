package cmd

import (
	"context"
	"log"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/time"
	"github.com/cristalhq/acmd"
	"go.uber.org/dig"
	"go.uber.org/fx"
)

// New command.
func New(name env.Name, version env.Version) *Command {
	return &Command{name: name, cmds: []acmd.Command{}, version: version}
}

// Command for application.
type Command struct {
	name    env.Name
	version env.Version
	cmds    []acmd.Command
}

// AddServer to the command.
func (c *Command) AddServer(name, description string, flags *FlagSet, opts ...fx.Option) {
	cmd := acmd.Command{
		Name:        name,
		Description: description,
		ExecFunc: func(ctx context.Context, args []string) error {
			if err := flags.Parse(args); err != nil {
				return err
			}

			opts = append(opts, fx.Provide(flags.Provide))

			return RunServer(ctx, name, opts...)
		},
	}

	c.cmds = append(c.cmds, cmd)
}

// AddClient to the command.
func (c *Command) AddClient(name, description string, flags *FlagSet, opts ...fx.Option) {
	cmd := acmd.Command{
		Name:        name,
		Description: description,
		ExecFunc: func(ctx context.Context, args []string) error {
			if err := flags.Parse(args); err != nil {
				return err
			}

			opts = append(opts, fx.Provide(flags.Provide))

			return RunClient(ctx, name, opts...)
		},
	}

	c.cmds = append(c.cmds, cmd)
}

// Run the command.
func (c *Command) Run(args ...string) error {
	if len(args) == 0 {
		args = SanitizeArgs(os.Args)
	}

	name := c.name.String()
	runner := acmd.RunnerOf(c.cmds, acmd.Config{
		AppName:        name,
		AppDescription: name,
		Version:        c.version.String(),
		Args:           args,
	})

	return runner.Run()
}

// ExitOnError will run the command and exit on error.
func (c *Command) ExitOnError(args ...string) {
	if err := c.Run(args...); err != nil {
		log.Printf("%s: failed to run: %v", c.name.String(), err)
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

func options(options []fx.Option) []fx.Option {
	return append(options, fx.StartTimeout(time.Minute), fx.StopTimeout(time.Minute), fx.NopLogger, fx.RecoverFromPanics())
}

func prefix(prefix string, err error) error {
	return errors.Prefix(prefix, dig.RootCause(err))
}
