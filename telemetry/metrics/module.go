package metrics

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires metrics provider and reader into Fx.
var Module = di.Module(
	di.Constructor(NewReader),
	di.Constructor(NewMeterProvider),
	di.Constructor(NewMeter),
)
