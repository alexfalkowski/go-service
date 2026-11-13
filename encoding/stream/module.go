package stream

import "github.com/alexfalkowski/go-signal/v2"

// Module provides the default streaming encoder/decoder registry as a *[Map].
var Module = signal.Module(
	signal.Constructor(NewMap),
)
