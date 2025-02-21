package logger

import (
	"context"

	"github.com/alexfalkowski/go-service/meta"
)

// Meta for logger.
func Meta(ctx context.Context) []Field {
	strings := meta.CamelStrings(ctx, "")
	fields := make([]Field, len(strings))
	index := 0

	for k, v := range strings {
		fields[index] = String(k, v)
		index++
	}

	return fields
}
