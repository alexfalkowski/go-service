package health

import (
	"github.com/alexfalkowski/go-service/health/transport/grpc"
	"github.com/alexfalkowski/go-service/health/transport/http"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	http.Module,
	grpc.Module,
	fx.Provide(NewServer),
)
