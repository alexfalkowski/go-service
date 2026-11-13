package validate

import "github.com/alexfalkowski/go-signal/v2"

// Module wires the shared Validator into [github.com/alexfalkowski/go-signal/v2].
var Module = signal.Module(
	signal.Constructor(NewValidator),
)
