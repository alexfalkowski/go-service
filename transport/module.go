package transport

import (
	"github.com/alexfalkowski/go-service/transport/grpc"
	gdatadog "github.com/alexfalkowski/go-service/transport/grpc/trace/opentracing/datadog"
	gjaeger "github.com/alexfalkowski/go-service/transport/grpc/trace/opentracing/jaeger"
	"github.com/alexfalkowski/go-service/transport/http"
	hdatadog "github.com/alexfalkowski/go-service/transport/http/trace/opentracing/datadog"
	hjaeger "github.com/alexfalkowski/go-service/transport/http/trace/opentracing/jaeger"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	"go.uber.org/fx"
)

var (
	// GRPCServerModule for fx.
	// nolint:gochecknoglobals
	GRPCServerModule = fx.Options(fx.Provide(grpc.NewServer), fx.Provide(grpc.UnaryServerInterceptor), fx.Provide(grpc.StreamServerInterceptor))

	// GRPCJaegerModule for fx.
	// nolint:gochecknoglobals
	GRPCJaegerModule = fx.Provide(gjaeger.NewTracer)

	// GRPCDataDogModule for fx.
	// nolint:gochecknoglobals
	GRPCDataDogModule = fx.Provide(gdatadog.NewTracer)

	// HTTPServerModule for fx.
	// nolint:gochecknoglobals
	HTTPServerModule = fx.Provide(http.NewServer)

	// GRPCJaegerModule for fx.
	// nolint:gochecknoglobals
	HTTPJaegerModule = fx.Provide(hjaeger.NewTracer)

	// GRPCDataDogModule for fx.
	// nolint:gochecknoglobals
	HTTPDataDogModule = fx.Provide(hdatadog.NewTracer)

	// NSQMsgPackMarshallerModule for fx.
	// nolint:gochecknoglobals
	NSQMsgPackMarshallerModule = fx.Provide(marshaller.NewMsgPack)
)
