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

// WithAttribute returns a copy of ctx with the given attribute stored under key.
func WithAttribute(ctx context.Context, key string, value Value) context.Context {
	return context.WithValue(ctx, meta, attributes(ctx).Add(key, value))
}

// Attribute returns the stored attribute value for key.
func Attribute(ctx context.Context, key string) Value {
	return attributes(ctx).Get(key)
}

// Map is a key-value map.
type Map map[string]string

// SnakeStrings returns all stored attributes as a string map with snake_cased keys.
//
// prefix is prepended to each exported key.
func SnakeStrings(ctx context.Context, prefix string) Map {
	return attributes(ctx).Strings(prefix, strcase.ToSnake)
}

// CamelStrings returns all stored attributes as a string map with lowerCamelCased keys.
//
// prefix is prepended to each exported key.
func CamelStrings(ctx context.Context, prefix string) Map {
	return attributes(ctx).Strings(prefix, strcase.ToLowerCamel)
}

// Strings returns all stored attributes as a string map with keys unchanged.
//
// prefix is prepended to each exported key.
func Strings(ctx context.Context, prefix string) Map {
	return attributes(ctx).Strings(prefix, func(s string) string { return s })
}
