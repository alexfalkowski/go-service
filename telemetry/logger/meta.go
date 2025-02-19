package logger

import (
	"context"

	"github.com/alexfalkowski/go-service/meta"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Meta for logger.
func Meta(ctx context.Context) []zapcore.Field {
	strings := meta.CamelStrings(ctx, "")
	fields := make([]zapcore.Field, len(strings))
	index := 0

	for k, v := range strings {
		fields[index] = zap.String(k, v)
		index++
	}

	return fields
}
