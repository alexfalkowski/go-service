package transport

import (
	"github.com/alexfalkowski/go-service/pkg/transport/grpc"
	"github.com/alexfalkowski/go-service/pkg/transport/http"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq"
	"go.uber.org/fx"
)

var (
	// GRPCServerModule for fx.
	GRPCServerModule = fx.Options(fx.Provide(grpc.NewServer), fx.Provide(grpc.NewConfig), fx.Provide(grpc.UnaryServerInterceptor), fx.Provide(grpc.StreamServerInterceptor))

	// HTTPServerModule for fx.
	HTTPServerModule = fx.Options(fx.Invoke(http.Register), fx.Provide(http.NewMux), fx.Provide(http.NewConfig))

	// HTTPClientModule for fx.
	HTTPClientModule = fx.Options(fx.Provide(http.NewClient))

	// NSQModule for fx.
	NSQModule = fx.Options(fx.Provide(nsq.NewConfig))
)
