package cli

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/flag"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	cmd "github.com/cristalhq/acmd"
	"go.uber.org/dig"
	"go.uber.org/fx"
)

// RegisterFunc for cmd.
type RegisterFunc = func(commander Commander)

// ApplicationOption for cmd.
type ApplicationOption interface {
	apply(opts *applicationOpts)
}

type applicationOpts struct {
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

// NewApplication for cmd.
func NewApplication(register RegisterFunc, opts ...ApplicationOption) *Application {
	ops := options(opts...)
	app := &Application{name: ops.name, version: ops.version}

	register(app)

	return app
}

// Application for cmd.
type Application struct {
	name    env.Name
	version env.Version
	cmds    []cmd.Command
}

// AddServer sub command.
func (a *Application) AddServer(name, description string, opts ...Option) *Command {
	flags := flag.NewFlagSet(name)
	cmd := cmd.Command{
		Name:        name,
		Description: description,
		ExecFunc: func(ctx context.Context, args []string) error {
			if err := flags.Parse(args); err != nil {
				return err
			}

			opts = append(opts, fx.Provide(flags.Provide), runtime.Module)
			app := fx.New(a.options(opts)...)
			done := app.Done()

			if err := app.Start(ctx); err != nil {
				return a.prefix(name, err)
			}

			<-done

			return a.prefix(name, app.Stop(ctx))
		},
	}

	a.cmds = append(a.cmds, cmd)

	return &Command{flags}
}

// AddClient sub command.
func (a *Application) AddClient(name, description string, opts ...Option) *Command {
	flags := flag.NewFlagSet(name)
	cmd := cmd.Command{
		Name:        name,
		Description: description,
		ExecFunc: func(ctx context.Context, args []string) error {
			if err := flags.Parse(args); err != nil {
				return err
			}

			opts = append(opts, fx.Provide(flags.Provide))
			app := fx.New(a.options(opts)...)

			if err := app.Start(ctx); err != nil {
				return a.prefix(name, err)
			}

			return a.prefix(name, app.Stop(ctx))
		},
	}

	a.cmds = append(a.cmds, cmd)

	return &Command{flags}
}

// Run the application.
func (a *Application) Run(ctx context.Context) error {
	name := a.name.String()
	runner := cmd.RunnerOf(a.cmds, cmd.Config{
		AppName:        name,
		AppDescription: name,
		Version:        a.version.String(),
		Args:           os.SanitizeArgs(os.Args),
		Context:        ctx,
	})

	return runner.Run()
}

// ExitOnError will run the application and exit on error.
func (a *Application) ExitOnError(ctx context.Context) {
	if err := a.Run(ctx); err != nil {
		logger.LogError("could not start", logger.Error(err))
		os.Exit(1)
	}
}

func (a *Application) options(options []Option) []Option {
	return append(options, fx.NopLogger)
}

func (a *Application) prefix(prefix string, err error) error {
	return errors.Prefix(prefix+": failed to run", dig.RootCause(err))
}

func options(opts ...ApplicationOption) *applicationOpts {
	ops := &applicationOpts{}
	for _, o := range opts {
		o.apply(ops)
	}

	if ops.name.IsEmpty() {
		ops.name = env.NewName(os.NewFS())
	}

	if ops.version.IsEmpty() {
		ops.version = env.NewVersion()
	}

	return ops
}
