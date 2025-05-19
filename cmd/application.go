package cmd

import (
	"context"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
)

// NewApplication for cmd.
func NewApplication(register func(*Command)) *Application {
	fs := os.NewFS()
	name := env.NewName(fs)
	version := env.NewVersion()
	cmd := NewCommand(name, version)

	register(cmd)

	return &Application{cmd: cmd}
}

// Application for cmd.
type Application struct {
	cmd *Command
}

// Run the application.
func (a *Application) Run(ctx context.Context, args ...string) error {
	return a.cmd.Run(ctx, args...)
}

// ExitOnError will run the application and exit on error.
func (a *Application) ExitOnError(ctx context.Context, args ...string) {
	a.cmd.ExitOnError(ctx, args...)
}
