package transport

import (
	"github.com/alexfalkowski/go-service/transport/grpc"
	gprometheus "github.com/alexfalkowski/go-service/transport/grpc/telemetry/metrics/prometheus"
	gtracer "github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport/http"
	hprometheus "github.com/alexfalkowski/go-service/transport/http/telemetry/metrics/prometheus"
	htracer "github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	nprometheus "github.com/alexfalkowski/go-service/transport/nsq/telemetry/metrics/prometheus"
	ntracer "github.com/alexfalkowski/go-service/transport/nsq/telemetry/tracer"
	"go.uber.org/fx"
)

var (
	// GRPCModule for fx.
	GRPCModule = fx.Options(
		fx.Provide(grpc.NewServer),
		fx.Provide(grpc.UnaryServerInterceptor),
		fx.Provide(grpc.StreamServerInterceptor),
		fx.Provide(gprometheus.NewServerCollector),
		fx.Provide(gprometheus.NewClientCollector),
		fx.Provide(gtracer.NewTracer),
	)

	// HTTPModule for fx.
	HTTPModule = fx.Options(
		fx.Provide(http.NewServer),
		fx.Provide(hprometheus.NewServerCollector),
		fx.Provide(hprometheus.NewClientCollector),
		fx.Provide(htracer.NewTracer),
	)

	// NSQModule for fx.
	NSQModule = fx.Options(
		fx.Provide(nprometheus.NewProducerCollector),
		fx.Provide(nprometheus.NewConsumerCollector),
		fx.Provide(ntracer.NewTracer),
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
