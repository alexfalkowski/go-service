package env

import "github.com/alexfalkowski/go-service/v2/di"

// Module for fx.
var Module = di.Module(
	di.Constructor(NewID),
	di.Constructor(NewUserAgent),
	di.Constructor(NewUserID),
)
