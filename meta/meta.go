package meta

import (
	"context"

	"github.com/iancoleman/strcase"
)

const meta = contextKey("meta")

type (
	// Converter for meta.
	Converter func(string) string

	// Map for meta.
	Map map[string]string

	contextKey string
)

// NoneConverter for meta.
var NoneConverter = func(s string) string { return s }

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
	attrs := attributes(ctx)
	attrMap := make(Map, len(attrs))

	for k, v := range attrs {
		s := v.String()
		if s == "" {
			continue
		}

		attrMap[key(prefix, converter(k))] = s
	}

	return attrMap
}

func key(prefix, key string) string {
	if prefix == "" {
		return key
	}

	return prefix + key
}

func attributes(ctx context.Context) map[string]Valuer {
	m := ctx.Value(meta)
	if m == nil {
		return make(map[string]Valuer)
	}

	return m.(map[string]Valuer)
}
