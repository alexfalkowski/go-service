package health

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires HTTP health endpoints into [go.uber.org/fx].
var Module = di.Module(
	di.Register(Register),
)
