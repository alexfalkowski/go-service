package meta

import (
	"context"
	"fmt"
	"strings"
)

// StringOrBlank for meta.
func StringOrBlank(s fmt.Stringer) string {
	if s != nil {
		return s.String()
	}

	return ""
}

// IsBlank checks for an empty string.
func IsBlank(actual fmt.Stringer) bool {
	return IsEqual(actual, "")
}

// IsEqual checks if actual is expected.
func IsEqual(actual fmt.Stringer, expected string) bool {
	return StringOrBlank(actual) == expected
}

// Value for meta.
type Value string

// String to satisfy fmt.Stringer.
func (v Value) String() string {
	return string(v)
}

// SafeValue for meta.
type SafeValue string

// String to satisfy fmt.Stringer.
func (v SafeValue) String() string {
	return strings.Repeat("*", len(string(v)))
}

type contextKey string

var meta = contextKey("meta")

// WithAttribute to meta.
func WithAttribute(ctx context.Context, key string, value fmt.Stringer) context.Context {
	attr := attributes(ctx)
	attr[key] = value

	return context.WithValue(ctx, meta, attr)
}

// Attribute of meta.
func Attribute(ctx context.Context, key string) fmt.Stringer {
	return attributes(ctx)[key]
}

// Attributes of meta.
func Attributes(ctx context.Context) map[string]fmt.Stringer {
	return attributes(ctx)
}

// Strings of meta.
func Strings(ctx context.Context) map[string]string {
	as := Attributes(ctx)
	ss := make(map[string]string, len(as))

	for k, v := range as {
		ss[k] = v.String()
	}

	return ss
}

func attributes(ctx context.Context) map[string]fmt.Stringer {
	m := ctx.Value(meta)
	if m == nil {
		return make(map[string]fmt.Stringer)
	}

	return m.(map[string]fmt.Stringer)
}
