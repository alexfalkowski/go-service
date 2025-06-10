package hooks

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/hooks"
)

// Module for fx.
var Module = di.Module(
	hooks.Module,
	di.Constructor(NewWebhook),
)
