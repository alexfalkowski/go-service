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
func (c *Command) AddServer(opts []fx.Option) {
	server := &cobra.Command{
		Use:          "server",
		Short:        "Start the server.",
		Long:         "Start the server.",
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			return RunServer(args, c.timeout, opts)
		},
	}

	c.root.AddCommand(server)
}

// AddWorker to the command.
func (c *Command) AddWorker(opts []fx.Option) {
	worker := &cobra.Command{
		Use:          "worker",
		Short:        "Start the worker.",
		Long:         "Start the worker.",
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			return RunServer(args, c.timeout, opts)
		},
	}

	c.root.AddCommand(worker)
}

// AddClient to the command.
func (c *Command) AddClient(opts []fx.Option) {
	worker := &cobra.Command{
		Use:          "client",
		Short:        "Start the client.",
		Long:         "Start the client.",
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			return RunClient(args, c.timeout, opts)
		},
	}

	c.root.AddCommand(worker)
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
