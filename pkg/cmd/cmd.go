package cmd

import (
	"context"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// New command with serve and worker.
func New(timeout time.Duration, serverOpts []fx.Option, workerOpts []fx.Option) (*cobra.Command, error) {
	cfg, err := NewConfig()
	if err != nil {
		return nil, err
	}

	rootCmd := &cobra.Command{
		Use:          strings.ToLower(cfg.Name),
		Short:        cfg.Description,
		Long:         cfg.Description,
		SilenceUsage: true,
	}

	serveCmd := &cobra.Command{
		Use:          "serve",
		Short:        "Serve the API.",
		Long:         "Serve the API.",
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			return RunServer(args, timeout, serverOpts)
		},
	}

	rootCmd.AddCommand(serveCmd)

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
