package logger

import (
	"log/slog"
	"unicode/utf8"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
)

// maxMetaValueLength keeps context metadata useful in logs without allowing one
// value to dominate every record.
const maxMetaValueLength = 1024

// Meta extracts context metadata and returns it as slog attributes.
//
// It reads metadata stored in the provided context (via the `meta` package) and converts
// it to camel-cased string key/value attributes with no prefix.
//
// Metadata values longer than 1024 bytes are truncated at a valid UTF-8 boundary.
//
// The returned attributes are intended to be appended to log records to provide
// consistent request/service context across log lines.
func Meta(ctx context.Context) []slog.Attr {
	metadata := meta.CamelStrings(ctx, meta.NoPrefix)
	fields := make([]slog.Attr, 0, len(metadata))
	for k, v := range metadata {
		fields = append(fields, slog.String(k, truncateMetaValue(v)))
	}
	return fields
}

func truncateMetaValue(value string) string {
	if len(value) <= maxMetaValueLength {
		return value
	}

	truncated := value[:maxMetaValueLength]
	for !utf8.ValidString(truncated) {
		_, size := utf8.DecodeLastRuneInString(truncated)
		truncated = truncated[:len(truncated)-size]
	}

	return truncated
}
