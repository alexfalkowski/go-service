package errors

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires telemetry error helpers into Fx.
var Module = di.Module(
	di.Constructor(NewHandler),
	di.Register(Register),
)
