package compress

import "github.com/alexfalkowski/go-signal/v2"

// Module provides the default compressor registry as a *[Map].
var Module = signal.Module(
	signal.Constructor(NewMap),
)
