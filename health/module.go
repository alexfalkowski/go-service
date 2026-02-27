package health

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires the health server into Fx/Dig.
//
// It provides a constructor for `*server.Server` (from github.com/alexfalkowski/go-health/v2/server)
// via `NewServer`, which is started/stopped automatically using Fx lifecycle hooks.
//
// Additional health check registrations are typically provided by other packages (including
// `health/checker`) and installed onto the returned server elsewhere in the application graph.
var Module = di.Module(
	di.Constructor(NewServer),
)
