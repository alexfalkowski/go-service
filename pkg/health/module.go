package health

import (
	"github.com/alexfalkowski/go-service/pkg/health/transport/grpc"
	"github.com/alexfalkowski/go-service/pkg/health/transport/http"
	"go.uber.org/fx"
)

var (
	// Module for fx.
	Module = fx.Options(fx.Invoke(http.Register), fx.Invoke(grpc.Register), fx.Provide(NewServer))
)
