package cli

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/flag"
)

// NewCommand creates a new command with the given name.
func NewCommand(name string) *Command {
	set := flag.NewFlagSet(name)
	return &Command{FlagSet: set}
}

// Command for cli.
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
