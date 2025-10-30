package meta

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/iancoleman/strcase"
)

const (
	// NoPrefix for strings.
	NoPrefix = strings.Empty

	meta = context.Key("meta")
)

// WithAttribute to meta.
func WithAttribute(ctx context.Context, key string, value Value) context.Context {
	return context.WithValue(ctx, meta, attributes(ctx).Add(key, value))
}

// Attribute of meta.
func Attribute(ctx context.Context, key string) Value {
	return attributes(ctx).Get(key)
}

// Map is a key-value map.
type Map map[string]string

// SnakeStrings for meta.
func SnakeStrings(ctx context.Context, prefix string) Map {
	return attributes(ctx).Strings(prefix, strcase.ToSnake)
}

// CamelStrings for meta.
func CamelStrings(ctx context.Context, prefix string) Map {
	return attributes(ctx).Strings(prefix, strcase.ToLowerCamel)
}

// Strings for meta.
func Strings(ctx context.Context, prefix string) Map {
	return attributes(ctx).Strings(prefix, func(s string) string { return s })
}
