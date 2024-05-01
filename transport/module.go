package transport

import (
	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/transport/events"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/transport/meta"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	grpc.Module,
	http.Module,
	events.Module,
	limiter.Module,
	meta.Module,
	fx.Provide(NewServers),
	fx.Invoke(Register),
)
