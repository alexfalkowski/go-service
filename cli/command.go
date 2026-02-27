package cli

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/flag"
)

// NewCommand creates a new CLI Command with the given name.
//
// The returned Command embeds a `*flag.FlagSet` that you can use to define CLI flags.
// Application subcommands call `(*flag.FlagSet).Parse` before starting DI, and the Command's
// `module` wires the parsed FlagSet into the DI container so constructors can read flag values.
func NewCommand(name string) *Command {
	set := flag.NewFlagSet(name)
	return &Command{FlagSet: set}
}

// Command wraps a `*flag.FlagSet` and provides DI wiring for CLI subcommands.
//
// The embedded FlagSet is intended to be configured with flags by the caller, then parsed by the
// subcommand execution path. The `module` method exposes providers for the filesystem/name/version
// metadata and the command's FlagSet.
type Command struct {
	*flag.FlagSet
}

func (c *Command) provide() *flag.FlagSet {
	return c.FlagSet
}

func (c *Command) module() di.Option {
	return di.Module(
		di.Constructor(provide),
		di.Constructor(c.provide),
		di.NoLogger,
	)
}
