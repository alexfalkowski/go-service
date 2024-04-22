package zap

import (
	"context"

	"github.com/alexfalkowski/go-service/meta"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Meta for zap.
func Meta(ctx context.Context) []zapcore.Field {
	fields := []zapcore.Field{}

	for k, v := range meta.CamelStrings(ctx, "") {
		fields = append(fields, zap.String(k, v))
	}

	return fields
}
