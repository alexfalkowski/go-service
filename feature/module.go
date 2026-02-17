package feature

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires OpenFeature client and registration into Fx.
var Module = di.Module(
	di.Constructor(NewClient),
	di.Register(Register),
)
