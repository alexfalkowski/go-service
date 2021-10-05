package cmd

import (
	"time"

	"go.uber.org/fx"
)

// Config for cmd.
type Config struct {
	Name        string
	Description string
	Timeout     time.Duration
	ServerOpts  []fx.Option
	WorkerOpts  []fx.Option
}
