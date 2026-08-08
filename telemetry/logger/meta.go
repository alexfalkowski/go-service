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
// Metadata values are bounded by limit at a valid UTF-8 boundary.
//
// The returned attributes are intended to be appended to log records to provide
// consistent request/service context across log lines.
func Meta(ctx context.Context, limit meta.Limit) []slog.Attr {
	attrs := meta.Attributes(ctx, limit)
	fields := make([]slog.Attr, 0, len(attrs))
	for key, value := range attrs {
		fields = append(fields, slog.String(key, value))
	}

	return fields
}
