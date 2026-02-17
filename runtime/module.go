package runtime

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires runtime helpers into Fx.
var Module = di.Module(
	di.Register(RegisterMemLimit),
)
