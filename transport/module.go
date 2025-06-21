package transport

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/events"
)

// Module for fx.
var Module = di.Module(
	grpc.Module,
	http.Module,
	events.Module,
	di.Constructor(NewServers),
	di.Register(Register),
)
