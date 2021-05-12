package logger

import (
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	"go.uber.org/fx"
)

var (
	// ZapLoggerModule for fx.
	ZapLoggerModule = fx.Options(fx.Provide(zap.NewConfig), fx.Provide(zap.NewLogger))

	// Module for fx.
	Module = fx.Options(ZapLoggerModule)
)
