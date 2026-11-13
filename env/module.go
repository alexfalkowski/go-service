package env

import "github.com/alexfalkowski/go-signal/v2"

// Module wires env-based service identity values into [github.com/alexfalkowski/go-signal/v2].
//
// It provides constructors for commonly used identity primitives derived from environment variables
// with sensible fallbacks:
//
//   - [ID] via [NewID] (SERVICE_ID or generated id)
//   - [UserAgent] via [NewUserAgent] ("<name>/<version>")
//   - [UserID] via [NewUserID] (SERVICE_USER_ID or service name)
//
// Note: this module does not provide [Name] or [Version] directly; those are commonly constructed by
// callers using [NewName] (requires *[os.FS]) and [NewVersion] (uses runtime metadata). Consumers that
// depend on [UserAgent] typically wire [Name] and [Version] elsewhere in their module graph.
var Module = signal.Module(
	signal.Constructor(NewID),
	signal.Constructor(NewUserAgent),
	signal.Constructor(NewUserID),
)
