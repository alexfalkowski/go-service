package env

import "github.com/alexfalkowski/go-service/v2/di"

// Module provides Fx wiring for env-based identity values.
//
// The constructors in this module typically prefer environment-variable overrides
// and otherwise fall back to derived defaults (see the constructor GoDocs).
var Module = di.Module(
	di.Constructor(NewID),
	di.Constructor(NewUserAgent),
	di.Constructor(NewUserID),
)
