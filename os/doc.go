// Package os provides OS and filesystem helpers used throughout go-service.
//
// The goal of this package is to centralize a small set of common OS concerns
// (environment variables, executable/home/config directories, process exit, and
// filesystem access) behind stable APIs that are easy to use in services and
// tests.
//
// # Filesystem wrapper
//
// The primary type in this package is FS, a thin wrapper around an avfs.VFS that
// provides go-service-specific helpers:
//
//   - Path normalization:
//     FS.CleanPath expands a leading "~" (when present) and then cleans the path.
//     FS.ReadFile and FS.WriteFile call CleanPath before accessing the filesystem.
//
//   - Trimming behavior:
//     FS.ReadFile trims leading/trailing whitespace from the bytes it returns.
//     FS.WriteFile trims leading/trailing whitespace from the bytes it writes.
//     This is convenient for reading and writing configuration fragments and
//     secrets where trailing newlines are common (for example when sourced from
//     files created by tooling).
//
//   - Existence and errors:
//     FS.PathExists reports whether a path exists.
//     FS.IsNotExist reports whether an error represents a missing path.
//
//   - Path helpers:
//     FS.PathExtension returns the file extension without the leading dot.
//     FS.ExecutableName and FS.ExecutableDir derive values from the running
//     executable.
//
// NewFS constructs an FS backed by the real OS filesystem.
//
// # “Source string” pattern
//
// FS.ReadSource implements the go-service “source string” convention, which
// allows configuration to reference values that may come from different sources.
// Supported forms are:
//
//   - "env:NAME"    reads the value of environment variable NAME.
//   - "file:/path"  reads bytes from the file at /path (via FS.ReadFile,
//     including path cleaning and trimming).
//   - otherwise     treats the string as the literal value.
//
// This pattern is used in multiple subsystems (for example secret material,
// tokens, and telemetry headers) to support flexible deployment environments.
//
// # Strict helpers and panics
//
// Some OS directory helpers are intentionally strict and will panic on unexpected
// OS errors by using runtime.Must:
//
//   - Executable
//   - UserHomeDir
//   - UserConfigDir
//
// This design assumes that inability to determine these values is not a
// recoverable runtime condition for typical service operation. If your use case
// requires error handling instead of panics, call the standard library
// equivalents directly.
//
// # Relationship to the standard library
//
// Several identifiers are thin wrappers or aliases of the standard library
// package os (for example Getenv/Setenv/Unsetenv, Exit, Args, Stdout). They exist
// to keep go-service code depending on go-service packages consistently, while
// still delegating to the underlying OS implementation.
package os
