package hooks

import (
	"github.com/alexfalkowski/go-service/v2/hooks"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	hooks.Module,
	fx.Provide(NewWebhook),
)
