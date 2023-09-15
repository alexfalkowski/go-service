package transport

import (
	"github.com/alexfalkowski/go-service/transport/grpc"
	gtel "github.com/alexfalkowski/go-service/transport/grpc/telemetry"
	"github.com/alexfalkowski/go-service/transport/http"
	htel "github.com/alexfalkowski/go-service/transport/http/telemetry"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	ntel "github.com/alexfalkowski/go-service/transport/nsq/telemetry"
	"go.uber.org/fx"
)

var (
	// GRPCModule for fx.
	GRPCModule = fx.Options(
		fx.Provide(grpc.NewServer),
		fx.Provide(grpc.UnaryServerInterceptor),
		fx.Provide(grpc.StreamServerInterceptor),
		fx.Provide(gtel.NewServerMetrics),
		fx.Provide(gtel.NewClientMetrics),
		fx.Provide(gtel.NewTracer),
	)

	// HTTPModule for fx.
	HTTPModule = fx.Options(
		fx.Provide(http.NewServer),
		fx.Provide(htel.NewServerMetrics),
		fx.Provide(htel.NewClientMetrics),
		fx.Provide(htel.NewTracer),
	)

	// NSQModule for fx.
	NSQModule = fx.Options(
		fx.Provide(ntel.NewProducerMetrics),
		fx.Provide(ntel.NewConsumerMetrics),
		fx.Provide(ntel.NewTracer),
		fx.Provide(marshaller.NewMsgPack),
	)

	// Module for fx.
	Module = fx.Options(
		GRPCModule,
		HTTPModule,
		NSQModule,
		fx.Invoke(Register),
	)
)
