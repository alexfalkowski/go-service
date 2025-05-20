package transport

import (
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/events"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/alexfalkowski/go-service/v2/transport/meta"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	grpc.Module,
	http.Module,
	events.Module,
	meta.Module,
	limiter.Module,
	fx.Provide(NewServers),
	fx.Invoke(Register),
)
