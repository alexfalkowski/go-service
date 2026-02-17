package logger

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires logger construction into Fx.
var Module = di.Module(
	di.Constructor(NewLogger),
	di.Constructor(convertLogger),
)
