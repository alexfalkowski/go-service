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
	// ApplicationOption configures how an Application is constructed.
	//
	// Options are applied in the order provided to NewApplication. If multiple options configure the
	// same setting, the last one wins.
	ApplicationOption interface {
		apply(opts *applicationOpts)
	}

	// ExitFunc is invoked when the application decides to terminate the process.
	//
	// It is used by (*Application).ExitOnError. The default is os.Exit.
	ExitFunc = func(code int)

	applicationOpts struct {
		exiter ExitFunc
	}
)

type applicationOptionFunc func(*applicationOpts)

func (f applicationOptionFunc) apply(o *applicationOpts) {
	f(o)
}

// WithApplicationExit sets the exit function used by an Application.
//
// This is primarily useful in tests to avoid terminating the test process.
func WithApplicationExit(exiter ExitFunc) ApplicationOption {
	return applicationOptionFunc(func(o *applicationOpts) {
		o.exiter = exiter
	})
}

// RegisterFunc registers subcommands on a Commander.
//
// A RegisterFunc is invoked by NewApplication to populate the application's command set.
// Implementations typically call commander.AddServer and/or commander.AddClient.
type RegisterFunc = func(commander Commander)

// NewApplication constructs an Application and invokes register to add subcommands.
//
// The returned Application is pre-populated with module-level Name and Version derived from the environment
// (see cli.Name and cli.Version).
func NewApplication(register RegisterFunc, opts ...ApplicationOption) *Application {
	options := options(opts...)
	app := &Application{name: Name, version: Version, exiter: options.exiter}

	register(app)
	return app
}

// Application is a command-line application composed of subcommands.
//
// An Application maintains a set of commands and delegates parsing/execution to the underlying command
// framework (github.com/cristalhq/acmd).
type Application struct {
	exiter  ExitFunc
	name    env.Name
	version env.Version
	cmds    []cmd.Command
}

// AddServer adds a long-running server subcommand with DI lifecycle wiring.
//
// The returned *Command embeds a *flag.FlagSet. The flag set is parsed before DI startup and is then
// provided into the DI container so constructors can consume parsed flag values.
//
// Execution semantics:
//   - parse the command args into the command's FlagSet
//   - build a DI application with the provided options, plus the command's module and runtime.Module
//   - start the DI application
//   - block until the DI application's Done channel is closed
//   - stop the DI application
//
// Any start/stop error is wrapped with the subcommand name for easier attribution.
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

// AddClient adds a short-lived client subcommand with DI lifecycle wiring.
//
// The returned *Command embeds a *flag.FlagSet. The flag set is parsed before DI startup and is then
// provided into the DI container so constructors can consume parsed flag values.
//
// Execution semantics:
//   - parse the command args into the command's FlagSet
//   - build a DI application with the provided options, plus the command's module
//   - start the DI application
//   - stop the DI application immediately after startup completes
//
// Any start/stop error is wrapped with the subcommand name for easier attribution.
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
//
// Run configures the underlying command runner with:
//   - app name/description derived from a.name
//   - version derived from a.version
//   - sanitized process arguments (see os.SanitizeArgs)
//   - the provided context
//
// It returns any execution error from the underlying runner or command ExecFunc.
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

// ExitOnError runs the application and terminates the process with exit status 1 if Run returns an error.
//
// The error is logged using the telemetry logger. The exit function is configurable via WithApplicationExit.
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
