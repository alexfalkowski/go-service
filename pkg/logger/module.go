package logger

import (
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	"go.uber.org/fx"
)

var (
	// ZapConfig for fx.
	ZapConfig = fx.Provide(zap.NewConfig)

	// ZapLogger for fx.
	ZapLogger = fx.Provide(zap.NewLogger)

	// Module for fx.
	Module = fx.Options(ZapConfig, ZapLogger)
)
