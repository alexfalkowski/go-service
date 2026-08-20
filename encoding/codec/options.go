package codec

// Option configures one [Encoder.Encode] or [Encoder.Decode] operation. An option that only
// applies to one of the two is ignored by the other.
type Option interface {
	apply(opts *options)
}

// Options reports settings resolved for one [Encoder.Encode] or [Encoder.Decode] operation.
type Options interface {
	// DiscardUnknown reports whether unknown source members are discarded on Decode.
	DiscardUnknown() bool

	// AllowPartial reports whether messages with missing required fields may be encoded or decoded.
	AllowPartial() bool
}

type options struct {
	discardUnknown bool
	allowPartial   bool
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func (o options) DiscardUnknown() bool {
	return o.discardUnknown
}

func (o options) AllowPartial() bool {
	return o.allowPartial
}

// Apply resolves opts for one [Encoder.Encode] or [Encoder.Decode] operation.
func Apply(opts ...Option) Options {
	resolved := options{}
	for _, opt := range opts {
		opt.apply(&resolved)
	}

	return resolved
}

// WithDiscardUnknown configures decoding to ignore source members that do not
// map to the destination. It preserves all other format validation, including
// unary trailing-data checks. It has no effect on Encode.
func WithDiscardUnknown() Option {
	return optionFunc(func(options *options) {
		options.discardUnknown = true
	})
}

// WithAllowPartial configures encoding or decoding to allow messages with missing required
// fields to be encoded or decoded instead of returning an error.
func WithAllowPartial() Option {
	return optionFunc(func(options *options) {
		options.allowPartial = true
	})
}
