package events

import (
	"github.com/alexfalkowski/go-service/v3/transport/http/hooks"
	"github.com/alexfalkowski/go-signal/v2"
)

// Module wires HTTP event receiver and hooks into [github.com/alexfalkowski/go-signal/v2].
var Module = signal.Module(
	signal.Constructor(NewReceiver),
	hooks.Module,
)
