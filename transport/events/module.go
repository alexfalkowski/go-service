package events

import (
	"github.com/alexfalkowski/go-service/transport/events/http"
	"github.com/alexfalkowski/go-service/transport/http/hooks"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	http.Module,
	hooks.Module,
)
