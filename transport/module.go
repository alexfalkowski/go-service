package transport

import (
	"github.com/alexfalkowski/go-service/transport/grpc"
	gtracer "github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport/http"
	htracer "github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	"go.uber.org/fx"
)

var (
	// GRPCModule for fx.
	GRPCModule = fx.Options(
		fx.Provide(grpc.NewServer),
		fx.Provide(gtracer.NewTracer),
	)

	// HTTPModule for fx.
	HTTPModule = fx.Options(
		fx.Provide(http.NewServer),
		fx.Provide(htracer.NewTracer),
	)

	// Module for fx.
	Module = fx.Options(
		GRPCModule,
		HTTPModule,
		fx.Provide(NewServers),
		fx.Invoke(Register),
	)
)
