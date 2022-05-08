package transport

import (
	"github.com/alexfalkowski/go-service/transport/grpc"
	gprometheus "github.com/alexfalkowski/go-service/transport/grpc/metrics/prometheus"
	gopentracing "github.com/alexfalkowski/go-service/transport/grpc/trace/opentracing"
	"github.com/alexfalkowski/go-service/transport/http"
	hprometheus "github.com/alexfalkowski/go-service/transport/http/metrics/prometheus"
	hopentracing "github.com/alexfalkowski/go-service/transport/http/trace/opentracing"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	"go.uber.org/fx"
)

var (
	// GRPCModule for fx.
	// nolint:gochecknoglobals
	GRPCModule = fx.Options(
		fx.Provide(grpc.NewServer),
		fx.Provide(grpc.UnaryServerInterceptor),
		fx.Provide(grpc.StreamServerInterceptor),
		fx.Provide(gprometheus.NewServerMetrics),
		fx.Provide(gprometheus.NewClientMetrics),
	)

	// GRPCOpentracingModule for fx.
	// nolint:gochecknoglobals
	GRPCOpentracingModule = fx.Provide(gopentracing.NewTracer)

	// HTTPModule for fx.
	// nolint:gochecknoglobals
	HTTPModule = fx.Options(
		fx.Provide(http.NewServer),
		fx.Provide(hprometheus.NewServerMetrics),
	)

	// HTTPOpentracingModule for fx.
	// nolint:gochecknoglobals
	HTTPOpentracingModule = fx.Provide(hopentracing.NewTracer)

	// NSQMsgPackMarshallerModule for fx.
	// nolint:gochecknoglobals
	NSQMsgPackMarshallerModule = fx.Provide(marshaller.NewMsgPack)
)
