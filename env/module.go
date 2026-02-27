package env

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires env-based service identity values into Fx/Dig.
//
// It provides constructors for commonly used identity primitives derived from environment variables
// with sensible fallbacks:
//
//   - ID via NewID (SERVICE_ID or generated id)
//   - UserAgent via NewUserAgent ("<name>/<version>")
//   - UserID via NewUserID (SERVICE_USER_ID or service name)
//
// Note: this module does not provide Name or Version directly; those are commonly constructed by
// callers using NewName (requires *os.FS) and NewVersion (uses runtime metadata). Consumers that
// depend on UserAgent typically wire Name and Version elsewhere in their module graph.
var Module = di.Module(
	di.Constructor(NewID),
	di.Constructor(NewUserAgent),
	di.Constructor(NewUserID),
)
