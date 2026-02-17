package health

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires HTTP health endpoints into Fx.
var Module = di.Module(
	di.Register(Register),
)
