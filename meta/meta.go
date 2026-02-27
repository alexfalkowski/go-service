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

// WithAttribute returns a copy of ctx with the given attribute stored under key.
//
// Attributes are stored on the context under a single internal key as a Storage map. Each call to
// WithAttribute returns a derived context containing an updated Storage.
//
// Rendering and export behavior is controlled by the provided Value. For example:
//   - Blank and Ignored values render as empty strings and are skipped by export helpers.
//   - Redacted values render as asterisks while preserving the underlying value in-context.
func WithAttribute(ctx context.Context, key string, value Value) context.Context {
	return context.WithValue(ctx, meta, attributes(ctx).Add(key, value))
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
