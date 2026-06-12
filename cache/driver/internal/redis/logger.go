package redis

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
)

type redisLogger struct {
	logger *logger.Logger
}

func (l redisLogger) Printf(ctx context.Context, format string, args ...any) {
	l.logger.LogAttrs(ctx, logger.LevelWarn, logger.NewText("redis: "+fmt.Sprintf(format, args...)))
}
