package hooks

import (
	"github.com/alexfalkowski/go-service/v3/hooks"
	"github.com/alexfalkowski/go-signal/v2"
)

// Module wires HTTP webhook handler helpers into [github.com/alexfalkowski/go-signal/v2].
var Module = signal.Module(
	hooks.Module,
	signal.Constructor(NewWebhook),
)
