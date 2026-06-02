package hooks

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/hooks"
)

// Module wires HTTP webhook handler helpers into [go.uber.org/fx].
var Module = di.Module(
	hooks.Module,
	di.Constructor(NewWebhook),
)
