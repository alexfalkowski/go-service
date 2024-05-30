package cmd

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/os"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"go.uber.org/fx"
)

// Command for application.
type Command struct {
	root *cobra.Command
}

// New command.
func New(version string) *Command {
	name := os.ExecutableName()

	root := &cobra.Command{
		Use:          name,
		Short:        name,
		Long:         name,
		SilenceUsage: true,
		Version:      version,
	}

	root.SetErrPrefix(name + ":")

	return &Command{root: root}
}

// Root command.
func (c *Command) Root() *cobra.Command {
	return c.root
}

// AddServer to the command.
func (c *Command) AddServer(opts ...fx.Option) *cobra.Command {
	return c.AddServerCommand("server", "Start the server.", opts...)
}

// AddClient to the command.
func (c *Command) AddClient(opts ...fx.Option) *cobra.Command {
	return c.AddClientCommand("client", "Start the client.", opts...)
}

// AddServerCommand to root.
func (c *Command) AddServerCommand(name, description string, opts ...fx.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:          name,
		Short:        description,
		Long:         description,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, _ []string) error {
			return RunServer(c.Context(), name, opts...)
		},
	}

	c.root.AddCommand(cmd)

	return cmd
}

// AddClientCommand to root.
func (c *Command) AddClientCommand(name, description string, opts ...fx.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:          name,
		Short:        description,
		Long:         description,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, _ []string) error {
			return RunClient(c.Context(), name, opts...)
		},
	}

	c.root.AddCommand(cmd)

	return cmd
}

// Run the command with a an arg.
func (c *Command) RunWithArgs(args []string) error {
	c.root.SetArgs(args)

	return c.root.Execute()
}

// Run the command.
func (c *Command) Run() error {
	return c.root.Execute()
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
