package pg

import "github.com/alexfalkowski/go-service/v2/di"

// Module provides the Fx module that wires PostgreSQL SQL database support.
var Module = di.Module(
	di.Register(Register),
	di.Constructor(Open),
)
