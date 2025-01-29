package proxy

import (
	"github.com/alexfalkowski/go-service/proxy/telemetry/logger/zap"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(zap.NewLogger),
	fx.Provide(NewServer),
)
