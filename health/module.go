package health

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires the health server into [go.uber.org/fx]/[go.uber.org/dig].
//
// It provides a constructor for `*server.Server` (from [github.com/alexfalkowski/go-health/v2/server])
// via [NewServer], which is started/stopped automatically using Fx lifecycle hooks.
//
// Module does not install HTTP or gRPC health endpoints by itself. Those endpoint registrations live in
// transport-specific health modules, while checker packages provide registrations that can be attached
// to the shared server.
var Module = di.Module(
	di.Constructor(NewServer),
)
