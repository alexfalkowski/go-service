package logger

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Meta for logger.
func Meta(ctx context.Context) []slog.Attr {
	strings := meta.CamelStrings(ctx, strings.Empty)
	fields := make([]slog.Attr, len(strings))
	index := 0
	for k, v := range strings {
		fields[index] = slog.String(k, v)
		index++
	}
	return fields
}
