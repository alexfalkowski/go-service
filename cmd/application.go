package cmd

import (
	"context"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
)

// RegisterFunc for cmd.
type RegisterFunc = func(command *Command)

// ApplicationOption for cmd.
type ApplicationOption interface {
	apply(opts *applicationOpts)
}

type applicationOpts struct {
	exit    os.ExitFunc
	name    env.Name
	version env.Version
}

type applicationOptionFunc func(*applicationOpts)

func (f applicationOptionFunc) apply(o *applicationOpts) {
	f(o)
}

// WithApplicationName for cmd.
func WithApplicationName(name env.Name) ApplicationOption {
	return applicationOptionFunc(func(o *applicationOpts) {
		o.name = name
	})
}

// WithApplicationVersion for cmd.
func WithApplicationVersion(version env.Version) ApplicationOption {
	return applicationOptionFunc(func(o *applicationOpts) {
		o.version = version
	})
}

// WithApplicationExit for cmd.
func WithApplicationExit(exit os.ExitFunc) ApplicationOption {
	return applicationOptionFunc(func(o *applicationOpts) {
		o.exit = exit
	})
}

// NewApplication for cmd.
func NewApplication(register RegisterFunc, opts ...ApplicationOption) *Application {
	ops := applicationOptions(opts...)

	cmd := NewCommand(ops.name, ops.version, ops.exit)
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

func applicationOptions(opts ...ApplicationOption) *applicationOpts {
	ops := &applicationOpts{}
	for _, o := range opts {
		o.apply(ops)
	}

	if !ops.name.IsSet() {
		ops.name = env.NewName(os.NewFS())
	}

	if !ops.version.IsSet() {
		ops.version = env.NewVersion()
	}

	if ops.exit == nil {
		ops.exit = os.NewExitFunc()
	}

	return ops
}
