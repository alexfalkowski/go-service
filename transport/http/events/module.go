package events

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/transport/http/hooks"
)

// Module wires HTTP event receiver and hooks into [go.uber.org/fx].
var Module = di.Module(
	di.Constructor(NewReceiver),
	hooks.Module,
)
