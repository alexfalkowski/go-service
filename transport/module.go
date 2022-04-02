package transport

import (
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
	"go.uber.org/fx"
)

var (
	// GRPCServerModule for fx.
	// nolint:gochecknoglobals
	GRPCServerModule = fx.Options(fx.Provide(grpc.NewServer), fx.Provide(grpc.UnaryServerInterceptor), fx.Provide(grpc.StreamServerInterceptor))

	// HTTPServerModule for fx.
	// nolint:gochecknoglobals
	HTTPServerModule = fx.Options(fx.Provide(http.NewMux), fx.Provide(http.NewServer))
)
