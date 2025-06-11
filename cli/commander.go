package cli

import "github.com/alexfalkowski/go-service/v2/di"

// Option is an alias of di.Option.
type Option = di.Option

// Commander allows adding different sub commands.
type Commander interface {
	// AddServer sub command.
	AddServer(name, description string, opts ...Option) Command

	// AddClient sub command.
	AddClient(name, description string, opts ...Option) Command
}
