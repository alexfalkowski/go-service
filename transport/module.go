package transport

import (
	"github.com/alexfalkowski/go-service/transport/grpc"
	gprometheus "github.com/alexfalkowski/go-service/transport/grpc/metrics/prometheus"
	gopentracing "github.com/alexfalkowski/go-service/transport/grpc/trace/opentracing"
	"github.com/alexfalkowski/go-service/transport/http"
	hprometheus "github.com/alexfalkowski/go-service/transport/http/metrics/prometheus"
	hopentracing "github.com/alexfalkowski/go-service/transport/http/trace/opentracing"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	nprometheus "github.com/alexfalkowski/go-service/transport/nsq/metrics/prometheus"
	nopentracing "github.com/alexfalkowski/go-service/transport/nsq/trace/opentracing"
	"go.uber.org/fx"
)

var (
	// GRPCModule for fx.
	GRPCModule = fx.Options(
		fx.Provide(grpc.NewServer),
		fx.Provide(grpc.UnaryServerInterceptor),
		fx.Provide(grpc.StreamServerInterceptor),
		fx.Provide(gprometheus.NewServerMetrics),
		fx.Provide(gprometheus.NewClientMetrics),
		fx.Provide(gopentracing.NewTracer),
	)

	// HTTPModule for fx.
	HTTPModule = fx.Options(
		fx.Provide(http.NewServer),
		fx.Provide(hprometheus.NewServerMetrics),
		fx.Provide(hprometheus.NewClientMetrics),
		fx.Provide(hopentracing.NewTracer),
	)

	// NSQModule for fx.
	NSQModule = fx.Options(
		fx.Provide(nprometheus.NewProducerMetrics),
		fx.Provide(nprometheus.NewConsumerMetrics),
		fx.Provide(nopentracing.NewTracer),
		fx.Provide(marshaller.NewMsgPack),
	)
)
