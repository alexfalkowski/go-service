package cli

import "go.uber.org/fx"

// Option is an alias of fx.Option.
type Option = fx.Option

// Commander allows adding different sub commands.
type Commander interface {
	// AddServer sub command.
	AddServer(name, description string, opts ...Option) *Command

	// AddClient sub command.
	AddClient(name, description string, opts ...Option) *Command
}
