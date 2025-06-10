package rand

import "github.com/alexfalkowski/go-service/v2/di"

// Module for fx.
var Module = di.Module(
	di.Constructor(NewReader),
	di.Constructor(NewGenerator),
)
