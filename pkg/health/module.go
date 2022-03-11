package health

import (
	"github.com/alexfalkowski/go-service/pkg/health/transport/grpc"
	"github.com/alexfalkowski/go-service/pkg/health/transport/http"
	"go.uber.org/fx"
)

var (
	// HTTPModule for fx.
	// nolint:gochecknoglobals
	HTTPModule = fx.Options(fx.Invoke(http.Register))

	// GRPCModule for fx.
	// nolint:gochecknoglobals
	GRPCModule = fx.Options(fx.Invoke(grpc.Register))

	// ServerModule for fx.
	// nolint:gochecknoglobals
	ServerModule = fx.Options(fx.Provide(NewServer))
)
