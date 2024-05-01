package events

import (
	"github.com/alexfalkowski/go-service/transport/events/http"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	http.Module,
)
