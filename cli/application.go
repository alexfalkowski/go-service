package cli

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	cmd "github.com/cristalhq/acmd"
)

type (
	// ApplicationOption for cli.
	ApplicationOption interface {
		apply(opts *applicationOpts)
	}

	// ExitFunc for cli.
	ExitFunc = func(code int)

	applicationOpts struct {
		exiter ExitFunc
	}
)

type applicationOptionFunc func(*applicationOpts)

func (f applicationOptionFunc) apply(o *applicationOpts) {
	f(o)
}

// WithApplicationExit for cli.
func WithApplicationExit(exiter ExitFunc) ApplicationOption {
	return applicationOptionFunc(func(o *applicationOpts) {
		o.exiter = exiter
	})
}

// RegisterFunc for cli.
type RegisterFunc = func(commander Commander)

// NewApplication for cli.
func NewApplication(register RegisterFunc, opts ...ApplicationOption) *Application {
	options := options(opts...)
	app := &Application{name: Name, version: Version, exiter: options.exiter}

	register(app)
	return app
}

// Application for cli.
type Application struct {
	exiter  ExitFunc
	name    env.Name
	version env.Version
	cmds    []cmd.Command
}

// AddServer sub command.
func (a *Application) AddServer(name, description string, opts ...Option) *Command {
	server := NewCommand(name)
	cmd := cmd.Command{
		Name:        name,
		Description: description,
		ExecFunc: func(ctx context.Context, args []string) error {
			if err := server.Parse(args); err != nil {
				return err
			}

			opts = append(opts, server.module(), runtime.Module)
			app := di.New(opts...)
			done := app.Done()

			if err := app.Start(ctx); err != nil {
				return a.prefix(name, err)
			}

			<-done

			return a.prefix(name, app.Stop(ctx))
		},
	}
	a.cmds = append(a.cmds, cmd)

	return server
}

// AddClient sub command.
func (a *Application) AddClient(name, description string, opts ...Option) *Command {
	client := NewCommand(name)
	cmd := cmd.Command{
		Name:        name,
		Description: description,
		ExecFunc: func(ctx context.Context, args []string) error {
			if err := client.Parse(args); err != nil {
				return err
			}

			opts = append(opts, client.module())
			app := di.New(opts...)

			if err := app.Start(ctx); err != nil {
				return a.prefix(name, err)
			}

			return a.prefix(name, app.Stop(ctx))
		},
	}
	a.cmds = append(a.cmds, cmd)

	return client
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
		logger.LogError(ctx, "could not start", logger.Error(err))
		a.exiter(1)
	}
}

func (a *Application) prefix(prefix string, err error) error {
	return errors.Prefix(prefix+": failed to run", di.RootCause(err))
}

func options(opts ...ApplicationOption) *applicationOpts {
	options := &applicationOpts{}
	for _, o := range opts {
		o.apply(options)
	}
	if options.exiter == nil {
		options.exiter = os.Exit
	}

	return options
}
