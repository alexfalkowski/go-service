package meta

import (
	"context"
	"fmt"
	"strings"
)

// ToValuer for meta.
func ToValuer(st fmt.Stringer) Valuer {
	if st != nil {
		return String(st.String())
	}

	return nil
}

// Valuer for meta.
type Valuer interface {
	Value() string

	fmt.Stringer
}

// ValueOrBlank for meta.
func ValueOrBlank(s Valuer) string {
	if s != nil {
		return s.Value()
	}

	return ""
}

// IsBlank checks for an empty string.
func IsBlank(actual Valuer) bool {
	return IsEqual(actual, "")
}

// IsEqual checks if actual is expected.
func IsEqual(actual Valuer, expected string) bool {
	return ValueOrBlank(actual) == expected
}

// Error for meta.
func Error(err error) Valuer {
	if err != nil {
		return String(err.Error())
	}

	return String("")
}

// String for meta.
type String string

// Value of the string.
func (v String) Value() string {
	return string(v)
}

// String to satisfy fmt.Stringer.
func (v String) String() string {
	return v.Value()
}

// Redacted for meta.
type Redacted string

// Value of the string.
func (v Redacted) Value() string {
	return string(v)
}

// String to satisfy fmt.Stringer.
func (v Redacted) String() string {
	return strings.Repeat("*", len(v.Value()))
}

type contextKey string

var meta = contextKey("meta")

// WithAttribute to meta.
func WithAttribute(ctx context.Context, key string, value Valuer) context.Context {
	attr := attributes(ctx)
	attr[key] = value

	return context.WithValue(ctx, meta, attr)
}

// Attribute of meta.
func Attribute(ctx context.Context, key string) Valuer {
	return attributes(ctx)[key]
}

// Strings of meta.
func Strings(ctx context.Context) map[string]string {
	as := attributes(ctx)
	ss := make(map[string]string, len(as))

	for k, v := range as {
		ss[k] = v.String()
	}

	return ss
}

func attributes(ctx context.Context) map[string]Valuer {
	m := ctx.Value(meta)
	if m == nil {
		return make(map[string]Valuer)
	}

	return m.(map[string]Valuer)
}
