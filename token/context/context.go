package context

import "context"

// Context is an alias of context.Context.
type Context = context.Context

type contextKey string

const options = contextKey("options")

// WithOpts for context.
func WithOpts(ctx context.Context, opts Options) context.Context {
	return context.WithValue(ctx, options, opts)
}

// Opts for context.
func Opts(ctx context.Context) Options {
	return ctx.Value(options).(Options)
}

// Opts for context.
func AddToOpts(ctx context.Context, key string, value any) context.Context {
	opts := Opts(ctx)
	opts[key] = value

	return WithOpts(ctx, opts)
}

// Options for context.
type Options map[string]any

// GetString from key.
func (o Options) GetString(key string) string {
	return o[key].(string)
}
