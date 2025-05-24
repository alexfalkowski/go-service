package events

import (
	"github.com/alexfalkowski/go-service/v2/transport/http/hooks"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewReceiver),
	hooks.Module,
)
