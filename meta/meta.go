package meta

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/iancoleman/strcase"
)

const (
	// NoPrefix is a convenience constant for passing an empty prefix to export helpers.
	NoPrefix = strings.Empty

	meta = context.Key("meta")
)

// WithAttributes returns a copy of ctx with all provided attributes stored.
//
// The storage update is copy-on-write and clones existing attributes only once before applying the
// full batch. Metadata writes use pairs so hot paths can set several attributes without repeated
// map copies and context wrappers.
func WithAttributes(ctx context.Context, pairs ...Pair) context.Context {
	if len(pairs) == 0 {
		return ctx
	}

	return context.WithValue(ctx, meta, attributes(ctx).AddPairs(pairs...))
}

// Attribute returns the stored attribute value for key.
//
// If no attribute is present for key, Attribute returns the zero-value Value. Callers can use
// Value.IsEmpty to distinguish empty values.
func Attribute(ctx context.Context, key string) Value {
	return attributes(ctx).Get(key)
}

// Map is a string key-value map of exported attributes.
//
// Export helpers return this type after rendering Values and applying key conversion and prefixing.
type Map map[string]string

// SnakeStrings returns all stored attributes as a string map with snake_cased keys.
//
// The prefix parameter is prepended to each exported key (if non-empty).
// Attributes whose rendered value is empty are skipped.
func SnakeStrings(ctx context.Context, prefix string) Map {
	return attributes(ctx).Strings(prefix, strcase.ToSnake)
}

// CamelStrings returns all stored attributes as a string map with lowerCamelCased keys.
//
// The prefix parameter is prepended to each exported key (if non-empty).
// Attributes whose rendered value is empty are skipped.
func CamelStrings(ctx context.Context, prefix string) Map {
	return attributes(ctx).Strings(prefix, strcase.ToLowerCamel)
}

// Strings returns all stored attributes as a string map with keys unchanged.
//
// The prefix parameter is prepended to each exported key (if non-empty).
// Attributes whose rendered value is empty are skipped.
func Strings(ctx context.Context, prefix string) Map {
	return attributes(ctx).Strings(prefix, func(s string) string { return s })
}
