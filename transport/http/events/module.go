package events

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/transport/http/hooks"
)

// Module for fx.
var Module = di.Module(
	di.Constructor(NewReceiver),
	hooks.Module,
)
