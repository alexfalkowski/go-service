package logger

import (
	"log/slog"
	"unicode/utf8"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
)

const maxMetaValueLength = 1024

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
		fields[index] = slog.String(k, truncateMetaValue(v))
		index++
	}
	return fields[:index]
}

func truncateMetaValue(value string) string {
	if len(value) <= maxMetaValueLength {
		return value
	}

	for index := range value {
		if index > maxMetaValueLength {
			_, size := utf8.DecodeLastRuneInString(value[:index])
			return value[:index-size]
		}
	}

	return value[:maxMetaValueLength]
}
