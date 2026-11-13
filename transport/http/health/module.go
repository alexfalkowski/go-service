package health

import "github.com/alexfalkowski/go-signal/v2"

// Module wires HTTP health endpoints into [github.com/alexfalkowski/go-signal/v2].
var Module = signal.Module(
	signal.Register(Register),
)
