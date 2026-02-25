package logger

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
)

// Meta extracts context metadata and returns it as slog attributes.
//
// It reads metadata stored in the provided context (via the `meta` package) and converts
// it to camel-cased string key/value attributes with no prefix.
//
// The returned attributes are intended to be appended to log records to provide
// consistent request/service context across log lines.
func Meta(ctx context.Context) []slog.Attr {
	strings := meta.CamelStrings(ctx, meta.NoPrefix)
	fields := make([]slog.Attr, len(strings))
	index := 0
	for k, v := range strings {
		fields[index] = slog.String(k, v)
		index++
	}
	return fields
}
