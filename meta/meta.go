package meta

import (
	"context"

	"github.com/iancoleman/strcase"
)

// Converter for meta.
type Converter func(string) string

// NoneConverter for meta.
var NoneConverter = func(s string) string { return s }

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

// SnakeStrings for meta.
func SnakeStrings(ctx context.Context, prefix string) map[string]string {
	return Strings(ctx, prefix, strcase.ToSnake)
}

// CamelStrings for meta.
func CamelStrings(ctx context.Context, prefix string) map[string]string {
	return Strings(ctx, prefix, strcase.ToLowerCamel)
}

// Strings for meta.
func Strings(ctx context.Context, prefix string, converter Converter) map[string]string {
	as := attributes(ctx)
	ss := map[string]string{}

	for k, v := range as {
		s := v.String()
		if s != "" {
			ss[prefix+converter(k)] = s
		}
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
