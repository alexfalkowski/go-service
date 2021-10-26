package transport

import (
	"github.com/alexfalkowski/go-service/pkg/transport/grpc"
	"github.com/alexfalkowski/go-service/pkg/transport/http"
	"go.uber.org/fx"
)

var (
	// GRPCServerModule for fx.
	GRPCServerModule = fx.Options(fx.Provide(grpc.NewServer), fx.Provide(grpc.UnaryServerInterceptor), fx.Provide(grpc.StreamServerInterceptor))

	// HTTPServerModule for fx.
	HTTPServerModule = fx.Options(fx.Provide(http.NewServer))
)
