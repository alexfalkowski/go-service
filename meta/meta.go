package meta

import (
	"context"

	"github.com/iancoleman/strcase"
)

type (
	// Converter for meta.
	Converter func(string) string

	// Map for meta.
	Map map[string]string

	contextKey string
)

var (
	// NoneConverter for meta.
	NoneConverter = func(s string) string { return s }

	meta = contextKey("meta")
)

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

// SnakeStrings for meta.
func SnakeStrings(ctx context.Context, prefix string) Map {
	return Strings(ctx, prefix, strcase.ToSnake)
}

// CamelStrings for meta.
func CamelStrings(ctx context.Context, prefix string) Map {
	return Strings(ctx, prefix, strcase.ToLowerCamel)
}

// Strings for meta.
func Strings(ctx context.Context, prefix string, converter Converter) Map {
	as := attributes(ctx)
	m := Map{}

	for k, v := range as {
		s := v.String()
		if s != "" {
			m[prefix+converter(k)] = s
		}
	}

	return m
}

func attributes(ctx context.Context) map[string]Valuer {
	m := ctx.Value(meta)
	if m == nil {
		return make(map[string]Valuer)
	}

	return m.(map[string]Valuer)
}
