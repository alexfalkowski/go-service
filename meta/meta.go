package meta

import (
	"unicode/utf8"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NoPrefix is a convenience constant for passing an empty prefix to export helpers.
const NoPrefix = strings.Empty

// contextKey is the package-private type used for metadata context storage.
type contextKey struct{}

// metaKey keeps metadata storage isolated from exported string-backed context keys.
var metaKey contextKey

// WithAttributes returns a copy of ctx with all provided attributes stored.
//
// The storage update is copy-on-write and clones existing attributes only once before applying the
// full batch. Metadata writes use pairs so hot paths can set several attributes without repeated
// map copies and context wrappers.
func WithAttributes(ctx context.Context, pairs ...Pair) context.Context {
	if len(pairs) == 0 {
		return ctx
	}

	return context.WithValue(ctx, metaKey, attributes(ctx).AddPairs(pairs...))
}

// Attribute returns the stored attribute value for key.
//
// If no attribute is present for key, Attribute returns the zero-value [Value]. The zero value renders
// like [String] with an empty value, so [Value.IsEmpty] treats missing attributes and explicitly stored
// empty values the same.
func Attribute(ctx context.Context, key string) Value {
	return attributes(ctx).Get(key)
}

// Map is a string key-value map of exported attributes.
//
// Export helpers return this type after rendering [Value] entries and applying key conversion and prefixing.
type Map map[string]string

// SnakeStrings returns all stored attributes as a string map with snake_cased keys.
//
// The prefix parameter is prepended to each exported key (if non-empty).
// Attributes whose rendered value is empty are skipped.
func SnakeStrings(ctx context.Context, prefix string) Map {
	return attributes(ctx).Strings(prefix, strings.ToSnake)
}

// CamelStrings returns all stored attributes as a string map with lowerCamelCased keys.
//
// The prefix parameter is prepended to each exported key (if non-empty).
// Attributes whose rendered value is empty are skipped.
func CamelStrings(ctx context.Context, prefix string) Map {
	return attributes(ctx).Strings(prefix, strings.ToLowerCamel)
}

// Limit bounds the length of an exported metadata value.
//
// Standard telemetry wiring resolves this value from telemetry configuration.
// Values are truncated at a valid UTF-8 boundary, while the original context
// metadata remains unchanged.
type Limit bytes.Size

// Bytes returns the limit as a byte count.
func (l Limit) Bytes() int {
	return int(l)
}

// Attributes returns context metadata as lowerCamelCase string attributes.
//
// Attribute values are bounded by limit so one metadata value cannot dominate a
// log record or telemetry payload. Metadata retained in ctx is not truncated.
func Attributes(ctx context.Context, limit Limit) Map {
	attrs := CamelStrings(ctx, NoPrefix)
	for key, value := range attrs {
		attrs[key] = truncate(value, limit)
	}

	return attrs
}

// Strings returns all stored attributes as a string map with keys unchanged.
//
// The prefix parameter is prepended to each exported key (if non-empty).
// Attributes whose rendered value is empty are skipped.
func Strings(ctx context.Context, prefix string) Map {
	return attributes(ctx).Strings(prefix, identity)
}

func identity(s string) string {
	return s
}

func truncate(value string, limit Limit) string {
	size := limit.Bytes()
	if len(value) <= size {
		return value
	}

	truncated := value[:size]
	for !utf8.ValidString(truncated) {
		_, size := utf8.DecodeLastRuneInString(truncated)
		truncated = truncated[:len(truncated)-size]
	}

	return truncated
}
