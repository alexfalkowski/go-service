package tracer

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires tracer construction into Fx.
var Module = di.Module(
	di.Constructor(NewTracer),
)
