package health

import "github.com/alexfalkowski/go-signal/v2"

// Module wires gRPC health server and registrations into [github.com/alexfalkowski/go-signal/v2].
var Module = signal.Module(
	signal.Constructor(NewServer),
	signal.Register(Register),
)
