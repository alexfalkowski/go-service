package cmd

import (
	"context"
	"log/slog"

	"github.com/alexfalkowski/go-service/cmd/flag"
	"github.com/alexfalkowski/go-service/env"
	se "github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/time"
	"github.com/cristalhq/acmd"
	"go.uber.org/dig"
	"go.uber.org/fx"
)

// NewCommand for cmd.
func NewCommand(name env.Name, version env.Version, exit os.ExitFunc) *Command {
	return &Command{name: name, exit: exit, version: version, cmds: []acmd.Command{}}
}

// Command for application.
type Command struct {
	name    env.Name
	version env.Version
	exit    os.ExitFunc
	cmds    []acmd.Command
}

// AddServer to the command.
func (c *Command) AddServer(name, description string, opts ...fx.Option) *flag.FlagSet {
	flags := flag.NewFlagSet(name)
	cmd := acmd.Command{
		Name:        name,
		Description: description,
		ExecFunc: func(ctx context.Context, args []string) error {
			if err := flags.Parse(args); err != nil {
				return err
			}

			opts = append(opts, fx.Provide(flags.Provide), runtime.Module)
			app := fx.New(fxOptions(opts)...)
			done := app.Done()

			if err := app.Start(ctx); err != nil {
				return prefix(name, err)
			}

			<-done

			return prefix(name, app.Stop(ctx))
		},
	}

	c.cmds = append(c.cmds, cmd)

	return flags
}

// AddClient to the command.
func (c *Command) AddClient(name, description string, opts ...fx.Option) *flag.FlagSet {
	flags := flag.NewFlagSet(name)
	cmd := acmd.Command{
		Name:        name,
		Description: description,
		ExecFunc: func(ctx context.Context, args []string) error {
			if err := flags.Parse(args); err != nil {
				return err
			}

			opts = append(opts, fx.Provide(flags.Provide))

			app := fx.New(fxOptions(opts)...)

			if err := app.Start(ctx); err != nil {
				return prefix(name, err)
			}

			return prefix(name, app.Stop(ctx))
		},
	}

	c.cmds = append(c.cmds, cmd)

	return flags
}

// Run the command.
func (c *Command) Run(ctx context.Context, args ...string) error {
	if len(args) == 0 {
		args = os.SanitizeArgs(os.Args)
	}

	name := c.name.String()
	runner := acmd.RunnerOf(c.cmds, acmd.Config{
		AppName:        name,
		AppDescription: name,
		Version:        c.version.String(),
		Args:           args,
		Context:        ctx,
	})

	return runner.Run()
}

// ExitOnError will run the command and exit on error.
func (c *Command) ExitOnError(ctx context.Context, args ...string) {
	if err := c.Run(ctx, args...); err != nil {
		slog.Error("could not start", logger.Error(err))
		c.exit(1)
	}
}

func fxOptions(options []fx.Option) []fx.Option {
	return append(options, fx.StartTimeout(time.Minute), fx.StopTimeout(time.Minute), fx.NopLogger)
}

func prefix(prefix string, err error) error {
	return se.Prefix(prefix+": failed to run", dig.RootCause(err))
}
