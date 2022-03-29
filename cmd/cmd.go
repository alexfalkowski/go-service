package cmd

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/os"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// New command with server and worker.
func New(timeout time.Duration, serverOpts []fx.Option, workerOpts []fx.Option) (*cobra.Command, error) {
	name, err := os.ExecutableName()
	if err != nil {
		return nil, err
	}

	rootCmd := &cobra.Command{
		Use:          name,
		Short:        name,
		Long:         name,
		SilenceUsage: true,
	}

	serverCmd := &cobra.Command{
		Use:          "server",
		Short:        "Start the server.",
		Long:         "Start the server.",
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			return RunServer(args, timeout, serverOpts)
		},
	}

	rootCmd.AddCommand(serverCmd)

	workerCmd := &cobra.Command{
		Use:          "worker",
		Short:        "Start the worker.",
		Long:         "Start the worker.",
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			return RunServer(args, timeout, workerOpts)
		},
	}

	rootCmd.AddCommand(workerCmd)

	return rootCmd, nil
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
