package bytes

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires this package's `config_size` validation rule into [go.uber.org/fx].
//
// Including this module contributes a [github.com/alexfalkowski/go-service/v2/config/validate.Validation] to
// the shared "validations" group consumed by
// [github.com/alexfalkowski/go-service/v2/config/validate.NewValidator].
var Module = di.Module(
	di.Constructor(newValidation),
)
