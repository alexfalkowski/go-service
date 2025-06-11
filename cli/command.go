package cli

import (
	"context"
	"slices"

	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/flag"
	"github.com/alexfalkowski/go-service/v2/runtime"
)

// Command provides a flag set and a way to execute the command.
type Command interface {
	// FlagSet of the command.
	FlagSet() *flag.FlagSet

	// AddInput to the flag set.
	AddInput(value string)

	// Exec the command.
	Exec(ctx context.Context, args []string) error
}

// NewServerCommand creates a command for long running server processes.
func NewServerCommand(name string, opts ...Option) *ServerCommand {
	return &ServerCommand{
		flags: flag.NewFlagSet(name),
		opts:  opts,
	}
}

// ServerCommand for cli.
type ServerCommand struct {
	flags *flag.FlagSet
	opts  []di.Option
}

// FlagSet of the command.
func (c *ServerCommand) FlagSet() *flag.FlagSet {
	return c.flags
}

// AddInput to the flag set.
func (c *ServerCommand) AddInput(value string) {
	c.flags.AddInput(value)
}

// Exec a server command with the provided context and arguments.
func (c *ServerCommand) Exec(ctx context.Context, args []string) error {
	if err := c.flags.Parse(args); err != nil {
		return err
	}

	opts := append(slices.Clone(c.opts), di.Constructor(Provide), di.Constructor(c.flags.Provide), runtime.Module, di.NoLogger)
	app := di.New(opts...)
	done := app.Done()

	if err := app.Start(ctx); err != nil {
		return err
	}

	<-done

	return app.Stop(ctx)
}

// NewServerCommand creates a command for short running child processes.
func NewClientCommand(name string, opts ...Option) *ClientCommand {
	return &ClientCommand{
		flags: flag.NewFlagSet(name),
		opts:  opts,
	}
}

// ClientCommand for cli.
type ClientCommand struct {
	flags *flag.FlagSet
	opts  []di.Option
}

// FlagSet of the command.
func (c *ClientCommand) FlagSet() *flag.FlagSet {
	return c.flags
}

// AddInput to the flag set.
func (c *ClientCommand) AddInput(value string) {
	c.flags.AddInput(value)
}

// Exec a client command with the provided context and arguments.
func (c *ClientCommand) Exec(ctx context.Context, args []string) error {
	if err := c.flags.Parse(args); err != nil {
		return err
	}

	opts := append(slices.Clone(c.opts), di.Constructor(Provide), di.Constructor(c.flags.Provide), di.NoLogger)
	app := di.New(opts...)

	if err := app.Start(ctx); err != nil {
		return err
	}

	return app.Stop(ctx)
}
