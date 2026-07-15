package mvc

// StaticOption configures static file response behavior.
type StaticOption func(*staticOptions)

type staticOptions struct {
	cacheControl string
}

// WithCacheControl sets the Cache-Control response header for a static route.
func WithCacheControl(value string) StaticOption {
	return func(options *staticOptions) {
		options.cacheControl = value
	}
}

func options(opts ...StaticOption) *staticOptions {
	options := &staticOptions{}
	for _, opt := range opts {
		opt(options)
	}

	return options
}
