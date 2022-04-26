package cmd

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/os"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// Command for application.
type Command struct {
	root    *cobra.Command
	timeout time.Duration
}

// New command.
func New(timeout time.Duration) *Command {
	name := os.ExecutableName()

	root := &cobra.Command{
		Use:          name,
		Short:        name,
		Long:         name,
		SilenceUsage: true,
	}

	return &Command{root: root, timeout: timeout}
}

// AddServer to the command.
func (c *Command) AddServer(opts []fx.Option) *cobra.Command {
	return c.AddServerCommand("server", "Start the server.", opts)
}

// AddWorker to the command.
func (c *Command) AddWorker(opts []fx.Option) *cobra.Command {
	return c.AddServerCommand("worker", "Start the worker.", opts)
}

// AddClient to the command.
func (c *Command) AddClient(opts []fx.Option) *cobra.Command {
	return c.AddClientCommand("client", "Start the client.", opts)
}

// AddServerCommand to root.
func (c *Command) AddServerCommand(name, description string, opts []fx.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:          name,
		Short:        description,
		Long:         description,
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			return RunServer(args, c.timeout, opts)
		},
	}

	c.root.AddCommand(cmd)

	return cmd
}

// AddClientCommand to root.
func (c *Command) AddClientCommand(name, description string, opts []fx.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:          name,
		Short:        description,
		Long:         description,
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			return RunClient(args, c.timeout, opts)
		},
	}

	c.root.AddCommand(cmd)

	return cmd
}

// Run the command with a an arg.
func (c *Command) RunWithArg(arg string) error {
	c.root.SetArgs([]string{arg})

	return c.root.Execute()
}

// Run the command.
func (c *Command) Run() error {
	return c.root.Execute()
}

// RunServer with args and a timeout.
func RunServer(args []string, timeout time.Duration, opts []fx.Option) error {
	app := fx.New(opts...)
	done := app.Done()

	startCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		return err
	}

	<-done

	stopCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return app.Stop(stopCtx)
}

// RunClient with args and a timeout.
func RunClient(args []string, timeout time.Duration, opts []fx.Option) error {
	app := fx.New(opts...)

	startCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		return err
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return app.Stop(stopCtx)
}
