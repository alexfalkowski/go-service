package transport

import (
	"github.com/alexfalkowski/go-service/nsq"
	"github.com/alexfalkowski/go-service/transport/grpc"
	gtracer "github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport/http"
	htracer "github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	ntracer "github.com/alexfalkowski/go-service/transport/nsq/telemetry/tracer"
	"go.uber.org/fx"
)

var (
	// GRPCModule for fx.
	GRPCModule = fx.Options(
		fx.Provide(grpc.NewServer),
		fx.Provide(grpc.UnaryServerInterceptor),
		fx.Provide(grpc.StreamServerInterceptor),
		fx.Provide(gtracer.NewTracer),
	)

	// HTTPModule for fx.
	HTTPModule = fx.Options(
		fx.Provide(http.NewServer),
		fx.Provide(htracer.NewTracer),
	)

	// NSQModule for fx.
	NSQModule = fx.Options(
		fx.Provide(ntracer.NewTracer),
		fx.Provide(nsq.NewMsgPackMarshaller),
	)

	// Module for fx.
	Module = fx.Options(
		GRPCModule,
		HTTPModule,
		NSQModule,
		fx.Invoke(Register),
	)
)
