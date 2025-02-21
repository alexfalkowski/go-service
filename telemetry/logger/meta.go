package logger

import (
	"context"
	"log/slog"

	"github.com/alexfalkowski/go-service/meta"
)

// Meta for logger.
func Meta(ctx context.Context) []slog.Attr {
	strings := meta.CamelStrings(ctx, "")
	fields := make([]slog.Attr, len(strings))
	index := 0

	for k, v := range strings {
		fields[index] = slog.String(k, v)
		index++
	}

	return fields
}
