package bytes

import "github.com/alexfalkowski/go-signal/v2"

// Module wires this package's `config_size` validation rule into [github.com/alexfalkowski/go-signal/v2].
//
// Including this module contributes a [github.com/alexfalkowski/go-service/v3/config/validate.Validation] to
// the shared "validations" group consumed by
// [github.com/alexfalkowski/go-service/v3/config/validate.NewValidator].
var Module = signal.Module(
	signal.Constructor(newValidation),
)
