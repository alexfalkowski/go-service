package health

import (
	"github.com/alexfalkowski/go-service/health/transport/grpc"
	"github.com/alexfalkowski/go-service/health/transport/http"
	"go.uber.org/fx"
)

var (
	// HTTPModule for fx.
	HTTPModule = fx.Options(fx.Invoke(http.Register))

	// GRPCModule for fx.
	GRPCModule = fx.Options(fx.Invoke(grpc.Register))

	// ServerModule for fx.

	ServerModule = fx.Options(fx.Provide(NewServer))
)
