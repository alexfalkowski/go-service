package hooks

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires webhook secret generation and hook construction into Fx.
var Module = di.Module(
	di.Constructor(NewGenerator),
	di.Constructor(NewHook),
)
