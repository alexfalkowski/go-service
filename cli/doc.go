// Package cli provides helpers for building command-line applications with go-service.
//
// This package wraps the command framework used by this module ([github.com/cristalhq/acmd]) and provides
// a small layer of conveniences for wiring subcommands using go-service DI (Fx/Dig via the `di` package).
//
// In typical service applications, this package is used together with `go-service-template` and the
// high-level module bundles from the `module` package.
//
// # Entry points
//
// Start with [NewApplication] to construct an [Application] and register subcommands via a [RegisterFunc].
// Most applications then call either:
//
//   - [Application.Run] to execute the CLI and return any error, or
//   - [Application.RunCode] to execute the CLI and return a process exit code.
//
// [Application.RunCode] returns [os.ExitCodeSuccess] on success, returns the requested
// non-zero shutdown exit code when the DI application shuts down with one, and
// returns [os.ExitCodeFailure] for other errors.
//
// [Application.Run] reads the go-service [os.Args] variable, sanitizes injected
// Go test harness flags such as `-test.v` and `-test.run=...`, and then passes
// the remaining arguments to the command runner.
//
// # Subcommands and DI wiring
//
// Subcommands are added via [Commander] methods:
//
//   - [Application.AddServer] creates a long-running server-style command. It starts the DI app and then
//     blocks until the DI app signals completion, stopping it afterwards.
//   - [Application.AddClient] creates a short-lived client-style command. It starts the DI app, then stops it
//     immediately after startup completes.
//
// Each added subcommand returns a *[Command], which embeds a `*flag.FlagSet`.
// Define command-specific flags on that `FlagSet` before execution. The
// command implementation parses it and wires the parsed flag set into DI so
// constructors can read parsed values. Command names must be unique across the
// application.
//
// # Environment-derived metadata
//
// The CLI name and version are derived from environment helpers in `env` and exposed via package variables
// ([Name] and [Version]) so they can be reused consistently.
package cli
