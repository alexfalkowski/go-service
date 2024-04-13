package transport

import (
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
	"go.uber.org/fx"
)

var (
	// GRPCModule for fx.
	GRPCModule = fx.Options(
		fx.Provide(grpc.NewServer),
	)

	// HTTPModule for fx.
	HTTPModule = fx.Options(
		fx.Provide(http.NewServer),
	)

	// Module for fx.
	Module = fx.Options(
		GRPCModule,
		HTTPModule,
		fx.Provide(NewServers),
		fx.Invoke(Register),
	)
)
