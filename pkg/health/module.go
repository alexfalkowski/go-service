package health

import (
	"github.com/alexfalkowski/go-service/pkg/health/transport/grpc"
	"github.com/alexfalkowski/go-service/pkg/health/transport/http"
	"go.uber.org/fx"
)

var (
	// HTTPModule for fx.
	HTTPModule = fx.Invoke(http.Register)

	// GRPCModule for fx.
	GRPCModule = fx.Invoke(grpc.Register)

	// ServerModule for fx.
	ServerModule = fx.Provide(NewServer)

	// Module for fx.
	Module = fx.Options(HTTPModule, GRPCModule, ServerModule)
)
