package runtime

import "github.com/alexfalkowski/go-service/v2/di"

// Module for fx.
var Module = di.Module(
	di.Register(RegisterMemLimit),
)
