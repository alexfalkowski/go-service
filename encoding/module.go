package encoding

import "github.com/alexfalkowski/go-service/v2/di"

// Module provides the default encoder registry as a *[Map].
var Module = di.Module(
	di.Constructor(NewMap),
)
