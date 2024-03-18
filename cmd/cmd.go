package cmd

import (
	"context"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/time"
	"github.com/spf13/cobra"
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

	root.PersistentFlags().StringVarP(&inputFlag, "input", "i", "env:CONFIG_FILE", "input config location (format kind:location, default env:CONFIG_FILE)")

	return &Command{root: root}
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
		RunE: func(_ *cobra.Command, _ []string) error {
			return RunServer(opts...)
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
		RunE: func(_ *cobra.Command, _ []string) error {
			return RunClient(opts...)
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

// RunServer with args and a timeout.
func RunServer(opts ...fx.Option) error {
	app := fx.New(opts...)
	done := app.Done()

	startCtx, cancel := context.WithTimeout(context.Background(), time.Timeout)
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		return err
	}

	<-done

	stopCtx, cancel := context.WithTimeout(context.Background(), time.Timeout)
	defer cancel()

	return app.Stop(stopCtx)
}

// RunClient with args and a timeout.
func RunClient(opts ...fx.Option) error {
	app := fx.New(opts...)

	startCtx, cancel := context.WithTimeout(context.Background(), time.Timeout)
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		return err
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), time.Timeout)
	defer cancel()

	return app.Stop(stopCtx)
}
