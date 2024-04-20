package transport

import (
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	grpc.Module,
	http.Module,
	fx.Provide(NewServers),
	fx.Invoke(Register),
)
