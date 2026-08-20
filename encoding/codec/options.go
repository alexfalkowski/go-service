package codec

// Option configures one [Encoder.Decode] operation.
type Option interface {
	apply(opts *options)
}

// Options reports settings resolved for one [Encoder.Decode] operation.
type Options interface {
	// DiscardUnknown reports whether unknown source members are discarded.
	DiscardUnknown() bool
}

type options struct {
	discardUnknown bool
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func (o options) DiscardUnknown() bool {
	return o.discardUnknown
}

// Apply resolves opts for one [Encoder.Decode] operation.
func Apply(opts ...Option) Options {
	resolved := options{}
	for _, opt := range opts {
		opt.apply(&resolved)
	}

	return resolved
}

// WithDiscardUnknown configures decoding to ignore source members that do not
// map to the destination. It preserves all other format validation, including
// unary trailing-data checks.
func WithDiscardUnknown() Option {
	return optionFunc(func(options *options) {
		options.discardUnknown = true
	})
}
