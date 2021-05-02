package health

import (
	"github.com/alexfalkowski/go-service/pkg/health/grpc"
	"github.com/alexfalkowski/go-service/pkg/health/http"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(fx.Invoke(http.Register), fx.Invoke(grpc.Register), fx.Provide(NewServer))
