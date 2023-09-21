package transport

import (
	"github.com/alexfalkowski/go-service/transport/grpc"
	gprometheus "github.com/alexfalkowski/go-service/transport/grpc/metrics/prometheus"
	gotel "github.com/alexfalkowski/go-service/transport/grpc/otel"
	"github.com/alexfalkowski/go-service/transport/http"
	hprometheus "github.com/alexfalkowski/go-service/transport/http/metrics/prometheus"
	hotel "github.com/alexfalkowski/go-service/transport/http/otel"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	nprometheus "github.com/alexfalkowski/go-service/transport/nsq/metrics/prometheus"
	notel "github.com/alexfalkowski/go-service/transport/nsq/otel"
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
		fx.Provide(gotel.NewTracer),
	)

	// HTTPModule for fx.
	HTTPModule = fx.Options(
		fx.Provide(http.NewServer),
		fx.Provide(hprometheus.NewServerCollector),
		fx.Provide(hprometheus.NewClientCollector),
		fx.Provide(hotel.NewTracer),
	)

	// NSQModule for fx.
	NSQModule = fx.Options(
		fx.Provide(nprometheus.NewProducerCollector),
		fx.Provide(nprometheus.NewConsumerCollector),
		fx.Provide(notel.NewTracer),
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
