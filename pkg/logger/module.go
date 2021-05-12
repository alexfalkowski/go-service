package logger

import (
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	"go.uber.org/fx"
)

var (
	// ZapModule for fx.
	ZapModule = fx.Options(fx.Provide(zap.NewConfig), fx.Provide(zap.NewLogger))
)
