package health

import "github.com/alexfalkowski/go-signal/v2"

// Module wires the health server into [github.com/alexfalkowski/go-signal/v2].
//
// It provides a constructor for `*server.Server` (from [github.com/alexfalkowski/go-health/v2/server])
// via [NewServer], which is started/stopped automatically using go-signal lifecycle hooks.
//
// Module does not install HTTP or gRPC health endpoints by itself. Those endpoint registrations live in
// transport-specific health modules, while checker packages provide registrations that can be attached
// to the shared server.
var Module = signal.Module(
	signal.Constructor(NewServer),
)
