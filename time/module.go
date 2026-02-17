package time

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires the network time provider into Fx.
var Module = di.Module(
	di.Constructor(NewNetwork),
)
