// Package cli provides helpers for building command-line applications with go-service.
//
// This package wraps the command framework used by this module (github.com/cristalhq/acmd) and provides
// a small layer of conveniences for wiring subcommands using go-service DI (Fx/Dig via the `di` package).
//
// # Entry points
//
// Start with `NewApplication` to construct an `Application` and register subcommands via a `RegisterFunc`.
// Most applications then call either:
//
//   - `(*Application).Run` to execute the CLI and return any error, or
//   - `(*Application).ExitOnError` to log failures and exit with a non-zero status code.
//
// # Subcommands and DI wiring
//
// Subcommands are added via `Commander` methods:
//
//   - `(*Application).AddServer` creates a long-running server-style command. It starts the DI app and then
//     blocks until the DI app signals completion, stopping it afterwards.
//   - `(*Application).AddClient` creates a short-lived client-style command. It starts the DI app, then stops it
//     immediately after startup completes.
//
// Each added subcommand returns a `*Command`, which embeds a `*flag.FlagSet`. You can define and parse flags on
// that `FlagSet`; the command implementation wires the flag set into DI so constructors can read parsed values.
//
// # Environment-derived metadata
//
// The CLI name and version are derived from environment helpers in `env` and exposed via package variables
// (`Name` and `Version`) so they can be reused consistently.
package cli
