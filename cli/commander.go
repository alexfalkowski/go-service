package cli

import "github.com/alexfalkowski/go-service/v2/di"

type (
	// Option is an alias for di.Option.
	Option = di.Option

	// Commander registers CLI subcommands on an application.
	//
	// Implementations typically add subcommands that build and run a DI application using go-service's
	// `di` package. The returned `*Command` embeds a `*flag.FlagSet` so you can define command-specific
	// flags before execution.
	Commander interface {
		// AddServer registers a long-running server-style subcommand.
		//
		// The subcommand:
		//   - parses command args into the returned `*Command`'s FlagSet,
		//   - starts a DI application built from opts plus server-specific wiring,
		//   - blocks until the DI application signals completion,
		//   - then stops the DI application.
		AddServer(name, description string, opts ...Option) *Command

		// AddClient registers a short-lived client-style subcommand.
		//
		// The subcommand:
		//   - parses command args into the returned `*Command`'s FlagSet,
		//   - starts a DI application built from opts plus client-specific wiring,
		//   - then stops the DI application immediately after startup completes.
		AddClient(name, description string, opts ...Option) *Command
	}
)
