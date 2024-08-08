package zap

import (
	"context"

	"github.com/alexfalkowski/go-service/meta"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Meta for zap.
func Meta(ctx context.Context) []zapcore.Field {
	strs := meta.CamelStrings(ctx, "")
	fields := make([]zapcore.Field, len(strs))
	cnt := 0

	for k, v := range strs {
		fields[cnt] = zap.String(k, v)
		cnt++
	}

	return fields
}
