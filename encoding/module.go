package encoding

import "github.com/alexfalkowski/go-signal/v2"

// Module provides the default encoder registry as a *[Map].
var Module = signal.Module(
	signal.Constructor(NewMap),
)
