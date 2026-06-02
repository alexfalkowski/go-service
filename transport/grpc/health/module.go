package health

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires gRPC health server and registrations into [go.uber.org/fx].
var Module = di.Module(
	di.Constructor(NewServer),
	di.Register(Register),
)
