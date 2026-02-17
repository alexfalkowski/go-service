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
	// ApplicationOption configures Application creation.
	ApplicationOption interface {
		apply(opts *applicationOpts)
	}

	// ExitFunc is invoked when the application decides to exit.
	ExitFunc = func(code int)

	applicationOpts struct {
		exiter ExitFunc
	}
)

type applicationOptionFunc func(*applicationOpts)

func (f applicationOptionFunc) apply(o *applicationOpts) {
	f(o)
}

// WithApplicationExit sets the exit function used by Application.
func WithApplicationExit(exiter ExitFunc) ApplicationOption {
	return applicationOptionFunc(func(o *applicationOpts) {
		o.exiter = exiter
	})
}

// RegisterFunc registers commands on a Commander.
type RegisterFunc = func(commander Commander)

// NewApplication constructs an Application and invokes register to add commands.
func NewApplication(register RegisterFunc, opts ...ApplicationOption) *Application {
	options := options(opts...)
	app := &Application{name: Name, version: Version, exiter: options.exiter}

	register(app)
	return app
}

// Application is a CLI application composed of subcommands.
type Application struct {
	exiter  ExitFunc
	name    env.Name
	version env.Version
	cmds    []cmd.Command
}

// AddServer adds a subcommand that runs a server with Fx lifecycle wiring.
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

// AddClient adds a subcommand that runs a client with Fx lifecycle wiring.
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

// Run executes the application using the configured command set.
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

// ExitOnError runs the application and exits with status 1 if Run returns an error.
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
