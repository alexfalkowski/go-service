package logger

import "github.com/alexfalkowski/go-service/v2/di"

// Module for fx.
var Module = di.Module(
	di.Constructor(NewLogger),
	di.Constructor(newLogger),
)
