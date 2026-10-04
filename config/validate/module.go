package validate

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires the shared Validator into [go.uber.org/fx].
var Module = di.Module(
	di.Constructor(NewValidator),
)
