package logger

import (
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	"go.uber.org/fx"
)

var (
	// ZapLogger for fx.
	ZapLogger = fx.Provide(zap.NewLogger)

	// Module for fx.
	Module = fx.Options(ZapLogger)
)
