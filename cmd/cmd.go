package cmd

import (
	"context"

	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/time"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"go.uber.org/fx"
)

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

// Command for application.
type Command struct {
	root *cobra.Command
}

// Root command.
func (c *Command) Root() *cobra.Command {
	return c.root
}

// AddServer to root.
func (c *Command) AddServer(name, description string, opts ...fx.Option) *cobra.Command {
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

// AddClient to root.
func (c *Command) AddClient(name, description string, opts ...fx.Option) *cobra.Command {
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

// SetArgs will set the actual arguments that is run the root command.
func (c *Command) SetArgs(args []string) {
	c.root.SetArgs(args)
}

// Run the command.
func (c *Command) Run() error {
	return c.root.Execute()
}

// ExitOnError will run the command and exit on error.
func (c *Command) ExitOnError() {
	if err := c.Run(); err != nil {
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
