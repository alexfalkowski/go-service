package otlp

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires this package's `otlp_cadence` and `otlp_batch_config` validation rules into
// [go.uber.org/fx].
//
// Including this module contributes two [github.com/alexfalkowski/go-service/v2/config/validate.Validation]
// values to the shared "validations" group consumed by
// [github.com/alexfalkowski/go-service/v2/config/validate.NewValidator].
var Module = di.Module(
	di.Constructor(newCadenceValidation),
	di.Constructor(newBatchConfigValidation),
)
