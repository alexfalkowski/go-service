package mvc

// StaticOption configures static file response behavior.
type StaticOption func(*staticOptions)

type staticOptions struct {
	cacheControl    string
	cacheValidators bool
}

// WithCacheControl sets the Cache-Control response header for a static route.
func WithCacheControl(value string) StaticOption {
	return func(options *staticOptions) {
		options.cacheControl = value
	}
}

// WithCacheValidators adds ETag and Last-Modified response headers for a static route.
//
// When a request's If-None-Match header matches the generated ETag, the route responds with
// 304 Not Modified and no body.
func WithCacheValidators() StaticOption {
	return func(options *staticOptions) {
		options.cacheValidators = true
	}
}

func options(opts ...StaticOption) *staticOptions {
	options := &staticOptions{}
	for _, opt := range opts {
		opt(options)
	}

	return options
}
