package health

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires the health server into Fx.
var Module = di.Module(
	di.Constructor(NewServer),
)
