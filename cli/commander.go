package cli

import "github.com/alexfalkowski/go-service/v2/di"

type (
	// Option is an alias for di.Option.
	Option = di.Option

	// Commander allows registering CLI subcommands.
	Commander interface {
		// AddServer adds a server subcommand.
		AddServer(name, description string, opts ...Option) *Command

		// AddClient adds a client subcommand.
		AddClient(name, description string, opts ...Option) *Command
	}
)
