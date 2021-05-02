package cmd

import (
	"context"
	"time"

	"go.uber.org/fx"
)

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

	if err := app.Stop(stopCtx); err != nil {
		return err
	}

	return nil
}
