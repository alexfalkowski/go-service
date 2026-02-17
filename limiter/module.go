package limiter

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires the default KeyMap into Fx.
var Module = di.Module(
	di.Constructor(NewKeyMap),
)
