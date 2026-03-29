package cli

import (
	"maps"
	"slices"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	cmd "github.com/cristalhq/acmd"
)

// ErrCommandRegistered indicates a subcommand name has already been registered on an Application.
var ErrCommandRegistered = errors.New("command already registered")

// Exit is invoked by (*Application).ExitOnError after logging a startup failure.
//
// Tests may replace this variable to avoid terminating the test process.
var Exit = os.Exit

// RegisterFunc registers subcommands on a Commander.
//
// A RegisterFunc is invoked by NewApplication to populate the application's command set.
// Implementations typically call commander.AddServer and/or commander.AddClient.
type RegisterFunc = func(commander Commander)

// NewApplication constructs an Application and invokes register to add subcommands.
//
// The returned Application is pre-populated with module-level Name and Version derived from the environment
// (see cli.Name and cli.Version).
func NewApplication(register RegisterFunc) *Application {
	app := &Application{name: Name, version: Version}
	register(app)

	return app
}

// Application is a command-line application composed of subcommands.
//
// An Application maintains a set of commands and delegates parsing/execution to the underlying command
// framework (github.com/cristalhq/acmd).
type Application struct {
	cmds    map[string]cmd.Command
	name    env.Name
	version env.Version
}

// AddServer adds a long-running server subcommand with DI lifecycle wiring.
//
// The returned *Command embeds a *flag.FlagSet. The flag set is parsed before DI startup and is then
// provided into the DI container so constructors can consume parsed flag values. Command names must be
// unique across the application.
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
	runtime.Must(a.register(cmd))

	return server
}

// AddClient adds a short-lived client subcommand with DI lifecycle wiring.
//
// The returned *Command embeds a *flag.FlagSet. The flag set is parsed before DI startup and is then
// provided into the DI container so constructors can consume parsed flag values. Command names must be
// unique across the application.
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
	runtime.Must(a.register(cmd))

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
	runner := cmd.RunnerOf(slices.Collect(maps.Values(a.cmds)), cmd.Config{
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
// The error is logged using the telemetry logger before Exit is invoked.
func (a *Application) ExitOnError(ctx context.Context) {
	if err := a.Run(ctx); err != nil {
		logger.LogError(ctx, "could not start", logger.Error(err))
		Exit(1)
	}
}

func (a *Application) register(command cmd.Command) error {
	if a.cmds == nil {
		a.cmds = make(map[string]cmd.Command)
	}

	if _, ok := a.cmds[command.Name]; ok {
		return errors.Prefix(command.Name, ErrCommandRegistered)
	}

	a.cmds[command.Name] = command

	return nil
}

func (a *Application) prefix(prefix string, err error) error {
	return errors.Prefix(prefix+": failed to run", di.RootCause(err))
}
