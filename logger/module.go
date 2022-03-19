package logger

import (
	"github.com/alexfalkowski/go-service/logger/zap"
	"go.uber.org/fx"
)

// ZapModule for fx.
// nolint:gochecknoglobals
var ZapModule = fx.Options(fx.Provide(zap.NewConfig), fx.Provide(zap.NewLogger))
